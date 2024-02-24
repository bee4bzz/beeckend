package repository

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type GormRepository struct {
	*repository.Repository[entity.User]
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		Repository: repository.NewRepository[entity.User](db),
	}
}

// Retrieve user from the database
func (r *GormRepository) Get(ID uint) (entity.User, error) {
	var user entity.User
	err := r.DB.Model(&entity.User{}).Preload(entity.CheptelsKey).First(&user, ID).Error
	return user, err
}
