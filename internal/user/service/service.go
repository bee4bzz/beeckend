package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	schema "github.com/gaetanDubuc/beeckend/internal/user"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type Repository interface {
	Update(ctx context.Context, user repository.Entity) error
}

type Service struct {
	Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		Repository: repository,
	}
}

// Update updates a user
//
// userID is the ID of the user that's updating
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.User, error) {
	user := &entity.User{
		Model: gorm.Model{
			ID: req.UserID,
		},
		Name: req.Name,
	}
	err := s.Repository.Update(ctx, user)
	if err != nil {
		return entity.User{}, err
	}
	return *user, err
}
