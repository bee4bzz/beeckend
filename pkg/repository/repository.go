package repository

import (
	"context"

	"gorm.io/gorm"
)

type Validator interface {
	Validate() error
}

type Repository[T Validator] struct {
	db *gorm.DB
}

func NewRepository[T Validator](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db: db,
	}
}

func (r *Repository[T]) DB() *gorm.DB {
	return r.db
}

// Retrieve entity from the database
func (r *Repository[T]) Get(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Where(entity).Model(entity).First(entity).Error
	return err
}

// Create a new entity in the database
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	err := (*entity).Validate()
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Model(entity).Create(entity).Error
	return err
}

// Update a entity in the database
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(entity).Updates(entity).Error
		if err != nil {
			return err
		}
		err = tx.Model(entity).Where(entity).First(entity).Error
		if err != nil {
			return err
		}
		if err = (*entity).Validate(); err != nil {
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
func (r *Repository[T]) SoftDelete(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Model(entity).Where(entity).Delete(entity).Error
}

// Hard delete a entity from the database
func (r *Repository[T]) HardDelete(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Unscoped().Model(entity).Where(entity).Delete(entity).Error
}
