package cheptel

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
	*repository.Repository[entity.Cheptel]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db:         db,
		Repository: repository.NewRepository[entity.Cheptel](db),
	}
}

// Retrieve Hive from the database
func (r *Repository) Get(ctx context.Context, ID uint) (entity.Cheptel, error) {
	var Cheptel entity.Cheptel
	err := r.db.Model(&entity.Cheptel{}).Preload(entity.UsersKey).First(&Cheptel, ID).Error
	return Cheptel, err
}
