package repository

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type GormRepository struct {
	*repository.Repository[entity.User]
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.User](db),
	}
}
