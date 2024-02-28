package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/user/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Update(ctx context.Context, user entity.User) (entity.User, error)
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
	user := entity.User{
		Model: gorm.Model{
			ID: req.UserID,
		},
		Name: req.Name,
	}
	return s.Repository.Update(ctx, user)
}

