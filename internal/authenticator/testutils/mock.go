package testutils

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	refreshtokenschema "github.com/gaetanDubuc/beeckend/internal/refresh-session/schema"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type Service struct {
	mock.Mock
}

func (s *Service) Login(ctx context.Context, req schema.LoginRequest) (schema.Session, error) {
	args := s.Called(req)
	return args.Get(0).(schema.Session), args.Error(1)
}
func (s *Service) RefreshSession(ctx context.Context, req refreshtokenschema.RefreshRequest) (schema.Session, error) {
	args := s.Called(req)
	return args.Get(0).(schema.Session), args.Error(1)
}
func (s *Service) Logout(ctx context.Context, req schema.LogoutRequest) error {
	args := s.Called(req)
	return args.Error(0)
}

type Middleware struct {
	mock.Mock
}

func (m *Middleware) AuthHandler(c *gin.Context) {
	args := m.Called()
	if len(args) > 0 {
		status := args.Get(0).(int)
		c.AbortWithStatus(status)
	}
}

func (m *Middleware) AuthHandlerWithoutExpiration(c *gin.Context) {
	args := m.Called()
	if len(args) > 0 {
		status := args.Get(0).(int)
		c.AbortWithStatus(status)
	}
}

func (m *Middleware) OnlyUnauthenticated(c *gin.Context) {
	args := m.Called()
	if len(args) > 0 {
		status := args.Get(0).(int)
		c.AbortWithStatus(status)
	}
}

func (m *Middleware) CurrentAuthenticatedUser(ctx context.Context) entity.User {
	return m.Called().Get(0).(entity.User)
}

type UserRepository struct {
	mock.Mock
}

func (r *UserRepository) Get(ctx context.Context, user *entity.User) error {
	args := r.Called(user)
	*user = args.Get(0).(entity.User)
	return args.Error(1)
}

type RefreshTokenRepository struct {
	mock.Mock
}

func (r *RefreshTokenRepository) Refresh(ctx context.Context, req refreshtokenschema.RefreshRequest) (string, error) {
	return r.Called(req).Get(0).(string), r.Called(req).Error(1)
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *entity.Token) (string, error) {
	args := r.Called(token)
	return args.Get(0).(string), args.Error(1)
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token *entity.Token) error {
	return r.Called(token).Error(0)
}
