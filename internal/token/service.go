package service

import (
	"context"
	"database/sql"
	e "errors"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/log"
)

var (
	ErrTokenInvalid = e.New("The token is invalid!")
	ErrTokenExpired = e.New("The token is expired!")
	ErrCoolDown     = e.New("wait until the cooldown is over")
)

type (
	Tokenable interface {
		GetHashedToken() string
		GetToken() string
		SetHashedToken(string)
		SetToken(string)
		GetCreatedAt() time.Time
		GetUpdatedAt() time.Time
	}

	Service[T Tokenable] struct {
		Repository[T]
		hasher Hasher
		logger log.Logger
	}

	Repository[T Tokenable] interface {
		Get(ctx context.Context, token T) error
		Create(ctx context.Context, token T) error
		Update(ctx context.Context, token T) error
		HardDelete(ctx context.Context, token T) error
	}

	// Hasher is the interface that wraps the basic crypto methods.
	Hasher interface {
		GenerateToken() string
		Hash(value string) (string, error)
		AreSameHash(value string, hashedValue string) bool
	}
)

// NewService creates a new token service.
func New[T Tokenable](
	repo Repository[T],
	hasher Hasher,
	logger log.Logger,
) *Service[T] {
	logger = logger.With(context.Background(), "service", "token")
	return &Service[T]{repo, hasher, logger}
}

// Create creates a new token.
func (s *Service[T]) Create(ctx context.Context, token T) error {
	// Generate token and store it.
	value := s.hasher.GenerateToken()
	hash, err := s.hasher.Hash(value)
	if err != nil {
		return err
	}

	token.SetHashedToken(hash)
	token.SetToken(value)

	err = s.Repository.Create(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service[T]) CreateOrUpdate(ctx context.Context, token T) error {
	// Get the current token. If no token exists then create one else update it with a new hash.
	err := s.Repository.Get(ctx,
		token,
	)

	if e.Is(err, sql.ErrNoRows) {
		err = s.Create(ctx, token)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// Generate token and store it.
	value := s.hasher.GenerateToken()
	hash, err := s.hasher.Hash(value)

	if err != nil {
		return err
	}

	token.SetHashedToken(hash)
	token.SetToken(value)

	err = s.Repository.Update(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service[T]) Confirm(ctx context.Context, expiration int, token T) error {
	err := s.Repository.Get(ctx, token)
	if err != nil {
		return err
	}

	if !s.hasher.AreSameHash(token.GetToken(), token.GetHashedToken()) {
		return ErrTokenInvalid
	}

	now := time.Now().UTC()
	if token.GetCreatedAt().Add(time.Duration(expiration)*time.Second).Before(now) &&
		token.GetUpdatedAt().Add(time.Duration(expiration)*time.Second).Before(now) {
		return ErrTokenExpired
	}

	return nil
}

func (s *Service[T]) ConfirmAndDelete(ctx context.Context, expiration int, token T) error {
	err := s.Confirm(ctx, expiration, token)
	if err != nil && err != ErrTokenExpired {
		return err
	}

	errDelete := s.Repository.HardDelete(ctx, token)
	if errDelete != nil {
		return errDelete
	}
	return err
}

func (s *Service[T]) CheckCoolDown(ctx context.Context, coolDownTime int, token T) error {
	if token.GetUpdatedAt().Add(time.Duration(coolDownTime)*time.Second).After(time.Now().UTC()) &&
		(token.GetCreatedAt() != token.GetUpdatedAt()) {
		return ErrCoolDown
	}
	return nil
}
