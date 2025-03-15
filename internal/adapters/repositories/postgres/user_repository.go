package postgres

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

var _ ports.UserRepository = &userRepository{}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, tx interface{}, user *domain.User) error {
	if tx == nil {
		return r.db.Create(user).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, tx interface{}, user *domain.User) error {
	if tx == nil {
		return r.db.Save(user).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, tx interface{}, id uuid.UUID) error {
	if tx == nil {
		return r.db.Delete(&domain.User{}, "id = ?", id).Error
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	return gormTx.Delete(&domain.User{}, "id = ?", id).Error
}
