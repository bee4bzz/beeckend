package service

import (
	"context"
	"errors"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	refreshtokenschema "github.com/gaetanDubuc/beeckend/internal/refresh-session/schema"
	"github.com/golang-jwt/jwt/v5"
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
		userRepository  UserRepository
		hasher          Hasher
		tokenExpiration int
		SigningMethod   jwt.SigningMethod
		Keyfunc         func(*jwt.Token) (interface{}, error)
		logger          log.Logger
	}
)

// New creates a new authentication service.
func New(
	userRepository UserRepository,
	hasher Hasher,
	tokenExpiration int,
	signingMethod jwt.SigningMethod,
	keyfunc func(*jwt.Token) (interface{}, error),
	logger log.Logger,
) *Service {
	return &Service{
		userRepository,
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
func (s *Service) Login(ctx context.Context, req schema.LoginRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}
	user := entity.User{
		Email: req.Username,
	}
	err := s.userRepository.Get(
		ctx,
		&user,
	)
	if err != nil {
		return "", err
	}

	if !s.hasher.AreSameHash(req.Password, user.HashedPassword) {
		return "", ErrWrongPassword
	}

	jwt, err := s.GenerateJWT(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

// GenerateJWT generates a JWT token for a given user.
func (s *Service) GenerateJWT(ctx context.Context, userID uint) (string, error) {
	claim := schema.MakeUserClaim(userID, s.tokenExpiration)

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
