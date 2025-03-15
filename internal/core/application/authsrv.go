package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/SOU9OUR-DCF/dcf-backend.git/pkg/jwt"
	password_util "github.com/SOU9OUR-DCF/dcf-backend.git/pkg/password"
)

type authService struct {
	txManager      ports.TransactionManager
	userRepo       ports.UserRepository
	restaurantRepo ports.RestaurantRepository
	volunteerRepo  ports.VolunteerRepository
	tokenCache     ports.TokenCache
	jwtService     *jwt.Service
}

func NewAuthService(
	txManager ports.TransactionManager,
	userRepo ports.UserRepository,
	restaurantRepo ports.RestaurantRepository,
	volunteerRepo ports.VolunteerRepository,
	tokenCache ports.TokenCache,
	jwtService *jwt.Service,
) ports.AuthService {
	return &authService{
		txManager:      txManager,
		userRepo:       userRepo,
		restaurantRepo: restaurantRepo,
		volunteerRepo:  volunteerRepo,
		tokenCache:     tokenCache,
		jwtService:     jwtService,
	}
}

func (s *authService) registerUser(ctx context.Context, tx interface{}, email, username, password string, userType domain.UserType) (*domain.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	existingUsername, _ := s.userRepo.GetByUsername(ctx, username)
	if existingUsername != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := password_util.Hash(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
		Type:      userType,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.userRepo.Create(ctx, tx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) createAuthResponse(ctx context.Context, user *domain.User, profile interface{}) (*domain.AuthResponse, domain.Token, error) {
	token, exp, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, domain.Token(""), err
	}

	if err := s.tokenCache.StoreToken(ctx, user.ID, token, time.Until(exp)); err != nil {
		return nil, domain.Token(""), err
	}

	return &domain.AuthResponse{
		ExpiresAt: exp,
		User:      *user,
		Profile:   profile,
	}, domain.NewToken(token), nil
}

func (s *authService) RegisterRestaurant(ctx context.Context, req domain.RestaurantRegisterRequest) (*domain.AuthResponse, domain.Token, error) {
	var user *domain.User
	var profile *domain.Restaurant

	err := s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		var err error
		// Register user with transaction
		user, err = s.registerUser(ctx, tx, req.Email, req.Username, req.Password, domain.UserTypeRestaurant)
		if err != nil {
			return err
		}

		// Create restaurant profile with the same transaction
		profile = &domain.Restaurant{
			UserID:        user.ID,
			Name:          req.Name,
			Address:       req.Address,
			ContactNumber: req.ContactNumber,
		}

		return s.restaurantRepo.Create(ctx, tx, profile)
	})

	if err != nil {
		return nil, domain.Token(""), err
	}

	return s.createAuthResponse(ctx, user, profile)
}

func (s *authService) RegisterVolunteer(ctx context.Context, req domain.VolunteerRegisterRequest) (*domain.AuthResponse, domain.Token, error) {
	var user *domain.User
	var profile *domain.Volunteer

	err := s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		var err error
		// Register user with transaction
		user, err = s.registerUser(ctx, tx, req.Email, req.Username, req.Password, domain.UserTypeVolunteer)
		if err != nil {
			return err
		}

		// Create volunteer profile with the same transaction
		profile = &domain.Volunteer{
			UserID:      user.ID,
			FullName:    req.FullName,
			PhoneNumber: req.PhoneNumber,
		}

		return s.volunteerRepo.Create(ctx, tx, profile)
	})

	if err != nil {
		return nil, domain.Token(""), err
	}

	return s.createAuthResponse(ctx, user, profile)
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, domain.Token, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.Token(""), errors.New("invalid credentials")
	}

	if !password_util.Verify(req.Password, user.Password) {
		return nil, domain.Token(""), errors.New("invalid credentials")
	}

	var profile interface{}
	switch user.Type {
	case domain.UserTypeRestaurant:
		profile, err = s.restaurantRepo.GetByUserID(ctx, user.ID)
	case domain.UserTypeVolunteer:
		profile, err = s.volunteerRepo.GetByUserID(ctx, user.ID)
	}

	if err != nil {
		return nil, domain.Token(""), err
	}

	return s.createAuthResponse(ctx, user, profile)
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*domain.User, interface{}, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, nil, err
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, nil, err
	}

	exists, err := s.tokenCache.TokenExists(ctx, token)
	if err != nil || !exists {
		return nil, nil, errors.New("token invalid or expired")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	var profile interface{}
	switch user.Type {
	case domain.UserTypeRestaurant:
		profile, err = s.restaurantRepo.GetByUserID(ctx, user.ID)
	case domain.UserTypeVolunteer:
		profile, err = s.volunteerRepo.GetByUserID(ctx, user.ID)
	}

	if err != nil {
		return nil, nil, err
	}

	return user, profile, nil
}

func (s *authService) RefreshToken(ctx context.Context, token string) (*domain.AuthResponse, domain.Token, error) {
	user, _, err := s.ValidateToken(ctx, token)
	if err != nil {
		return nil, domain.Token(""), err
	}

	if err := s.tokenCache.InvalidateToken(ctx, user.ID); err != nil {
		return nil, domain.Token(""), err
	}

	newToken, exp, err := s.jwtService.GenerateToken(user.ID.String())
	if err != nil {
		return nil, domain.Token(""), err
	}

	if err := s.tokenCache.StoreToken(ctx, user.ID, newToken, time.Until(exp)); err != nil {
		return nil, domain.Token(""), err
	}

	return &domain.AuthResponse{
		ExpiresAt: exp,
		User:      *user,
	}, domain.Token(""), nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return err
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return err
	}

	return s.tokenCache.InvalidateToken(ctx, userID)
}
