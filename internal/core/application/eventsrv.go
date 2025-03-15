package application

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
)

type eventService struct {
	txManager      ports.TransactionManager
	eventRepo      ports.EventRepository
	restaurantRepo ports.RestaurantRepository
}

func NewEventService(
	txManager ports.TransactionManager,
	eventRepo ports.EventRepository,
	restaurantRepo ports.RestaurantRepository,
) ports.EventService {
	return &eventService{
		txManager:      txManager,
		eventRepo:      eventRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, event *domain.Event) error {
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

func (s *eventService) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid event ID: %w", err)
	}

	return s.eventRepo.GetByID(ctx, eventID)
}

func (s *eventService) GetUpcomingEvents(ctx context.Context, restaurantID string, limit, offset int) ([]*domain.Event, int, error) {
	rid, err := uuid.Parse(restaurantID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	return s.eventRepo.GetByRestaurantID(ctx, rid, string(domain.EventStatusUpcoming), limit, offset)
}

func (s *eventService) GetTodayEvents(ctx context.Context, restaurantID string) ([]*domain.Event, error) {
	rid, err := uuid.Parse(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	return s.eventRepo.GetTodayEvents(ctx, rid)
}

func (s *eventService) UpdateEvent(ctx context.Context, event *domain.Event) error {
	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.eventRepo.Update(ctx, tx, event)
	})
}

func (s *eventService) UpdateEventStatus(ctx context.Context, id string, status domain.EventStatus) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.eventRepo.UpdateStatus(ctx, tx, eventID, string(status))
	})
}

func (s *eventService) UpdateGuestCount(ctx context.Context, id string, count int) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		return s.eventRepo.UpdateGuestCount(ctx, tx, eventID, count)
	})
}

func (s *eventService) UpdateMealsServed(ctx context.Context, id string, count int) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		// Update event meals served
		if err := s.eventRepo.UpdateMealsServed(ctx, tx, eventID, count); err != nil {
			return err
		}

		// Get the event to update restaurant stats
		event, err := s.eventRepo.GetByID(ctx, eventID)
		if err != nil {
			return err
		}

		// Update restaurant total meals served
		restaurant, err := s.restaurantRepo.GetByID(ctx, event.RestaurantID)
		if err != nil {
			return err
		}

		return s.restaurantRepo.UpdateStats(ctx, tx, restaurant.ID, restaurant.TotalEvents, restaurant.MealsServed+count, restaurant.Rating)
	})
}

func (s *eventService) DeleteEvent(ctx context.Context, id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid event ID: %w", err)
	}

	return s.txManager.WithTransaction(ctx, func(ctx context.Context, tx interface{}) error {
		// Get the event to update restaurant stats
		event, err := s.eventRepo.GetByID(ctx, eventID)
		if err != nil {
			return err
		}

		// Delete the event
		if err := s.eventRepo.Delete(ctx, tx, eventID); err != nil {
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
