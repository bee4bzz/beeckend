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

	_, err := s.Get(ctx, schema.GetRequest{
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
func (s *Service) Delete(ctx context.Context, req schema.GetRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	_, err := s.Get(ctx, schema.GetRequest{
		UserID:    req.UserID,
		CheptelID: req.CheptelID,
		HiveID:    req.HiveID,
	})

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
