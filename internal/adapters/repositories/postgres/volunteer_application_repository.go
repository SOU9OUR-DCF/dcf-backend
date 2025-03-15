package postgres

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type volunteerApplicationRepository struct {
	db *gorm.DB
}

func NewVolunteerApplicationRepository(db *gorm.DB) ports.VolunteerApplicationRepository {
	return &volunteerApplicationRepository{db: db}
}

func (r *volunteerApplicationRepository) Create(ctx context.Context, tx interface{}, application *domain.VolunteerApplication) error {
	if tx == nil {
		return r.db.Create(application).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Create(application).Error
}

func (r *volunteerApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.VolunteerApplication, error) {
	var app domain.VolunteerApplication
	if err := r.db.Where("id = ?", id).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *volunteerApplicationRepository) GetByEventID(ctx context.Context, eventID uuid.UUID, status string) ([]*domain.VolunteerApplication, error) {
	var apps []*domain.VolunteerApplication
	query := r.db.Where("event_id = ?", eventID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *volunteerApplicationRepository) GetByRestaurantID(ctx context.Context, restaurantID uuid.UUID, status string) ([]*domain.VolunteerApplication, error) {
	var apps []*domain.VolunteerApplication

	query := r.db.Joins("JOIN events ON volunteer_applications.event_id = events.id").
		Where("events.restaurant_id = ?", restaurantID)

	if status != "" {
		query = query.Where("volunteer_applications.status = ?", status)
	}

	if err := query.Preload("Volunteer").Preload("Event").Find(&apps).Error; err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *volunteerApplicationRepository) GetByVolunteerID(ctx context.Context, volunteerID uuid.UUID) ([]*domain.VolunteerApplication, error) {
	var apps []*domain.VolunteerApplication
	if err := r.db.Where("volunteer_id = ?", volunteerID).Find(&apps).Error; err != nil {
		return nil, err
	}
	return apps, nil
}

func (r *volunteerApplicationRepository) UpdateStatus(ctx context.Context, tx interface{}, id uuid.UUID, status string) error {
	if tx == nil {
		return r.db.Model(&domain.VolunteerApplication{}).
			Where("id = ?", id).
			Update("status", status).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Model(&domain.VolunteerApplication{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *volunteerApplicationRepository) Delete(ctx context.Context, tx interface{}, id uuid.UUID) error {
	if tx == nil {
		return r.db.Delete(&domain.VolunteerApplication{}, "id = ?", id).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Delete(&domain.VolunteerApplication{}, "id = ?", id).Error
}
