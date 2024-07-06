package auth

import (
	"context"

	"4d63.com/optional"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/auth"
	refreshtokenschema "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/refresh-token/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
)

type (
	JWTGenerator interface {
		GenerateJWT(ctx context.Context, claim jwt.Claims) (string, error)
	}

	UserRepository interface {
		Get(ctx context.Context, UUID uuid.UUID) (entity.User, error)
		GetFromUsername(ctx context.Context, username string) (entity.User, error)
	}

	Hasher interface {
		AreSameHash(value string, hashedValue string) bool
	}

	RefreshTokenService interface {
		Refresh(ctx context.Context, req refreshtokenschema.RefreshRequest) (string, error)
		Create(ctx context.Context, req refreshtokenschema.CreateRequest) (string, error)
		DeleteFromUser(ctx context.Context, req refreshtokenschema.DeleteFromUserRequest) error
	}

	Service struct {
		userRepository      UserRepository
		refreshTokenService RefreshTokenService
		jwtGenerator        JWTGenerator
		hasher              Hasher
		tokenExpiration     int
		logger              log.Logger
	}
)

// NewService creates a new authentication service.
func NewService(
	userRepository UserRepository,
	refreshTokenService RefreshTokenService,
	jwtGenerator JWTGenerator,
	hasher Hasher,
	tokenExpiration int,
	logger log.Logger,
) *Service {
	return &Service{
		userRepository,
		refreshTokenService,
		jwtGenerator,
		hasher,
		tokenExpiration,
		logger,
	}
}

// Login authenticates a user from its username and password
// and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s *Service) Login(ctx context.Context, req schema.LoginRequest) (schema.Session, error) {
	if err := req.Validate(); err != nil {
		return schema.Session{}, err
	}
	user, err := s.userRepository.GetFromUsername(
		ctx,
		req.Username,
	)
	if err != nil {
		return schema.Session{}, err
	}

	if !s.hasher.AreSameHash(req.Password, user.Password) {
		return schema.Session{}, auth.ErrWrongPassword
	}

	jwt, err := s.GenerateJWT(ctx, user)
	if err != nil {
		return schema.Session{}, err
	}

	refreshToken, err := s.refreshTokenService.Create(
		ctx,
		refreshtokenschema.CreateRequest{
			UserUUID: user.UUID,
		},
	)
	if err != nil {
		return schema.Session{}, err
	}

	return schema.Session{
		Token:        jwt,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshSession(ctx context.Context, req refreshtokenschema.RefreshRequest) (schema.Session, error) {
	if err := req.Validate(); err != nil {
		return schema.Session{}, err
	}
	newRefreshJWT, err := s.refreshTokenService.Refresh(ctx, req)

	if err != nil {
		return schema.Session{}, err
	}

	user, err := s.userRepository.Get(
		ctx,
		req.UserUUID,
	)
	if err != nil {
		return schema.Session{}, err
	}

	newJWT, err := s.GenerateJWT(ctx, user)
	if err != nil {
		return schema.Session{}, err
	}

	return schema.Session{newJWT, newRefreshJWT}, nil
}

func (s *Service) Logout(ctx context.Context, req schema.LogoutRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return s.refreshTokenService.DeleteFromUser(ctx, refreshtokenschema.DeleteFromUserRequest{
		req.UserUUID,
	})
}

// GenerateJWT generates a JWT token for a given user.
func (s *Service) GenerateJWT(ctx context.Context, user entity.User) (string, error) {
	claim := schema.MakeUserClaim(user, s.tokenExpiration)

	JWT, err := s.jwtGenerator.GenerateJWT(ctx, claim)
	if err != nil {
		return "", err
	}
	return JWT, nil
}

type Claim interface {
	Valid() error
	GetUUID() uuid.UUID
	GetAdministrator() optional.Optional[bool]
	GetConfirmed() optional.Optional[bool]
}
