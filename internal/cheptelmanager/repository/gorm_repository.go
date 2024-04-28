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
	err := r.db.Model(user).Association(entity.CheptelsKey).Find(cheptel)
	if err != nil {
		return err
	} else if cheptel.CreatedAt.IsZero() {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormRepository) Create(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	return r.db.Model(user).Omit(entity.CheptelsKey + ".*").Association(entity.CheptelsKey).Append(cheptel)
}

func (r *GormRepository) Update(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	return r.db.Model(user).Omit(entity.CheptelsKey + ".*").Association(entity.CheptelsKey).Replace(cheptel)
}

func (r *GormRepository) SoftDelete(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	return r.db.Model(user).Association(entity.CheptelsKey).Delete(cheptel)
}
