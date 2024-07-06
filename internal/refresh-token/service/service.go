package service

import (
	"context"
	"time"
	ti "time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/refresh-token/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	tokenschema "gitlab.com/fogo-dev/infrastructure/web-api/internal/token/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
	val "gitlab.com/fogo-dev/infrastructure/web-api/pkg/validation"
)

type JWTGenerator interface {
	GenerateJWT(ctx context.Context, claim jwt.Claims) (string, error)
}

type UserRepository interface {
	Get(ctx context.Context, userUUID uuid.UUID) (entity.User, error)
	GetFromUsername(ctx context.Context, username string) (entity.User, error)
	Update(ctx context.Context, user entity.User, exclude ...string) error
}

type TokenService interface {
	Create(ctx context.Context, req tokenschema.CreateRequest) (entity.Token, error)
	ConfirmAndDelete(ctx context.Context, req tokenschema.ConfirmRequest) error
	DeleteBy(ctx context.Context, filters map[string]any, tableName string) error
}

type Hasher interface {
	AreSameHash(value string, hashedValue string) bool
}

// NewService creates a new user confirmation service.
func NewService(
	tokenService TokenService,
	jwtGenerator JWTGenerator,
	hasher Hasher,
	expiration int,
	logger log.Logger,
) *Service {
	return &Service{
		tokenService,
		jwtGenerator,
		expiration,
	}
}

type Service struct {
	TokenService
	jwtGenerator JWTGenerator
	expiration   int
}

func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}
	token, err := s.TokenService.Create(
		ctx,
		tokenschema.CreateRequest{
			OwnerUUID: req.UserUUID,
			OwnerType: entity.RefreshTokenType,
		},
	)
	if err != nil {
		return "", err
	}
	jwt, err := s.GenerateJWT(ctx, token)
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
		tokenschema.ConfirmRequest{
			OwnerUUID:  req.UserUUID,
			OwnerType:  entity.RefreshTokenType,
			Expiration: s.expiration,
			TokenUUID:  req.TokenUUID,
			Token:      req.Token,
		},
	)

	if err != nil {
		return "", err
	}

	jwt, err := s.Create(
		ctx,
		schema.CreateRequest{
			UserUUID: req.UserUUID,
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
	err := s.TokenService.DeleteBy(
		ctx,
		map[string]any{"owner_UUID": req.UserUUID.String()},
		string(entity.RefreshTokenType),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GenerateJWT(ctx context.Context, token entity.Token) (string, error) {
	claim := makeTokenClaim(token, s.expiration)
	JWT, err := s.jwtGenerator.GenerateJWT(ctx, claim)
	if err != nil {
		return "", err
	}
	return JWT, nil
}

func makeTokenClaim(refreshToken entity.Token, expiration int) *RefreshJWTClaims {
	return &RefreshJWTClaims{
		UUID:  refreshToken.UUID,
		Token: refreshToken.Token,
		Exp:   time.Now().UTC().Add(ti.Duration(expiration) * ti.Second).Unix(),
	}
}

type RefreshJWTClaims struct {
	UUID  uuid.UUID `json:"UUID"`
	Token string    `json:"token"`
	Exp   int64     `json:"exp"`
}

// Validate RefreshJWTClaims structure.
func (c *RefreshJWTClaims) Valid() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.UUID, val.NotNilUUID),
		validation.Field(&c.Token, validation.Required),
		validation.Field(&c.Exp, validation.Required, validation.Min(time.Now().Unix()).Exclusive().Error("Token is expired")),
	)
}
