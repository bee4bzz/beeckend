package service

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/user/schema"
	"gorm.io/gorm"
)

type Repository interface {
	Update(ctx context.Context, user *entity.User) error
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
func (s *Service) Update(ctx context.Context, req schema.UpdateRequest) (entity.User, error) {
	user := entity.User{
		Model: gorm.Model{
			ID: req.UserID,
		},
		Name: req.Name,
	}
	err := s.Repository.Update(ctx, &user)
	if err != nil {
		return entity.User{}, err
	}
	return user, err
}
