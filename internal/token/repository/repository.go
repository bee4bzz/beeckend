package repository

import (
	"context"
	"database/sql"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/errors"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/dbcontext"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
)

type Repository struct {
	DB *dbcontext.DB
}

// NewRepository creates a Repository using `db`.
func NewRepository(db *dbcontext.DB, logger log.Logger) *Repository {
	return &Repository{db}
}

// Get reads the entity with the specified ID from the database.
func (r *Repository) Get(ctx context.Context, token *entity.Token) error {
	err := r.DB.With(ctx).Select().Model(token.UUID, token)
	return err
}

func (r Repository) GetBy(ctx context.Context, filters map[string]interface{}, token *entity.Token) error {
	err := r.DB.With(ctx).Select().Where(dbx.HashExp(filters)).One(token)
	return err
}

func (r Repository) Query(ctx context.Context, tableName entity.TokenType, filters []map[string]interface{}, tokens *[]entity.Token) error {
	var exps []dbx.Expression
	for _, f := range filters {
		exps = append(exps, dbx.HashExp(f))
	}
	var ts []entity.Token
	err := r.DB.With(ctx).Select().Where(dbx.Or(exps...)).From(string(tableName)).All(&ts)
	*tokens = ts
	return err
}

// Create saves a new token record in the database.
func (r *Repository) Create(ctx context.Context, token *entity.Token, exclude ...string) error {
	err := r.Get(ctx, token)
	if err == nil {
		return errors.ErrDuplicate
	} else if err != sql.ErrNoRows {
		return err
	}
	err = r.DB.With(ctx).Model(token).Exclude(exclude...).Insert()
	if err != nil {
		return err
	}
	err = r.Get(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

// Update saves the changes to an entity in the database.
func (r Repository) Update(ctx context.Context, token *entity.Token, exclude ...string) error {
	err := r.DB.With(ctx).Model(token).Exclude(exclude...).Update()
	if err != nil {
		return err
	}
	err = r.Get(ctx, token)
	if err != nil {
		return err
	}
	return err
}

// Delete deletes an entity with the specified ID from the database.
func (r Repository) Delete(ctx context.Context, token *entity.Token) error {
	err := r.Get(ctx, token)
	if err != nil {
		return err
	}
	return r.DB.With(ctx).Model(token).Delete()
}
