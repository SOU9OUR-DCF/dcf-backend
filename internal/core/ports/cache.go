package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenCache interface {
	StoreToken(ctx context.Context, userID uuid.UUID, token string, expiration time.Duration) error
	GetToken(ctx context.Context, userID uuid.UUID) (string, error)
	InvalidateToken(ctx context.Context, userID uuid.UUID) error
	TokenExists(ctx context.Context, token string) (bool, error)
}
