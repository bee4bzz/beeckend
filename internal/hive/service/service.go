package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/hive/schema"
	"gorm.io/gorm"
)

type CheptelManager interface {
	OnlyMember(ctx context.Context, cheptelID, userID uint) error
}

type Repository interface {
	Get(ctx context.Context, hive *entity.Hive) error
	QueryByUser(ctx context.Context, user entity.User, hives *[]entity.Hive) error
	Create(ctx context.Context, hive *entity.Hive) error
	Update(ctx context.Context, hive *entity.Hive) error
	SoftDelete(ctx context.Context, hive *entity.Hive) error
}

type Service struct {
	Repository
	cheptelManager CheptelManager
}

func NewService(repository Repository, cheptelManager CheptelManager) *Service {
	return &Service{
		Repository:     repository,
		cheptelManager: cheptelManager,
	}
}

func (s *Service) QueryByUser(ctx context.Context, req schema.QueryRequest) ([]entity.Hive, error) {
	if err := req.Validate(); err != nil {
		return []entity.Hive{}, err
	}

	hives := []entity.Hive{}

	err := s.Repository.QueryByUser(ctx, entity.User{Model: gorm.Model{ID: req.UserID}}, &hives)
	if err != nil {
		return []entity.Hive{}, err
	}

	return hives, nil
}

// Create creates an hive
func (s *Service) Create(ctx context.Context, req schema.CreateRequest) (entity.Hive, error) {
	if err := req.Validate(); err != nil {
		return entity.Hive{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.Hive{}, err
	}

	hive := entity.Hive{
		Model: gorm.Model{
			ID: req.HiveID,
		},
		CheptelID: req.CheptelID,
		Name:      req.Name,
	}

	err = s.Repository.Create(ctx, &hive)
	if err != nil {
		return entity.Hive{}, err
	}
	return hive, err
}

// Update updates an hive
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.Hive, error) {
	if err := req.Validate(); err != nil {
		return entity.Hive{}, err
	}

	err := s.cheptelManager.OnlyMember(ctx, req.CheptelID, req.UserID)
	if err != nil {
		return entity.Hive{}, err
	}

	if req.NewCheptelID != 0 {
		err := s.cheptelManager.OnlyMember(ctx, req.NewCheptelID, req.UserID)
		if err != nil {
			return entity.Hive{}, err
		}
	}

	hive := entity.Hive{
		Model: gorm.Model{
			ID: req.HiveID,
		},
		CheptelID: req.NewCheptelID,
		Name:      req.NewName,
	}

	err = s.Repository.Update(ctx, &hive)
	if err != nil {
		return entity.Hive{}, err
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

	hive := entity.Hive{
		Model: gorm.Model{
			ID: req.HiveID,
		},
		CheptelID: req.CheptelID,
	}

	return s.Repository.SoftDelete(ctx, &hive)
}
