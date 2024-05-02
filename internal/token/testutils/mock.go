package testutils

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/token/schema"
)

type Repository struct {
	mock.Mock
}

func (r *Repository) Get(ctx context.Context, token *entity.Token) error {
	args := r.Called(token)
	*token = args.Get(0).(entity.Token)
	return args.Error(1)
}

func (r *Repository) GetBy(ctx context.Context, filters map[string]any, token *entity.Token) error {
	args := r.Called(filters, token)
	*token = args.Get(0).(entity.Token)
	return args.Error(1)
}

func (r *Repository) Create(ctx context.Context, token *entity.Token, exclude ...string) error {
	args := r.Called(token.UUID, token.OwnerUUID, token.OwnerType, token.HashedToken, token.Token)
	*token = args.Get(0).(entity.Token)
	return args.Error(1)
}

func (r *Repository) Delete(ctx context.Context, token *entity.Token) error {
	args := r.Called(token)
	return args.Error(0)
}

type Service struct {
	mock.Mock
}

func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (entity.Token, error) {
	args := s.Called(req)
	t := args.Get(0).(entity.Token)
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return t, args.Error(1)
}

func (s *Service) ConfirmToken(ctx context.Context, req schema.ConfirmRequest) error {
	args := s.Called(req)
	return args.Error(0)
}
