package cheptel

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	dbcontext "github.com/gaetanDubuc/beeckend/pkg/dbcontext"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
)

type Repository struct {
	*repository.Repository[entity.Cheptel]
}

func NewRepository() *Repository {
	return &Repository{
		Repository: repository.NewRepository[entity.Cheptel](),
	}
}

// Retrieve Hive from the database
func (r *Repository) Get(ctx context.Context, ID uint) (entity.Cheptel, error) {
	db := dbcontext.FromContext(ctx)
	var Cheptel entity.Cheptel
	err := db.Model(&entity.Cheptel{}).Preload(entity.UsersKey).First(&Cheptel, ID).Error
	return Cheptel, err
}
