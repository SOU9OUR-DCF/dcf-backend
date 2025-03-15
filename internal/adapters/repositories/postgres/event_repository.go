package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) ports.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, tx interface{}, event *domain.Event) error {
	if tx == nil {
		return r.db.Create(event).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Create(event).Error
}

func (r *eventRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	var event domain.Event
	if err := r.db.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) GetByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string, limit, offset int) ([]*domain.Event, int, error) {
	var events []*domain.Event
	var count int64

	query := r.db.Model(&domain.Event{}).Where("restaurant_id = ?", restaurantID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("date DESC, start_time DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, int(count), nil
}

func (r *eventRepository) GetTodayEvents(ctx context.Context, restaurantID uuid.UUID) ([]*domain.Event, error) {
	var events []*domain.Event
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	if err := r.db.Where("restaurant_id = ? AND start_time >= ? AND start_time < ?",
		restaurantID, today, tomorrow).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, tx interface{}, event *domain.Event) error {
	if tx == nil {
		return r.db.Save(event).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Save(event).Error
}

func (r *eventRepository) UpdateStatus(ctx context.Context, tx interface{}, id uuid.UUID, status string) error {
	if tx == nil {
		return r.db.Model(&domain.Event{}).
			Where("id = ?", id).
			Update("status", status).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Model(&domain.Event{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *eventRepository) UpdateGuestCount(ctx context.Context, tx interface{}, id uuid.UUID, count int) error {
	if tx == nil {
		return r.db.Model(&domain.Event{}).
			Where("id = ?", id).
			Update("guest_count", count).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Model(&domain.Event{}).
		Where("id = ?", id).
		Update("guest_count", count).Error
}

func (r *eventRepository) UpdateMealsServed(ctx context.Context, tx interface{}, id uuid.UUID, count int) error {
	if tx == nil {
		return r.db.Model(&domain.Event{}).
			Where("id = ?", id).
			Update("meals_served", count).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Model(&domain.Event{}).
		Where("id = ?", id).
		Update("meals_served", count).Error
}

func (r *eventRepository) Delete(ctx context.Context, tx interface{}, id uuid.UUID) error {
	if tx == nil {
		return r.db.Delete(&domain.Event{}, "id = ?", id).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Delete(&domain.Event{}, "id = ?", id).Error
}

func (r *eventRepository) GetUpcomingEvents(ctx context.Context) ([]*domain.Event, error) {
	var events []*domain.Event

	// Get events that are upcoming and have not reached max volunteers
	if err := r.db.Where("status = ?", domain.EventStatusUpcoming).
		Where("start_time > ?", time.Now()).
		Order("start_time asc").
		Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}
