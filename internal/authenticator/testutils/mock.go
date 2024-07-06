package testutils

import (
	"context"
	"errors"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/mock"
	refreshtokenschema "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/refresh-token/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
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

func (m *Middleware) AuthHandler(c *routing.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyUnauthenticated(c *routing.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyConfirmedUser(c *routing.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyUnconfirmedUser(c *routing.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyRootAdmin(c *routing.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) CurrentAuthenticatedUser(ctx context.Context) entity.User {
	return m.Called().Get(0).(entity.User)
}

func (m *Middleware) IngressHeaderAuthentication(c *routing.Context) error {
	return m.Called().Error(0)
}

func (m *Middleware) OnlyExpiredSession(c *routing.Context) error {
	return m.Called().Error(0)

}

type JWTGenerator struct {
	mock.Mock
}

func (j *JWTGenerator) GenerateJWT(ctx context.Context, claim jwt.Claims) (string, error) {
	c, ok := claim.(*schema.JWTClaims)
	if !ok {
		return "", errors.New("invalid claim")
	}
	args := j.Called(c.UUID)
	return args.String(0), args.Error(1)
}
