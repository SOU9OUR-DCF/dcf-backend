package postgres

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type eventVolunteerRepository struct {
	db *gorm.DB
}

func NewEventVolunteerRepository(db *gorm.DB) ports.EventVolunteerRepository {
	return &eventVolunteerRepository{db: db}
}

func (r *eventVolunteerRepository) Create(ctx context.Context, tx interface{}, eventVolunteer *domain.EventVolunteer) error {
	if tx == nil {
		return r.db.Create(eventVolunteer).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Create(eventVolunteer).Error
}

func (r *eventVolunteerRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.EventVolunteer, error) {
	var eventVolunteers []*domain.EventVolunteer
	if err := r.db.Where("event_id = ?", eventID).Find(&eventVolunteers).Error; err != nil {
		return nil, err
	}
	return eventVolunteers, nil
}

func (r *eventVolunteerRepository) GetByVolunteerID(ctx context.Context, volunteerID uuid.UUID) ([]*domain.EventVolunteer, error) {
	var eventVolunteers []*domain.EventVolunteer
	if err := r.db.Where("volunteer_id = ?", volunteerID).Find(&eventVolunteers).Error; err != nil {
		return nil, err
	}
	return eventVolunteers, nil
}

func (r *eventVolunteerRepository) UpdateCheckIn(ctx context.Context, tx interface{}, id uuid.UUID, checkedIn bool) error {
	if tx == nil {
		return r.db.Model(&domain.EventVolunteer{}).
			Where("id = ?", id).
			Update("checked_in", checkedIn).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Model(&domain.EventVolunteer{}).
		Where("id = ?", id).
		Update("checked_in", checkedIn).Error
}

func (r *eventVolunteerRepository) Delete(ctx context.Context, tx interface{}, id uuid.UUID) error {
	if tx == nil {
		return r.db.Delete(&domain.EventVolunteer{}, "id = ?", id).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Delete(&domain.EventVolunteer{}, "id = ?", id).Error
}

func (r *eventVolunteerRepository) CountByEventID(ctx context.Context, eventID uuid.UUID) (int, error) {
	var count int64
	if err := r.db.Model(&domain.EventVolunteer{}).Where("event_id = ?", eventID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
