package service

import (
	"context"
	e "errors"

	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/errors"
	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	log "github.com/gaetanDubuc/beeckend/pkg/log"
	"gorm.io/gorm"
)

type CheptelManager interface {
	OnlyMember(ctx context.Context, cheptelID, userID uint) error
}

type Repository interface {
	Get(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error
	Create(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error
	Update(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error
	SoftDelete(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error
}

type Service struct {
	Repository
	logger *log.Logger
}

func NewService(repository Repository, logger *log.Logger) *Service {
	return &Service{
		Repository: repository,
		logger:     logger.Named("cheptel manager"),
	}
}

// Create creates a cheptel association with a user
func (s *Service) Create(ctx context.Context, req schema.Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	err := s.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return err
	}

	cheptel := entity.Cheptel{
		Model: gorm.Model{
			ID: req.CheptelID,
		},
	}

	return s.Repository.Create(ctx, &entity.User{Model: gorm.Model{ID: req.MemberID}}, &cheptel)
}

// Delete deletes a cheptel association with a user
func (s *Service) SoftDelete(ctx context.Context, req schema.Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	err := s.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return err
	}

	cheptel := entity.Cheptel{
		Model: gorm.Model{
			ID: req.CheptelID,
		},
	}

	return s.Repository.SoftDelete(ctx, &entity.User{Model: gorm.Model{ID: req.MemberID}}, &cheptel)
}

func (s *Service) OnlyMember(ctx context.Context, cheptelID, userID uint) error {
	err := s.Repository.Get(
		ctx,
		&entity.User{Model: gorm.Model{ID: userID}},
		&entity.Cheptel{
			Model: gorm.Model{
				ID: cheptelID,
			},
		})
	if e.Is(err, gorm.ErrRecordNotFound) {
		return errors.ErrNotMember
	}
	return err
}
