package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/internal/refresh-token/schema"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Get(ctx context.Context, userUUID uuid.UUID) (entity.User, error)
	GetFromUsername(ctx context.Context, username string) (entity.User, error)
	Update(ctx context.Context, user entity.User, exclude ...string) error
}

type TokenService interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	ConfirmAndDelete(ctx context.Context, expiration int, token *entity.RefreshToken) error
	HardDelete(ctx context.Context, token *entity.RefreshToken) error
}

type Hasher interface {
	AreSameHash(value string, hashedValue string) bool
}

// NewService creates a new user confirmation service.
func New(
	tokenService TokenService,
	hasher Hasher,
	expiration int,
	signingMethod jwt.SigningMethod,
	keyfunc func(*jwt.Token) (interface{}, error),
	logger log.Logger,
) *Service {
	return &Service{
		tokenService,
		expiration,
		signingMethod,
		keyfunc,
	}
}

type Service struct {
	TokenService
	expiration    int
	SigningMethod jwt.SigningMethod
	Keyfunc       func(*jwt.Token) (interface{}, error)
}

func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}
	token := entity.RefreshToken{
		Token: entity.Token{
			OwnerID: req.UserID,
		},
	}
	err := s.TokenService.Create(
		ctx,
		&token,
	)
	if err != nil {
		return "", err
	}
	jwt, err := s.GenerateJWT(ctx, req.UserID, token.GetToken())
	if err != nil {
		return "", err
	}
	return jwt, err
}

func (s *Service) Refresh(ctx context.Context, req schema.RefreshRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	err := s.TokenService.ConfirmAndDelete(
		ctx,
		s.expiration,
		&entity.RefreshToken{
			Model: gorm.Model{
				ID: req.TokenID,
			},
			Token: entity.Token{
				OwnerID: req.UserID,
				Token:   req.Token,
			},
		},
	)

	if err != nil {
		return "", err
	}

	jwt, err := s.Create(
		ctx,
		schema.CreateRequest{
			UserID: req.UserID,
		},
	)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (s *Service) DeleteFromUser(ctx context.Context, req schema.DeleteFromUserRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	token := entity.RefreshToken{
		Token: entity.Token{
			OwnerID: req.UserID,
		},
	}

	// We delete it after the confirmation and before the expiration check
	// to avoid the token to be used again.
	return s.TokenService.HardDelete(ctx, &token)
}

// GenerateJWT generates a JWT token for a given user.
func (s *Service) GenerateJWT(ctx context.Context, userID uint, token string) (string, error) {
	claim := schema.MakeUserClaim(userID, s.expiration, token)

	tok := jwt.NewWithClaims(s.SigningMethod, claim)
	key, err := s.Keyfunc(tok)
	if err != nil {
		return "", err
	}
	JWT, err := tok.SignedString(key)
	if err != nil {
		return "", err
	}
	return JWT, nil
}
