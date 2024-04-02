package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/hivenote/schema"
	log "github.com/gaetanDubuc/beeckend/pkg/log"
	"gorm.io/gorm"
)

type HiveRepository interface {
	Get(ctx context.Context, hive *entity.Hive) error
}

type CheptelManager interface {
	OnlyMember(ctx context.Context, cheptelID, userID uint) error
}

type Repository interface {
	Get(ctx context.Context, hivenote *entity.HiveNote) error
	QueryByUser(ctx context.Context, user *entity.User, hiveNotes *[]entity.HiveNote) error
	Create(ctx context.Context, hivenote *entity.HiveNote) error
	Update(ctx context.Context, hivenote *entity.HiveNote) error
	SoftDelete(ctx context.Context, hivenote *entity.HiveNote) error
}

type Service struct {
	Repository
	hiveRepository HiveRepository
	cheptelManager CheptelManager
	logger         *log.Logger
}

func NewService(repository Repository, hiveRepository HiveRepository, cheptelManager CheptelManager, logger *log.Logger) *Service {
	return &Service{
		Repository:     repository,
		hiveRepository: hiveRepository,
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

	err := s.Repository.QueryByUser(ctx, &entity.User{Model: gorm.Model{ID: req.UserID}}, &hiveNotes)
	if err != nil {
		logger.Error("An error occured when getting the hive's notes from the repository")
		return []entity.HiveNote{}, err
	}

	logger.Debugf("hive's notes are retrieved")
	return hiveNotes, nil
}

// Create creates a hive note
func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (entity.HiveNote, error) {
	if err := req.Validate(); err != nil {
		return entity.HiveNote{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.HiveNote{}, err
	}

	// check if the hive exists in the cheptel.
	err = s.hiveRepository.Get(ctx, &entity.Hive{Model: gorm.Model{ID: req.HiveID}, CheptelID: req.CheptelID})
	if err != nil {
		return entity.HiveNote{}, err
	}

	hiveNote := entity.HiveNote{
		Model: gorm.Model{
			ID: req.HiveNoteID,
		},
		HiveID:      req.HiveID,
		Name:        req.Name,
		NBRisers:    req.NBRisers,
		Operation:   req.Operation,
		Observation: req.Observation,
	}

	err = s.Repository.Create(ctx, &hiveNote)
	if err != nil {
		return entity.HiveNote{}, err
	}
	return hiveNote, err
}

// Update updates a hive note
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.HiveNote, error) {
	logger := s.logger.Named("Update")
	logger.Debugf("User %v update the hive note", req.UserID)

	if err := req.Validate(); err != nil {
		return entity.HiveNote{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.HiveNote{}, err
	}

	// check if the hive exists in the cheptel.
	err = s.hiveRepository.Get(ctx,
		&entity.Hive{
			Model:     gorm.Model{ID: req.HiveID},
			CheptelID: req.CheptelID})
	if err != nil {
		return entity.HiveNote{}, err
	}

	// check if the hive note exists in the hive.
	err = s.Repository.Get(
		ctx,
		&entity.HiveNote{Model: gorm.Model{ID: req.HiveNoteID}, HiveID: req.HiveID},
	)
	if err != nil {
		return entity.HiveNote{}, err
	}

	hiveNote := entity.HiveNote{
		Model: gorm.Model{
			ID: req.HiveNoteID,
		},
		Name:        req.NewName,
		NBRisers:    req.NewNBRisers,
		Operation:   req.NewOperation,
		Observation: req.NewObservation,
	}

	err = s.Repository.Update(ctx, &hiveNote)
	if err != nil {
		return entity.HiveNote{}, err
	}
	return hiveNote, err
}

// Delete deletes a hive note
func (s *Service) SoftDelete(ctx context.Context, req schema.Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return err
	}

	// check if the hive exists in the cheptel.
	err = s.hiveRepository.Get(ctx, &entity.Hive{Model: gorm.Model{ID: req.HiveID}, CheptelID: req.CheptelID})
	if err != nil {
		return err
	}

	hiveNote := entity.HiveNote{
		Model: gorm.Model{
			ID: req.HiveNoteID,
		},
		HiveID: req.HiveID,
	}

	return s.Repository.SoftDelete(ctx, &hiveNote)
}
