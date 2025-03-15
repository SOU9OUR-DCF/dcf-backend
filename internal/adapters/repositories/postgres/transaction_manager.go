package postgres

import (
	"context"
	"fmt"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/ports"
	"gorm.io/gorm"
)

type transactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) ports.TransactionManager {
	return &transactionManager{
		db: db,
	}
}

func (tm *transactionManager) BeginTx(ctx context.Context) (interface{}, error) {
	tx := tm.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return tx, nil
}

func (tm *transactionManager) CommitTx(tx interface{}) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if err := gormTx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (tm *transactionManager) RollbackTx(tx interface{}) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	if err := gormTx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

// WithTransaction provides a convenient way to execute operations within a transaction
func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx interface{}) error) error {
	tx, err := tm.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tm.RollbackTx(tx)
			panic(r) // re-throw panic after rollback
		}
	}()

	if err := fn(ctx, tx); err != nil {
		if rbErr := tm.RollbackTx(tx); rbErr != nil {
			return fmt.Errorf("error: %v, rollback error: %v", err, rbErr)
		}
		return err
	}

	if err := tm.CommitTx(tx); err != nil {
		return err
	}

	return nil
}
