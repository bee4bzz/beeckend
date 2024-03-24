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

func (r *GormRepository) QueryCheptelAlbumsByUser(ctx context.Context, user *entity.User, albums *[]entity.Album) error {
	err := r.DB().WithContext(ctx).Preload(entity.CheptelsKey + "." + entity.AlbumsKey).Find(user).Error
	if err != nil {
		return err
	}
	for _, cheptel := range user.Cheptels {
		*albums = append(*albums, cheptel.Albums...)

	}
	return err
}

func (r *GormRepository) QueryHiveNoteAlbumsByUser(ctx context.Context, user *entity.User, albums *[]entity.Album) error {
	err := r.DB().WithContext(ctx).Preload(
		entity.CheptelsKey + "." + entity.HivesKey + "." + entity.HiveNotesKey + "." + entity.AlbumsKey,
	).Find(user).Error
	if err != nil {
		return err
	}
	for _, cheptel := range user.Cheptels {
		for _, hive := range cheptel.Hives {
			for _, note := range hive.Notes {
				*albums = append(*albums, note.Albums...)
			}
		}
	}
	return err
}
