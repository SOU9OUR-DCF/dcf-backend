package application

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
)

type restaurantService struct {
	txManager      ports.TransactionManager
	restaurantRepo ports.RestaurantRepository
	eventRepo      ports.EventRepository
	volunteerRepo  ports.VolunteerRepository
	appRepo        ports.VolunteerApplicationRepository
	eventVolRepo   ports.EventVolunteerRepository
}

func NewRestaurantService(
	txManager ports.TransactionManager,
	restaurantRepo ports.RestaurantRepository,
	eventRepo ports.EventRepository,
	volunteerRepo ports.VolunteerRepository,
	appRepo ports.VolunteerApplicationRepository,
	eventVolRepo ports.EventVolunteerRepository,
) ports.RestaurantService {
	return &restaurantService{
		txManager:      txManager,
		restaurantRepo: restaurantRepo,
		eventRepo:      eventRepo,
		volunteerRepo:  volunteerRepo,
		appRepo:        appRepo,
		eventVolRepo:   eventVolRepo,
	}
}

func (s *restaurantService) GetRestaurantByUserID(ctx context.Context, userID string) (*domain.Restaurant, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return s.restaurantRepo.GetByUserID(ctx, uid)
}

func (s *restaurantService) GetRestaurantStats(ctx context.Context, restaurantID string) (map[string]interface{}, error) {
	rid, err := uuid.Parse(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	restaurant, err := s.restaurantRepo.GetByID(ctx, rid)
	if err != nil {
		return nil, err
	}

	// Get upcoming events count
	upcomingEvents, _, err := s.eventRepo.GetByRestaurantID(ctx, rid, string(domain.EventStatusUpcoming), 100, 0)
	if err != nil {
		return nil, err
	}

	// Get volunteer count
	var totalVolunteers int
	for _, event := range upcomingEvents {
		volunteers, err := s.eventVolRepo.GetByEventID(ctx, event.ID)
		if err != nil {
			return nil, err
		}
		totalVolunteers += len(volunteers)
	}

	// Get pending applications
	pendingApps, err := s.appRepo.GetByRestaurantID(ctx, rid, "pending")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_events":       restaurant.TotalEvents,
		"meals_served":       restaurant.MealsServed,
		"rating":             restaurant.Rating,
		"upcoming_events":    len(upcomingEvents),
		"volunteers_engaged": totalVolunteers,
		"pending_apps":       len(pendingApps),
	}, nil
}

func (s *restaurantService) UpdateRestaurant(ctx context.Context, restaurant *domain.Restaurant) error {
	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.restaurantRepo.Update(ctx, tx, restaurant)
	})
}

func (s *restaurantService) CreateEvent(ctx context.Context, event *domain.Event) error {
	// Set initial status
	event.Status = domain.EventStatusUpcoming

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		// Create the event
		if err := s.eventRepo.Create(ctx, tx, event); err != nil {
			return err
		}

		// Update restaurant stats
		restaurant, err := s.restaurantRepo.GetByID(ctx, event.RestaurantID)
		if err != nil {
			return err
		}

		return s.restaurantRepo.UpdateStats(ctx, tx, restaurant.ID, restaurant.TotalEvents+1, restaurant.MealsServed, restaurant.Rating)
	})
}

func (s *restaurantService) UpdateEvent(ctx context.Context, event *domain.Event) error {
	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.eventRepo.Update(ctx, tx, event)
	})
}

func (s *restaurantService) DeleteEvent(ctx context.Context, eventID string) error {
	id, err := uuid.Parse(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		// Get the event first to get restaurant ID
		event, err := s.eventRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}

		// Delete the event
		if err := s.eventRepo.Delete(ctx, tx, id); err != nil {
			return err
		}

		// Update restaurant stats
		restaurant, err := s.restaurantRepo.GetByID(ctx, event.RestaurantID)
		if err != nil {
			return err
		}

		return s.restaurantRepo.UpdateStats(ctx, tx, restaurant.ID, restaurant.TotalEvents-1, restaurant.MealsServed, restaurant.Rating)
	})
}
