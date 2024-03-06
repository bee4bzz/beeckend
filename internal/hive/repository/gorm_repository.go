package repository

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type GormRepository struct {
	*repository.Repository[entity.Hive]
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.Hive](db),
	}
}

func (r *GormRepository) QueryByUser(ctx context.Context, user *entity.User) error {
	return r.DB().WithContext(ctx).Preload(entity.HivesKey).Preload(entity.CheptelsKey).Find(user).Error
}
