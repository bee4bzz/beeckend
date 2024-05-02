package repository

import (
	"context"

	dbx "github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	*repository.Repository[entity.Cheptel]
}

func NewGormRepository(db *dbx.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.Cheptel](db),
	}
}

func (r *GormRepository) QueryByUser(ctx context.Context, user *entity.User, cheptels *[]entity.Cheptel) error {
	return r.DB().Model(user).Preload(clause.Associations).Association(entity.CheptelsKey).Find(cheptels)
}
