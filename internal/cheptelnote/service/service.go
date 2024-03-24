package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/cheptelnote/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	log "github.com/gaetanDubuc/beeckend/pkg/log"
	"gorm.io/gorm"
)

type CheptelManager interface {
	OnlyMember(ctx context.Context, cheptelID, userID uint) error
}

type Repository interface {
	Get(ctx context.Context, cheptelNote *entity.CheptelNote) error
	QueryByUser(ctx context.Context, user *entity.User, cheptelNotes *[]entity.CheptelNote) error
	Create(ctx context.Context, cheptelNote *entity.CheptelNote) error
	Update(ctx context.Context, cheptelNote *entity.CheptelNote) error
	SoftDelete(ctx context.Context, cheptelNote *entity.CheptelNote) error
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
		logger:         logger.Named("cheptelNote"),
	}
}

func (s *Service) QueryByUser(ctx context.Context, req schema.QueryRequest) ([]entity.CheptelNote, error) {
	logger := s.logger.Named("QueryByUser")
	logger.Debugf("User %v query its cheptelNotes", req.UserID)

	if err := req.Validate(); err != nil {
		logger.Error("the request is invalid")
		return []entity.CheptelNote{}, err
	}

	cheptelNotes := []entity.CheptelNote{}

	err := s.Repository.QueryByUser(ctx, &entity.User{Model: gorm.Model{ID: req.UserID}}, &cheptelNotes)
	if err != nil {
		logger.Error("An error occured when getting the cheptelNotes from the repository")
		return []entity.CheptelNote{}, err
	}

	logger.Debugf("cheptelNotes are retrieved")
	return cheptelNotes, nil
}

// Create creates a cheptelNote
func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (entity.CheptelNote, error) {
	if err := req.Validate(); err != nil {
		return entity.CheptelNote{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.CheptelNote{}, err
	}

	cheptelNote := entity.CheptelNote{
		Model: gorm.Model{
			ID: req.NoteID,
		},
		CheptelID:        req.CheptelID,
		Name:             req.Name,
		TemperatureDay:   req.TemperatureDay,
		TemperatureNight: req.TemperatureNight,
		Weather:          req.Weather,
		Flora:            req.Flora,
		State:            req.State,
		Observation:      req.Observation,
	}

	err = s.Repository.Create(ctx, &cheptelNote)
	if err != nil {
		return entity.CheptelNote{}, err
	}
	return cheptelNote, err
}

// Update updates a cheptelNote
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.CheptelNote, error) {
	if err := req.Validate(); err != nil {
		return entity.CheptelNote{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.CheptelNote{}, err
	}

	if req.NewCheptelID != 0 {
		err := s.cheptelManager.OnlyMember(ctx, req.NewCheptelID, req.UserID)
		if err != nil {
			return entity.CheptelNote{}, err
		}
	}

	cheptelNote := entity.CheptelNote{
		Model: gorm.Model{
			ID: req.NoteID,
		},
		CheptelID: req.NewCheptelID,
		Name:      req.NewName,
	}

	err = s.Repository.Update(ctx, &cheptelNote)
	if err != nil {
		return entity.CheptelNote{}, err
	}
	return cheptelNote, err
}

// Delete deletes a cheptelNote
func (s *Service) SoftDelete(ctx context.Context, req schema.Request) error {
	if err := req.Validate(); err != nil {
		return err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return err
	}

	cheptelNote := entity.CheptelNote{
		Model: gorm.Model{
			ID: req.NoteID,
		},
		CheptelID: req.CheptelID,
	}

	return s.Repository.SoftDelete(ctx, &cheptelNote)
}
