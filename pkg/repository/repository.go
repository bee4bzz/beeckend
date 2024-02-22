package repository

import "gorm.io/gorm"

type Validator interface {
	Validate() error
}

type Repository[T Validator] struct {
	DB *gorm.DB
}

func NewRepository[T Validator](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		DB: db,
	}
}

// Retrieve entity from the database
func (r *Repository[T]) Get(ID uint) (T, error) {
	var entity T
	err := r.DB.Model(&entity).First(&entity, ID).Error
	return entity, err
}

// Create a new entity in the database
func (r *Repository[T]) Create(entity T) (T, error) {
	if err := entity.Validate(); err != nil {
		return entity, err
	}
	err := r.DB.Create(&entity).Error
	if err != nil {
		var zero T
		return zero, err
	}
	return entity, err
}

// Update a entity in the database
func (r *Repository[T]) Update(entity T) (T, error) {
	if err := entity.Validate(); err != nil {
		return entity, err
	}
	err := r.DB.Save(&entity).Error
	return entity, err
}

// Delete a entity from the database
func (r *Repository[T]) SoftDelete(entity T) error {
	return r.DB.Delete(&entity).Error
}

// Hard delete a entity from the database
func (r *Repository[T]) HardDelete(entity T) error {
	return r.DB.Unscoped().Delete(&entity).Error
}
