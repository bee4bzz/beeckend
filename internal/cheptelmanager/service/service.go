package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	log "github.com/gaetanDubuc/beeckend/pkg/log"
	"gorm.io/gorm"
)

type CheptelManager interface {
	OnlyMember(ctx context.Context, cheptelID, userID uint) error
}

type Repository interface {
	Get(ctx context.Context, cheptel *entity.Cheptel) error
	QueryByUser(ctx context.Context, user entity.User, cheptels *[]entity.Cheptel) error
	Create(ctx context.Context, cheptel *entity.Cheptel) error
	Update(ctx context.Context, cheptel *entity.Cheptel) error
	SoftDelete(ctx context.Context, cheptel *entity.Cheptel) error
}

type Service struct {
	Repository
	cheptelManager CheptelManager
	logger         *log.Logger
}

func NewService(repository Repository, cheptelManager CheptelManager, logger *log.Logger) *Service {
	return &Service{
		Repository:     repository,
		cheptelManager: cheptelManager,
		logger:         logger.Named("cheptel"),
	}
}

func (s *Service) QueryByUser(ctx context.Context, req schema.QueryRequest) ([]entity.Cheptel, error) {
	logger := s.logger.Named("QueryByUser")
	logger.Debugf("User %v query its cheptels", req.UserID)

	if err := req.Validate(); err != nil {
		logger.Error("the request is invalid")
		return []entity.Cheptel{}, err
	}

	cheptels := []entity.Cheptel{}

	err := s.Repository.QueryByUser(ctx, entity.User{Model: gorm.Model{ID: req.UserID}}, &cheptels)
	if err != nil {
		logger.Error("An error occured when getting the cheptels from the repository")
		return []entity.Cheptel{}, err
	}

	logger.Debugf("cheptels are retrieved")
	return cheptels, nil
}

// Create creates a cheptel
func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (entity.Cheptel, error) {
	if err := req.Validate(); err != nil {
		return entity.Cheptel{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.Cheptel{}, err
	}

	cheptel := entity.Cheptel{
		Model: gorm.Model{
			ID: req.CheptelID,
		},
		Name: req.Name,
	}

	err = s.Repository.Create(ctx, &cheptel)
	if err != nil {
		return entity.Cheptel{}, err
	}
	return cheptel, err
}

// Update updates a cheptel
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.Cheptel, error) {
	if err := req.Validate(); err != nil {
		return entity.Cheptel{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.Cheptel{}, err
	}

	cheptel := entity.Cheptel{
		Model: gorm.Model{
			ID: req.CheptelID,
		},
		Name: req.NewName,
	}

	err = s.Repository.Update(ctx, &cheptel)
	if err != nil {
		return entity.Cheptel{}, err
	}
	return cheptel, err
}

// Delete deletes an cheptel
func (s *Service) SoftDelete(ctx context.Context, req schema.Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return err
	}

	cheptel := entity.Cheptel{
		Model: gorm.Model{
			ID: req.CheptelID,
		},
	}

	return s.Repository.SoftDelete(ctx, &cheptel)
}
