package service

import "github.com/gaetanDubuc/beeckend/internal/entity"

type Repository interface {
	Create(user entity.User) (entity.User, error)
	Update(user entity.User) (entity.User, error)
	Get(ID uint) (entity.User, error)
	SoftDelete(user entity.User) error
}

type Service struct {
	Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		Repository: repository,
	}
}

func (s *Service) Update(user entity.User) (entity.User, error) {

}
