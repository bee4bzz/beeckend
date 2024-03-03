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
	Update(ctx context.Context, hive *entity.Hive) error
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

func (s *Service) Get(ctx context.Context, req schema.GetRequest) (entity.Hive, error) {
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
	}

	err = s.Repository.Get(ctx, &hive)
	if err != nil {
		return entity.Hive{}, err
	}

	return hive, nil
}

// Update updates an hive
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.Hive, error) {
	if err := req.Validate(); err != nil {
		return entity.Hive{}, err
	}

	hive, err := s.Get(ctx, schema.GetRequest{
		UserID:    req.UserID,
		CheptelID: req.CheptelID,
		HiveID:    req.HiveID,
	})

	if err != nil {
		return entity.Hive{}, err
	}

	if req.NewCheptelID != 0 {
		err := s.cheptelManager.OnlyMember(ctx, req.NewCheptelID, req.UserID)
		if err != nil {
			return entity.Hive{}, err
		}
	}

	hive = entity.Hive{
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
