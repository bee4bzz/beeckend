package repository

import (
	"context"

	"gorm.io/gorm"
)

type Entity interface {
	Validate() error
}

type Repository[T Entity] struct {
	db *gorm.DB
}

func NewRepository[T Entity](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db: db,
	}
}

// Retrieve entity from the database
func (r *Repository[T]) Get(ctx context.Context, entity T) (T, error) {
	var empty T
	err := r.db.WithContext(ctx).Model(&empty).First(&entity).Error
	if err != nil {
		return empty, err
	}
	return entity, err
}

// Create a new entity in the database
func (r *Repository[T]) Create(ctx context.Context, entity T) (T, error) {
	var empty T
	err := entity.Validate()
	if err != nil {
		return empty, err
	}
	err = r.db.WithContext(ctx).Create(&entity).Error
	if err != nil {
		return empty, err
	}
	return entity, err
}

// Update a entity in the database
func (r *Repository[T]) Update(ctx context.Context, entity T) (T, error) {
	var empty T
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&entity).Updates(&entity).Error
		if err != nil {
			return err
		}
		err = tx.Model(&empty).First(&entity).Error
		if err != nil {
			return err
		}
		if err = entity.Validate(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return empty, err
	}
	return entity, nil
}

// Delete a entity from the database
func (r *Repository[T]) SoftDelete(ctx context.Context, entity T) error {
	return r.db.WithContext(ctx).Delete(&entity).Error
}

// Hard delete a entity from the database
func (r *Repository[T]) HardDelete(ctx context.Context, entity T) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity).Error
}
