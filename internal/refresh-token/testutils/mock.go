package testutils

import (
	"context"

	"github.com/stretchr/testify/mock"
	refreshtokenschema "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/refresh-token/schema"
)

type Repository struct {
	mock.Mock
}

func (m *Repository) Refresh(ctx context.Context, req refreshtokenschema.RefreshRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *Repository) Create(ctx context.Context, req refreshtokenschema.CreateRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *Repository) DeleteFromUser(ctx context.Context, req refreshtokenschema.DeleteFromUserRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
