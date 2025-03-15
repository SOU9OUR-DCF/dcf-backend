package ports

import (
	"context"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, tx interface{}, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, tx interface{}, user *domain.User) error
	Delete(ctx context.Context, tx interface{}, id uuid.UUID) error
}

type RestaurantRepository interface {
	Create(ctx context.Context, tx interface{}, restaurant *domain.Restaurant) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Restaurant, error)
	Update(ctx context.Context, tx interface{}, restaurant *domain.Restaurant) error
	UpdateStats(ctx context.Context, tx interface{}, id uuid.UUID, totalEvents, mealsServed int, rating float64) error
	Delete(ctx context.Context, tx interface{}, id uuid.UUID) error
}

type VolunteerRepository interface {
	Create(ctx context.Context, tx interface{}, volunteer *domain.Volunteer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Volunteer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Volunteer, error)
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.Volunteer, error)
	Update(ctx context.Context, tx interface{}, volunteer *domain.Volunteer) error
	Delete(ctx context.Context, tx interface{}, id uuid.UUID) error
	GetNearbyVolunteers(ctx context.Context, latitude, longitude float64, radiusKm int) ([]*domain.Volunteer, error)
	UpdateStats(ctx context.Context, tx interface{}, id uuid.UUID, tasksCompleted, hoursVolunteered, mealsServed, reputationPoints int) error
	CountByRestaurantID(ctx context.Context, restaurantID uuid.UUID) (int, error)
}

type EventRepository interface {
	Create(ctx context.Context, tx interface{}, event *domain.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string, limit, offset int) ([]*domain.Event, int, error)
	GetTodayEvents(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Event, error)
	Update(ctx context.Context, tx interface{}, event *domain.Event) error
	UpdateStatus(ctx context.Context, tx interface{}, id uuid.UUID, status string) error
	UpdateGuestCount(ctx context.Context, tx interface{}, id uuid.UUID, count int) error
	UpdateMealsServed(ctx context.Context, tx interface{}, id uuid.UUID, count int) error
	Delete(ctx context.Context, tx interface{}, id uuid.UUID) error
	GetUpcomingEvents(ctx context.Context) ([]*domain.Event, error)
}

type VolunteerApplicationRepository interface {
	Create(ctx context.Context, tx interface{}, application *domain.VolunteerApplication) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.VolunteerApplication, error)
	GetByVolunteerID(ctx context.Context, volunteerID uuid.UUID) ([]*domain.VolunteerApplication, error)
	GetByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]*domain.VolunteerApplication, error)
	UpdateStatus(ctx context.Context, tx interface{}, id uuid.UUID, status string) error
	Delete(ctx context.Context, tx interface{}, id uuid.UUID) error
}

type EventVolunteerRepository interface {
	Create(ctx context.Context, tx interface{}, eventVolunteer *domain.EventVolunteer) error
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.EventVolunteer, error)
	GetByVolunteerID(ctx context.Context, volunteerID uuid.UUID) ([]*domain.EventVolunteer, error)
	UpdateCheckIn(ctx context.Context, tx interface{}, id uuid.UUID, checkedIn bool) error
	Delete(ctx context.Context, tx interface{}, id uuid.UUID) error
	CountByEventID(ctx context.Context, eventID uuid.UUID) (int, error)
}

type TransactionManager interface {
	BeginTx(ctx context.Context) (interface{}, error)
	CommitTx(tx interface{}) error
	RollbackTx(tx interface{}) error
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx interface{}) error) error
}
