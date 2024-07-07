package testutils

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	refreshtokenschema "github.com/gaetanDubuc/beeckend/internal/refresh-token/schema"
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

func (m *Middleware) AuthHandler(c *gin.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyUnauthenticated(c *gin.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyConfirmedUser(c *gin.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyUnconfirmedUser(c *gin.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyRootAdmin(c *gin.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) CurrentAuthenticatedUser(ctx context.Context) entity.User {
	return m.Called().Get(0).(entity.User)
}

func (m *Middleware) IngressHeaderAuthentication(c *gin.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyExpiredSession(c *gin.Context) error {
	return m.Called().Error(0)

}
