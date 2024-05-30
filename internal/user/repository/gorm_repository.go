package repository

import (
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
)

type GormRepository struct {
	*repository.Repository[entity.User]
}

func NewGormRepository(db *db.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.User](db),
	}
}
