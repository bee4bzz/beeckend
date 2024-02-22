package hive

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) *repository.Repository[entity.Hive] {
	return repository.NewRepository[entity.Hive](db)
}
