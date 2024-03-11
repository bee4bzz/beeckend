package repository

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) Get(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	return r.db.Model(user).Association(entity.CheptelsKey).Find(&cheptel)
}
