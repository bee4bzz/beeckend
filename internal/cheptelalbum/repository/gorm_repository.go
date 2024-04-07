package repository

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type GormRepository struct {
	*repository.Repository[entity.Album]
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.Album](db),
	}
}

func (r *GormRepository) QueryByOwnerIDs(ctx context.Context, IDs []any, albumOwner entity.AlbumOwner, albums *[]entity.Album) error {
	err := r.DB().WithContext(ctx).Where(&entity.Album{
		OwnerType: albumOwner,
	}).Where("owner_id IN (?)", IDs).Find(albums).Error
	return err
}
