package testutils

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (r *Repository) Get(ctx context.Context, album *entity.CheptelAlbum) error {
	args := r.Called(album)
	*album = args.Get(0).(entity.CheptelAlbum)
	return args.Error(1)
}

func (r *Repository) QueryByOwnerIDs(
	ctx context.Context,
	albums *[]entity.CheptelAlbum,
	ownerIDs ...uint,
) error {
	args := r.Called(albums, ownerIDs)
	*albums = args.Get(0).([]entity.CheptelAlbum)
	return args.Error(1)
}

func (r *Repository) Create(ctx context.Context, album *entity.CheptelAlbum) error {
	args := r.Called(album)
	*album = args.Get(0).(entity.CheptelAlbum)
	return args.Error(1)
}

func (r *Repository) Update(ctx context.Context, album *entity.CheptelAlbum) error {
	args := r.Called(album)
	*album = args.Get(0).(entity.CheptelAlbum)
	return args.Error(1)
}
func (r *Repository) SoftDelete(ctx context.Context, album *entity.CheptelAlbum) error {
	args := r.Called(album)
	return args.Error(0)
}
