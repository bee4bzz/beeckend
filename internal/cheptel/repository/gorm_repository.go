package repository

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	*repository.Repository[entity.Cheptel]
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.Cheptel](db),
	}
}

func (r *GormRepository) QueryByUser(ctx context.Context, user entity.User, cheptels *[]entity.Cheptel) error {
	err := r.DB().WithContext(ctx).Preload(entity.CheptelsKey + "." + clause.Associations).Find(&user).Error
	if err != nil {
		return err
	}
	*cheptels = user.Cheptels
	return err
}
