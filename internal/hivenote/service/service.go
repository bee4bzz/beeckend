package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/hivenote/schema"
	log "github.com/gaetanDubuc/beeckend/pkg/log"
	"gorm.io/gorm"
)

type CheptelManager interface {
	OnlyMember(ctx context.Context, cheptelID, userID uint) error
}

type Repository interface {
	Get(ctx context.Context, hive *entity.HiveNote) error
	QueryByUser(ctx context.Context, user entity.User, hiveNotes *[]entity.HiveNote) error
	Create(ctx context.Context, hive *entity.HiveNote) error
	Update(ctx context.Context, hive *entity.HiveNote) error
	SoftDelete(ctx context.Context, hive *entity.HiveNote) error
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
		logger:         logger.Named("hive"),
	}
}

func (s *Service) QueryByUser(ctx context.Context, req schema.QueryRequest) ([]entity.HiveNote, error) {
	logger := s.logger.Named("QueryByUser")
	logger.Debugf("User %v query its hive's notes", req.UserID)

	if err := req.Validate(); err != nil {
		logger.Error("the request is invalid")
		return []entity.HiveNote{}, err
	}

	hiveNotes := []entity.HiveNote{}

	err := s.Repository.QueryByUser(ctx, entity.User{Model: gorm.Model{ID: req.UserID}}, &hiveNotes)
	if err != nil {
		logger.Error("An error occured when getting the hive's notes from the repository")
		return []entity.HiveNote{}, err
	}

	logger.Debugf("hive's notes are retrieved")
	return hiveNotes, nil
}

// Create creates an hive
func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (entity.HiveNote, error) {
	if err := req.Validate(); err != nil {
		return entity.HiveNote{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.HiveNote{}, err
	}

	hive := entity.HiveNote{
		Model: gorm.Model{
			ID: req.HiveNoteID,
		},
		HiveID: req.HiveID,
		Name:   req.Name,
	}

	err = s.Repository.Create(ctx, &hive)
	if err != nil {
		return entity.HiveNote{}, err
	}
	return hive, err
}

// Update updates an hive
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.HiveNote, error) {
	if err := req.Validate(); err != nil {
		return entity.HiveNote{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.HiveNote{}, err
	}

	if req.NewHiveID != 0 {
		err := s.cheptelManager.OnlyMember(ctx, req.NewCheptelID, req.UserID)
		if err != nil {
			return entity.HiveNote{}, err
		}
	}

	hive := entity.HiveNote{
		Model: gorm.Model{
			ID: req.HiveNoteID,
		},
		HiveID: req.NID,
		Name:   req.NewName,
	}

	err = s.Repository.Update(ctx, &hive)
	if err != nil {
		return entity.HiveNote{}, err
	}
	return hive, err
}

// Delete deletes an hive
func (s *Service) SoftDelete(ctx context.Context, req schema.Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return err
	}

	hive := entity.HiveNote{
		Model: gorm.Model{
			ID: req.HiveNoteID,
		},
		CheptelID: req.CheptelID,
	}

	return s.Repository.SoftDelete(ctx, &hive)
}
