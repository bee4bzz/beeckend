package hive

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type Repository struct {
	*repository.Repository[entity.Hive]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Repository: repository.NewRepository[entity.Hive](db),
	}
}

// Retrieve Hive from the database
func (r *Repository) Get(ID uint) (entity.Hive, error) {
	var Hive entity.Hive
	err := r.DB.Model(&entity.Hive{}).Preload(entity.HiveNotesKey).First(&Hive, ID).Error
	return Hive, err
}
