package repository

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type GormRepository struct {
	*repository.Repository[entity.CheptelAlbum]
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		repository.NewRepository[entity.CheptelAlbum](db),
	}
}

func (r *GormRepository) QueryByOwnerIDs(ctx context.Context, albums *[]entity.CheptelAlbum, OwnerIDs ...uint) error {
	err := r.DB().WithContext(ctx).Where("owner_id IN (?)", OwnerIDs).Find(albums).Error
	return err
}
