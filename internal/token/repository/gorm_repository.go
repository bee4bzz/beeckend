package repository

import (
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
)

type Repository struct {
	*repository.Repository[entity.Token]
}

// NewRepository creates a Repository using `db`.
func NewRepository(db *db.DB) *Repository {
	return &Repository{
		Repository: repository.NewRepository[entity.Token](db),
	}
}
