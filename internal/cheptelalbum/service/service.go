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
	Get(ctx context.Context, album *entity.Album) error
	QueryByUser(ctx context.Context, user *entity.User, albums *[]entity.Album) error
	Create(ctx context.Context, album *entity.Album) error
	Update(ctx context.Context, album *entity.Album) error
	SoftDelete(ctx context.Context, album *entity.Album) error
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
		logger:         logger.Named("cheptel album"),
	}
}

func (s *Service) QueryByUser(ctx context.Context, req schema.QueryRequest) ([]entity.Album, error) {
	logger := s.logger.Named("QueryByUser")
	logger.Debugf("User %v query its albums", req.UserID)

	if err := req.Validate(); err != nil {
		logger.Error("the request is invalid")
		return []entity.Album{}, err
	}

	albums := []entity.Album{}

	err := s.Repository.QueryByUser(ctx, &entity.User{Model: gorm.Model{ID: req.UserID}}, &albums)
	if err != nil {
		logger.Error("An error occured when getting the albums from the repository")
		return []entity.Album{}, err
	}

	logger.Debugf("albums are retrieved")
	return albums, nil
}

// Create creates an album
func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (entity.Album, error) {
	if err := req.Validate(); err != nil {
		return entity.Album{}, err
	}

	album := entity.Album{
		Name:      req.Name,
		CheptelID: req.CheptelID,
	}

	if err := s.Repository.Create(ctx, &album); err != nil {
		return entity.Album{}, err
	}

	return album, nil
}
