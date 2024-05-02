package repository

import (
	"context"

	dbx "github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	*repository.Repository[entity.Hive]
}

func NewGormRepository(db *dbx.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.Hive](db),
	}
}

func (r *GormRepository) QueryByUser(ctx context.Context, user *entity.User, hives *[]entity.Hive) error {
	err := r.With(ctx).Joins(
		"INNER JOIN user_cheptels ON user_cheptels.\"user_id\" = ? AND user_cheptels.\"cheptel_id\" = hives.\"cheptel_id\"", user.ID).
		Preload(clause.Associations).Find(hives).Error

	return err
}
