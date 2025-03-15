package postgres

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type volunteerRepository struct {
	db *gorm.DB
}

func NewVolunteerRepository(db *gorm.DB) ports.VolunteerRepository {
	return &volunteerRepository{db: db}
}

func (r *volunteerRepository) Create(ctx context.Context, tx interface{}, volunteer *domain.Volunteer) error {
	if tx == nil {
		return r.db.Create(volunteer).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Create(volunteer).Error
}

func (r *volunteerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Volunteer, error) {
	var volunteer domain.Volunteer
	if err := r.db.Where("id = ?", id).First(&volunteer).Error; err != nil {
		return nil, err
	}
	return &volunteer, nil
}

func (r *volunteerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Volunteer, error) {
	var volunteer domain.Volunteer
	if err := r.db.Where("user_id = ?", userID).First(&volunteer).Error; err != nil {
		return nil, err
	}
	return &volunteer, nil
}

func (r *volunteerRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]*domain.Volunteer, error) {
	var volunteers []*domain.Volunteer
	if err := r.db.Joins("JOIN event_volunteers ON volunteers.id = event_volunteers.volunteer_id").
		Where("event_volunteers.event_id = ?", eventID).
		Find(&volunteers).Error; err != nil {
		return nil, err
	}
	return volunteers, nil
}

func (r *volunteerRepository) Update(ctx context.Context, tx interface{}, volunteer *domain.Volunteer) error {
	if tx == nil {
		return r.db.Save(volunteer).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Save(volunteer).Error
}

func (r *volunteerRepository) Delete(ctx context.Context, tx interface{}, id uuid.UUID) error {
	if tx == nil {
		return r.db.Delete(&domain.Volunteer{}, "id = ?", id).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Delete(&domain.Volunteer{}, "id = ?", id).Error
}

func (r *volunteerRepository) GetNearbyVolunteers(ctx context.Context, latitude, longitude float64, radiusKm int) ([]*domain.Volunteer, error) {
	// This is a simplified implementation. In a real application, you would use
	// a spatial query or a more sophisticated distance calculation.
	// For now, we'll just return all volunteers as "nearby"
	var volunteers []*domain.Volunteer
	if err := r.db.Find(&volunteers).Error; err != nil {
		return nil, err
	}
	return volunteers, nil
}

func (r *volunteerRepository) UpdateStats(ctx context.Context, tx interface{}, id uuid.UUID, tasksCompleted, hoursVolunteered, mealsServed, reputationPoints int) error {
	updates := map[string]interface{}{
		"tasks_completed":   tasksCompleted,
		"hours_volunteered": hoursVolunteered,
		"meals_served":      mealsServed,
		"reputation_points": reputationPoints,
	}

	if tx == nil {
		return r.db.Model(&domain.Volunteer{}).
			Where("id = ?", id).
			Updates(updates).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Model(&domain.Volunteer{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *volunteerRepository) CountByRestaurantID(ctx context.Context, restaurantID uuid.UUID) (int, error) {
	var count int64
	if err := r.db.Model(&domain.Volunteer{}).
		Joins("JOIN event_volunteers ON volunteers.id = event_volunteers.volunteer_id").
		Joins("JOIN events ON event_volunteers.event_id = events.id").
		Where("events.restaurant_id = ?", restaurantID).
		Distinct("volunteers.id").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
