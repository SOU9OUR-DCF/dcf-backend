package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type TokenCacheMock struct {
	mock.Mock
}

func (m *TokenCacheMock) StoreToken(ctx context.Context, userID uuid.UUID, token string, expiration time.Duration) error {
	args := m.Called(ctx, userID, token, expiration)
	return args.Error(0)
}

func (m *TokenCacheMock) GetToken(ctx context.Context, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *TokenCacheMock) InvalidateToken(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *TokenCacheMock) TokenExists(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}
