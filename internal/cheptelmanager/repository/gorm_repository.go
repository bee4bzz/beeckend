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
	return r.db.Model(user).Association(entity.CheptelsKey).Find(cheptel)
}

func (r *GormRepository) Query(ctx context.Context, user *entity.User, cheptels *[]entity.Cheptel) error {
	return r.db.Model(user).Association(entity.CheptelsKey).Find(cheptels)
}

func (r *GormRepository) Create(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	return r.db.Model(user).Association(entity.CheptelsKey).Append(cheptel)
}

func (r *GormRepository) Delete(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	return r.db.Model(user).Association(entity.CheptelsKey).Delete(cheptel)
}

func (r *GormRepository) Update(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	return r.db.Model(user).Association(entity.CheptelsKey).Replace(cheptel)
}
