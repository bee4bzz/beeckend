package cheptel

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type Repository struct {
	*repository.Repository[entity.Cheptel]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Repository: repository.NewRepository[entity.Cheptel](db),
	}
}

// Retrieve Hive from the database
func (r *Repository) Get(ID uint) (entity.Cheptel, error) {
	var Cheptel entity.Cheptel
	err := r.DB.Model(&entity.Cheptel{}).Preload(entity.UsersKey).First(&Cheptel, ID).Error
	return Cheptel, err
}
