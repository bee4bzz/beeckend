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

func (r *GormRepository) QueryByUser(ctx context.Context, user entity.User, hives *[]entity.Hive) error {
	err := r.DB().WithContext(ctx).Preload(entity.CheptelsKey + "." + entity.HivesKey).Find(&user).Error
	if err != nil {
		return err
	}
	for _, cheptel := range user.Cheptels {
		*hives = append(*hives, cheptel.Hives...)
	}
	return err
}
