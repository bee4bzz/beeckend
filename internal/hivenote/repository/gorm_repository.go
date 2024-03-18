package repository

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type GormRepository struct {
	*repository.Repository[entity.HiveNote]
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.HiveNote](db),
	}
}

func (r *GormRepository) QueryByUser(ctx context.Context, user *entity.User, hiveNotes *[]entity.HiveNote) error {
	err := r.DB().WithContext(ctx).Preload(entity.CheptelsKey + "." + entity.HivesKey + "." + entity.HiveNotesKey).Find(user).Error
	if err != nil {
		return err
	}
	for _, cheptel := range user.Cheptels {
		for _, hive := range cheptel.Hives {
			*hiveNotes = append(*hiveNotes, hive.Notes...)
		}
	}
	return err
}
