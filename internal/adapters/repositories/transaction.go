package repositories

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"gorm.io/gorm"
)

type gormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) ports.TransactionManager {
	return &gormTransactionManager{
		db: db,
	}
}

func (tm *gormTransactionManager) BeginTx(ctx context.Context) (interface{}, error) {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return tx, nil
}

func (tm *gormTransactionManager) CommitTx(tx interface{}) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if err := gormTx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (tm *gormTransactionManager) RollbackTx(tx interface{}) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if err := gormTx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

func (tm *gormTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx interface{}) error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		return fn(ctx, tx)
	})
}
