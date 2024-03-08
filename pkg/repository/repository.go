// The repository package is a generic package using gorm to handle the database
// operations. It is used by the service package to handle the database
// operations.
package repository

import (
	"context"

	"gorm.io/gorm"
)

type Validator interface {
	Validate() error
}

// NewRepository returns a new repository
func NewRepository[T Validator](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db: db,
	}
}

type Repository[T Validator] struct {
	db *gorm.DB
}

// DB returns the gorm database
func (r *Repository[T]) DB() *gorm.DB {
	return r.db
}

// Retrieve entity from the database
//
// This function will retrieve the entity from the database.
// The entity passed as argument will be used as a filter to retrieve the entity.
func (r *Repository[T]) Get(ctx context.Context, entity *T) error {
	var empty T
	err := r.db.WithContext(ctx).Model(&empty).Where(entity).First(&entity).Error
	return err
}

// Create a new entity in the database
//
// This function will validate the entity to make sure it respect the business
// rules. Then it will create the entity in the database.
// The created entity can be retrieved from the entity passed as argument.
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	err := (*entity).Validate()
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Model(entity).Create(entity).Error
	return err
}

// Update a entity in the database
//
// This function will update the entity in the database and then retrieve it
// to make sure the update was successful. Only a partial update is possible
// to optimize the query. Then we valid the entity to make sure it respect
// the business rules.
// A transaction is used to make sure the update and the validation are
// atomic.
// The updated entity can be retrieved from the entity passed as argument.
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

// Soft Delete a entity from the database
//
// The entity is still in the database but is marked as deleted and
// can't be retrieved from the database.
// The entity can't be retrieved through the pointer passed to the function.
func (r *Repository[T]) SoftDelete(ctx context.Context, entity *T) error {
	var empty T
	return r.db.WithContext(ctx).Model(&empty).Where(entity).Delete(&empty).Error
}

// Hard delete a entity from the database
//
// The entity can't be retrieved from the database.
func (r *Repository[T]) HardDelete(ctx context.Context, entity *T) error {
	var empty T
	return r.db.WithContext(ctx).Unscoped().Model(&empty).Where(entity).Delete(&empty).Error
}
