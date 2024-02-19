package user

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
	"gorm.io/gorm"
)

type Repository struct {
	*repository.Repository[entity.User]
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Repository: repository.NewRepository[entity.User](db),
	}
}

// Retrieve user from the database
func (r *Repository) Get(ID uint) (entity.User, error) {
	var user entity.User
	err := r.DB.Model(&entity.User{}).Preload(entity.CheptelTableName).First(&user, ID).Error
	return user, err
}
