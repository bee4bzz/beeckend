package repository

import (
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
)

type Repository struct {
	*repository.Repository[entity.RefreshToken]
}

// New creates a Repository using `db`.
func New(db *db.DB) *Repository {
	return &Repository{
		Repository: repository.NewRepository[entity.RefreshToken](db),
	}
}
