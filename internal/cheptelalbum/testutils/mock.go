package testutils

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (r *Repository) Get(ctx context.Context, album *entity.Album) error {
	args := r.Called(album)
	return args.Error(0)
}

func (r *Repository) QueryByOwnerIDs(ctx context.Context, ownerIDs []any, ownerType entity.AlbumType, albums *[]entity.Album) error {
	args := r.Called(ownerIDs, ownerType, albums)
	return args.Error(0)
}

func (r *Repository) Create(ctx context.Context, album *entity.Album) error {
	args := r.Called(album)
	return args.Error(0)
}

func (r *Repository) Update(ctx context.Context, album *entity.Album) error {
	args := r.Called(album)
	return args.Error(0)
}
func (r *Repository) SoftDelete(ctx context.Context, album *entity.Album) error {
	args := r.Called(album)
	return args.Error(0)
}
