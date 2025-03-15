package application

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
)

type userService struct {
	txManager      ports.TransactionManager
	userRepo       ports.UserRepository
	restaurantRepo ports.RestaurantRepository
	volunteerRepo  ports.VolunteerRepository
}

func NewUserService(
	txManager ports.TransactionManager,
	userRepo ports.UserRepository,
	restaurantRepo ports.RestaurantRepository,
	volunteerRepo ports.VolunteerRepository,
) ports.UserService {
	return &userService{
		txManager:      txManager,
		userRepo:       userRepo,
		restaurantRepo: restaurantRepo,
		volunteerRepo:  volunteerRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*domain.User, interface{}, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, uid)
	if err != nil {
		return nil, nil, err
	}

	// Get the appropriate profile based on user type
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

func (s *userService) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.userRepo.Update(ctx, tx, user)
	})
}

func (s *userService) GetUserProfile(ctx context.Context, userID string, userType domain.UserType) (interface{}, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	switch userType {
	case domain.UserTypeRestaurant:
		return s.restaurantRepo.GetByUserID(ctx, uid)
	case domain.UserTypeVolunteer:
		return s.volunteerRepo.GetByUserID(ctx, uid)
	default:
		return nil, fmt.Errorf("unsupported user type: %s", userType)
	}
}

func (s *userService) UpdateUserProfile(ctx context.Context, userID string, profile interface{}) error {
	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		uid, err := uuid.Parse(userID)
		if err != nil {
			return fmt.Errorf("invalid user ID: %w", err)
		}

		switch p := profile.(type) {
		case *domain.Restaurant:
			if p.UserID != uid {
				return fmt.Errorf("profile does not belong to user")
			}
			return s.restaurantRepo.Update(ctx, tx, p)
		case *domain.Volunteer:
			if p.UserID != uid {
				return fmt.Errorf("profile does not belong to user")
			}
			return s.volunteerRepo.Update(ctx, tx, p)
		default:
			return fmt.Errorf("unsupported profile type")
		}
	})
}
