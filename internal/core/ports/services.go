package ports

import (
	"context"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
)

type AuthService interface {
	RegisterRestaurant(ctx context.Context, req domain.RestaurantRegisterRequest) (*domain.AuthResponse, domain.Token, error)
	RegisterVolunteer(ctx context.Context, req domain.VolunteerRegisterRequest) (*domain.AuthResponse, domain.Token, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, domain.Token, error)
	ValidateToken(ctx context.Context, token string) (*domain.User, interface{}, error)
	RefreshToken(ctx context.Context, token string) (*domain.AuthResponse, domain.Token, error)
	Logout(ctx context.Context, token string) error
}

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.User, interface{}, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	GetUserProfile(ctx context.Context, userID string, userType domain.UserType) (interface{}, error)
	UpdateUserProfile(ctx context.Context, userID string, profile interface{}) error
}

type RestaurantService interface {
	GetRestaurantByUserID(ctx context.Context, userID string) (*domain.Restaurant, error)
	GetRestaurantStats(ctx context.Context, restaurantID string) (map[string]interface{}, error)
	UpdateRestaurant(ctx context.Context, restaurant *domain.Restaurant) error
}

type EventService interface {
	CreateEvent(ctx context.Context, event *domain.Event) error
	GetEventByID(ctx context.Context, id string) (*domain.Event, error)
	GetUpcomingEvents(ctx context.Context, restaurantID string, limit, offset int) ([]*domain.Event, int, error)
	GetTodayEvents(ctx context.Context, restaurantID string) ([]*domain.Event, error)
	UpdateEvent(ctx context.Context, event *domain.Event) error
	UpdateEventStatus(ctx context.Context, id string, status domain.EventStatus) error
	UpdateGuestCount(ctx context.Context, id string, count int) error
	UpdateMealsServed(ctx context.Context, id string, count int) error
	DeleteEvent(ctx context.Context, id string) error
}

type VolunteerService interface {
	GetVolunteerByUserID(ctx context.Context, userID string) (*domain.Volunteer, error)
	GetEventVolunteers(ctx context.Context, eventID string) ([]*domain.Volunteer, error)
	GetPendingApplications(ctx context.Context, restaurantID string) ([]*domain.VolunteerApplication, error)
	ApproveApplication(ctx context.Context, applicationID string) error
	DeclineApplication(ctx context.Context, applicationID string) error
	GetVolunteerCount(ctx context.Context, restaurantID string) (int, error)
	GetVolunteerDashboard(ctx context.Context, volunteerID string) (map[string]interface{}, error)
	GetUpcomingTasks(ctx context.Context, volunteerID string) ([]map[string]interface{}, error)
	GetNearbyOpportunities(ctx context.Context, volunteerID string) ([]map[string]interface{}, error)
	GetVolunteerBadges(ctx context.Context, volunteerID string) ([]map[string]interface{}, error)
	ApplyForEvent(ctx context.Context, volunteerID string, eventID string, role string) error
	CheckInForEvent(ctx context.Context, volunteerID string, eventVolunteerID string) error
}
