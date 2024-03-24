package repository

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository struct {
	*repository.Repository[entity.CheptelNote]
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.CheptelNote](db),
	}
}

func (r *GormRepository) QueryByUser(ctx context.Context, user *entity.User, cheptelNotes *[]entity.CheptelNote) error {
	err := r.DB().WithContext(ctx).Preload(entity.CheptelsKey + "." + entity.CheptelNotesKey + "." + clause.Associations).Find(user).Error
	if err != nil {
		return err
	}
	for _, cheptel := range user.Cheptels {
		*cheptelNotes = append(*cheptelNotes, cheptel.Notes...)
	}
	return err
}
