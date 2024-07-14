package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/token/schema"
	"github.com/stretchr/testify/mock"
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
	fmt.Print(args.Get(0).(entity.Token))
	*token = args.Get(0).(entity.Token)
	return args.Error(1)
}

func (r *Repository) Create(ctx context.Context, token *entity.Token, exclude ...string) error {
	args := r.Called(token.ID, token.OwnerID, token.OwnerType, token.HashedToken, token.Token)
	*token = args.Get(0).(entity.Token)
	return args.Error(1)
}

func (r *Repository) Update(ctx context.Context, token *entity.Token, exclude ...string) error {
	args := r.Called(token.ID, token.OwnerID, token.OwnerType, token.HashedToken, token.Token)
	return args.Error(0)
}

func (r *Repository) Delete(ctx context.Context, token *entity.Token) error {
	args := r.Called(token)
	return args.Error(0)
}

func (r *Repository) DeleteBy(ctx context.Context, filters map[string]any, tableName string) error {
	args := r.Called(filters, tableName)
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

func (s *Service) ConfirmAndDelete(ctx context.Context, req schema.ConfirmRequest) error {
	args := s.Called(req)
	return args.Error(0)
}

func (s *Service) Confirm(ctx context.Context, req schema.ConfirmRequest) error {
	args := s.Called(req)
	return args.Error(0)
}

func (s *Service) CreateOrUpdate(ctx context.Context, req schema.CreateRequest) (entity.Token, error) {
	args := s.Called(req)
	return args.Get(0).(entity.Token), args.Error(1)
}

func (s *Service) CheckTimeOut(ctx context.Context, token entity.Token, coolDownTime int) error {
	args := s.Called(token, coolDownTime)
	return args.Error(0)
}

type Hasher struct {
	mock.Mock
}

func (m *Hasher) GenerateToken() string {
	ret := m.Called()
	return ret.Get(0).(string)
}

func (m *Hasher) Hash(value string) (string, error) {
	ret := m.Called(value)
	return ret.String(0), ret.Error(1)
}

func (m *Hasher) AreSameHash(value string, hashedValue string) bool {
	ret := m.Called(value, hashedValue)
	return ret.Bool(0)
}
