package auth

import (
	"context"
	"errors"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	refreshtokenschema "github.com/gaetanDubuc/beeckend/internal/refresh-token/schema"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	ErrWrongPassword = errors.New("wrong password")
)

type (
	UserRepository interface {
		Get(ctx context.Context, user *entity.User) error
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
		hasher              Hasher
		tokenExpiration     int
		SigningMethod       jwt.SigningMethod
		Keyfunc             func(*jwt.Token) (interface{}, error)
		logger              log.Logger
	}
)

// NewService creates a new authentication service.
func NewService(
	userRepository UserRepository,
	refreshTokenService RefreshTokenService,
	hasher Hasher,
	tokenExpiration int,
	signingMethod jwt.SigningMethod,
	keyfunc func(*jwt.Token) (interface{}, error),
	logger log.Logger,
) *Service {
	return &Service{
		userRepository,
		refreshTokenService,
		hasher,
		tokenExpiration,
		signingMethod,
		keyfunc,
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
	user := entity.User{
		Email: req.Username,
	}
	err := s.userRepository.Get(
		ctx,
		&user,
	)
	if err != nil {
		return schema.Session{}, err
	}

	if !s.hasher.AreSameHash(req.Password, user.HashedPassword) {
		return schema.Session{}, ErrWrongPassword
	}

	jwt, err := s.GenerateJWT(ctx, user)
	if err != nil {
		return schema.Session{}, err
	}

	refreshJWT, err := s.refreshTokenService.Create(
		ctx,
		refreshtokenschema.CreateRequest{
			UserID: user.ID,
		},
	)
	if err != nil {
		return schema.Session{}, err
	}

	return schema.Session{
		JWT:        jwt,
		RefreshJWT: refreshJWT,
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
	user := entity.User{
		Model: gorm.Model{
			ID: req.UserID,
		},
	}
	err = s.userRepository.Get(
		ctx,
		&user,
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
		req.UserID,
	})
}

// GenerateJWT generates a JWT token for a given user.
func (s *Service) GenerateJWT(ctx context.Context, user entity.User) (string, error) {
	claim := schema.MakeUserClaim(user, s.tokenExpiration)

	token := jwt.NewWithClaims(s.SigningMethod, claim)
	key, err := s.Keyfunc(token)
	if err != nil {
		return "", err
	}
	JWT, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return JWT, nil
}
