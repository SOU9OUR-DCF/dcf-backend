package postgres

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type restaurantRepository struct {
	db *gorm.DB
}

func NewRestaurantRepository(db *gorm.DB) ports.RestaurantRepository {
	return &restaurantRepository{db: db}
}

func (r *restaurantRepository) Create(ctx context.Context, tx interface{}, restaurant *domain.Restaurant) error {
	if tx == nil {
		return r.db.Create(restaurant).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Create(restaurant).Error
}

func (r *restaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error) {
	var restaurant domain.Restaurant
	if err := r.db.Where("id = ?", id).First(&restaurant).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *restaurantRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Restaurant, error) {
	var restaurant domain.Restaurant
	if err := r.db.Where("user_id = ?", userID).First(&restaurant).Error; err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *restaurantRepository) Update(ctx context.Context, tx interface{}, restaurant *domain.Restaurant) error {
	if tx == nil {
		return r.db.Save(restaurant).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Save(restaurant).Error
}

func (r *restaurantRepository) UpdateStats(ctx context.Context, tx interface{}, id uuid.UUID, totalEvents, mealsServed int, rating float64) error {
	updates := map[string]interface{}{
		"total_events": totalEvents,
		"meals_served": mealsServed,
		"rating":       rating,
	}

	if tx == nil {
		return r.db.Model(&domain.Restaurant{}).
			Where("id = ?", id).
			Updates(updates).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Model(&domain.Restaurant{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *restaurantRepository) Delete(ctx context.Context, tx interface{}, id uuid.UUID) error {
	if tx == nil {
		return r.db.Delete(&domain.Restaurant{}, "id = ?", id).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Delete(&domain.Restaurant{}, "id = ?", id).Error
}
