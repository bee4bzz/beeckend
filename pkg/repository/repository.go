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
func (r *Repository[T]) Get(ctx context.Context, entity Entity) error {
	return r.db.WithContext(ctx).Model(entity).First(entity).Error
}

// Create a new entity in the database
func (r *Repository[T]) Create(ctx context.Context, entity Entity) error {
	err := entity.Validate()
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update a entity in the database
func (r *Repository[T]) Update(ctx context.Context, entity Entity) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(entity).Updates(entity).Error
		if err != nil {
			return err
		}
		err = tx.Model(entity).First(entity).Error
		if err != nil {
			return err
		}
		if err = entity.Validate(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Delete a entity from the database
func (r *Repository[T]) SoftDelete(ctx context.Context, entity Entity) error {
	return r.db.WithContext(ctx).Model(entity).Delete(entity).Error
}

// Hard delete a entity from the database
func (r *Repository[T]) HardDelete(ctx context.Context, entity Entity) error {
	return r.db.WithContext(ctx).Unscoped().Model(entity).Delete(entity).Error
}
