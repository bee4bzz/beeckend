package service

import (
	"context"
	"database/sql"
	e "errors"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	errors "github.com/gaetanDubuc/beeckend/internal/token"
	"github.com/gaetanDubuc/beeckend/internal/token/schema"
	"gorm.io/gorm"
)

type (
	Service struct {
		TokenRepository
		hasher Hasher
		logger log.Logger
	}

	TokenRepository interface {
		Get(ctx context.Context, token *entity.Token) error
		Create(ctx context.Context, token *entity.Token, exclude ...string) error
		Update(ctx context.Context, token *entity.Token, exclude ...string) error
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

	token := entity.Token{
		OwnerID:     req.OwnerID,
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

func (s *Service) CreateOrUpdate(ctx context.Context, req schema.CreateRequest) (entity.Token, error) {
	if err := req.Validate(); err != nil {
		return entity.Token{}, err
	}

	token := &entity.Token{
		OwnerID:   req.OwnerID,
		OwnerType: req.OwnerType,
	}

	// Get the current token. If no token exists then create one else update it with a new hash.
	err := s.TokenRepository.Get(ctx,
		token,
	)

	if e.Is(err, sql.ErrNoRows) {
		*token, err = s.Create(ctx, req)
		if err != nil {
			return entity.Token{}, err
		}
		return *token, nil
	} else if err != nil {
		return entity.Token{}, err
	}

	// Generate token and store it.
	tokenValue := s.hasher.GenerateToken()
	hash, err := s.hasher.Hash(tokenValue)

	if err != nil {
		return entity.Token{}, err
	}

	token.Token = tokenValue
	token.HashedToken = hash
	token.UpdatedAt = time.Now().UTC()

	err = s.TokenRepository.Update(ctx, token)
	if err != nil {
		return entity.Token{}, err
	}

	return *token, nil
}

func (s *Service) Confirm(ctx context.Context, req schema.ConfirmRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	token := entity.Token{
		Model: gorm.Model{
			ID: req.TokenID,
		},
		OwnerID:   req.OwnerID,
		OwnerType: req.OwnerType,
	}

	err := s.TokenRepository.Get(ctx, &token)
	if err != nil {
		return err
	}

	if !s.hasher.AreSameHash(req.Token, token.HashedToken) {
		return errors.ErrTokenInvalid
	}

	if token.CreatedAt.Add(time.Duration(req.Expiration)*time.Second).Before(time.Now().UTC()) &&
		token.UpdatedAt.Add(time.Duration(req.Expiration)*time.Second).Before(time.Now().UTC()) {
		return errors.ErrTokenExpired
	}

	return nil
}

func (s *Service) ConfirmAndDelete(ctx context.Context, req schema.ConfirmRequest) error {
	err := s.Confirm(ctx, req)
	if err != nil && err != errors.ErrTokenExpired {
		return err
	}

	token := entity.Token{
		Model: gorm.Model{
			ID: req.TokenID,
		},
		OwnerID:   req.OwnerID,
		OwnerType: req.OwnerType,
	}

	// We delete it after the confirmation and before the expiration check
	// to avoid the token to be used again.
	errDelete := s.TokenRepository.Delete(ctx, &token)
	if errDelete != nil {
		return errDelete
	}
	return err
}

func (s *Service) CheckTimeOut(ctx context.Context, token entity.Token, coolDownTime int) error {
	if token.UpdatedAt.Add(time.Duration(coolDownTime)*time.Second).After(time.Now().UTC()) &&
		(token.CreatedAt != token.UpdatedAt) {
		return errors.ErrCoolDown
	}
	return nil
}
