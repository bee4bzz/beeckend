package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/cheptelalbum/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	log "github.com/gaetanDubuc/beeckend/pkg/log"
	"gorm.io/gorm"
)

type (
	CheptelManager interface {
		OnlyMember(ctx context.Context, cheptelID, userID uint) error
	}

	CheptelRepository interface {
		QueryByUser(ctx context.Context, user *entity.User, cheptels *[]entity.Cheptel) error
	}

	Repository interface {
		Get(ctx context.Context, album *entity.Album) error
		QueryByOwnerIDs(ctx context.Context, ownerIDs []any, ownerType entity.AlbumType, albums *[]entity.Album) error
		Create(ctx context.Context, album *entity.Album) error
		Update(ctx context.Context, album *entity.Album) error
		SoftDelete(ctx context.Context, album *entity.Album) error
	}

	Service struct {
		Repository
		cheptelRepository CheptelRepository
		cheptelManager    CheptelManager
		logger            *log.Logger
	}
)

func NewService(repository Repository, cheptelRepository CheptelRepository, cheptelManager CheptelManager, logger *log.Logger) *Service {
	return &Service{
		Repository:        repository,
		cheptelRepository: cheptelRepository,
		cheptelManager:    cheptelManager,
		logger:            logger.Named("cheptel album"),
	}
}

func (s *Service) QueryByUser(ctx context.Context, req schema.QueryRequest) ([]entity.Album, error) {
	logger := s.logger.Named("QueryByUser")
	logger.Debugf("User %v query its cheptel albums", req.UserID)

	if err := req.Validate(); err != nil {
		logger.Error("the request is invalid")
		return []entity.Album{}, err
	}

	cheptels := []entity.Cheptel{}

	err := s.cheptelRepository.QueryByUser(ctx, &entity.User{Model: gorm.Model{ID: req.UserID}}, &cheptels)
	if err != nil {
		logger.Error("An error occured when getting the cheptels from the repository")
		return []entity.Album{}, err
	}

	cheptelIDs := []any{}
	for _, cheptel := range cheptels {
		cheptelIDs = append(cheptelIDs, cheptel.ID)
	}

	albums := []entity.Album{}

	err = s.Repository.QueryByOwnerIDs(ctx, cheptelIDs, entity.Cheptels, &albums)
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

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.Album{}, err
	}

	album := entity.Album{
		Name:        req.Name,
		Observation: req.Observation,
		OwnerID:     req.CheptelID,
		OwnerType:   entity.Cheptels,
	}

	if err := s.Repository.Create(ctx, &album); err != nil {
		return entity.Album{}, err
	}

	return album, nil
}

// Update updates an album
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.Album, error) {
	if err := req.Validate(); err != nil {
		return entity.Album{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.Album{}, err
	}

	album := entity.Album{
		Model: gorm.Model{
			ID: req.AlbumID,
		},
		OwnerID:   req.CheptelID,
		OwnerType: entity.Cheptels,
	}

	err = s.Repository.Get(ctx, &album)
	if err != nil {
		return entity.Album{}, err
	}

	if req.NewCheptelID != 0 {
		err := s.cheptelManager.OnlyMember(ctx, req.NewCheptelID, req.UserID)
		if err != nil {
			return entity.Album{}, err
		}
	}

	album = entity.Album{
		Model: gorm.Model{
			ID: req.AlbumID,
		},
		OwnerID:   req.NewCheptelID,
		OwnerType: entity.Cheptels,
		Name:      req.NewName,
	}

	err = s.Repository.Update(ctx, &album)
	if err != nil {
		return entity.Album{}, err
	}
	return album, err
}

// Delete deletes an album
func (s *Service) Delete(ctx context.Context, req schema.Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return err
	}

	album := entity.Album{
		Model:     gorm.Model{ID: req.AlbumID},
		OwnerID:   req.CheptelID,
		OwnerType: entity.Cheptels,
	}

	return s.Repository.SoftDelete(ctx, &album)
}
