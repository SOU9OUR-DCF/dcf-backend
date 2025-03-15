package mocks

import (
	"context"

	"github.com/SOU9OUR-DCF/dcf-backend.git/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Register(ctx context.Context, req domain.RegisterRequest) (*domain.AuthResponse, domain.Token, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, domain.Token(""), args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Get(1).(domain.Token), args.Error(2)
}

func (m *AuthServiceMock) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, domain.Token, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, domain.Token(""), args.Error(1)
	}
	return args.Get(0).(*domain.AuthResponse), args.Get(1).(domain.Token), args.Error(2)
}

func (m *AuthServiceMock) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *AuthServiceMock) RefreshToken(ctx context.Context, token string) (*domain.AuthResponse, domain.Token, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, domain.Token(""), args.Error(2)
	}
	return args.Get(0).(*domain.AuthResponse), args.Get(1).(domain.Token), args.Error(2)
}

func (m *AuthServiceMock) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
