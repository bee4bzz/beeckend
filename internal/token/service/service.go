package service

import (
	"context"

	"time"

	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	errors "gitlab.com/fogo-dev/infrastructure/web-api/internal/token"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/token/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
)

type (
	Service struct {
		TokenRepository
		hasher Hasher
		logger log.Logger
	}

	TokenRepository interface {
		Get(ctx context.Context, token *entity.Token) error
		GetBy(ctx context.Context, filters map[string]any, token *entity.Token) error
		Create(ctx context.Context, token *entity.Token, exclude ...string) error
		Delete(ctx context.Context, token *entity.Token) error
	}

	// Hasher is the interface that wraps the basic crypto methods.
	Hasher interface {
		GenerateToken() string
		Hash(value string) (string, error)
		AreSameHash(value string, hashedValue string) bool
	}
)

// NewService creates a new token service.
func NewService(
	repo TokenRepository,
	hasher Hasher,
	logger log.Logger,
) *Service {
	logger = logger.With(context.Background(), "service", "token")
	return &Service{repo, hasher, logger}
}

// Create creates a new token.
func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (entity.Token, error) {
	if err := req.Validate(); err != nil {
		return entity.Token{}, err
	}

	var emptyToken entity.Token

	// Generate token and store it.
	value := s.hasher.GenerateToken()
	hash, err := s.hasher.Hash(value)
	if err != nil {
		return emptyToken, err
	}

	now := time.Now().UTC()
	token := entity.Token{
		Base: entity.Base{
			UUID:      req.UUID.Else(entity.GenerateUUID()),
			CreatedAt: now,
			UpdatedAt: now,
		},
		OwnerUUID:   req.OwnerUUID,
		OwnerType:   req.OwnerType,
		HashedToken: hash,
		Token:       value,
	}

	err = s.TokenRepository.Create(ctx, &token)
	if err != nil {
		return emptyToken, err
	}

	return token, nil
}

func (s *Service) ConfirmToken(ctx context.Context, req schema.ConfirmRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	token := entity.Token{
		OwnerType: req.OwnerType,
	}

	err := s.TokenRepository.GetBy(ctx, map[string]any{
		token.PrimaryKey(): req.TokenUUID,
		token.ForeignKey(): req.OwnerUUID,
	}, &token)
	if err != nil {
		return err
	}

	if !s.hasher.AreSameHash(req.Token, token.HashedToken) {
		return errors.ErrTokenInvalid
	}

	// We delete it after the confirmation and before the check of expiration
	// to avoid the token to be used again.
	err = s.TokenRepository.Delete(ctx, &token)
	if err != nil {
		return err
	}

	if token.CreatedAt.Add(time.Duration(req.Expiration)*time.Second).Before(time.Now().UTC()) &&
		token.UpdatedAt.Add(time.Duration(req.Expiration)*time.Second).Before(time.Now().UTC()) {
		return errors.ErrTokenExpiredResp
	}

	return nil
}
