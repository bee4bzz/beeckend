package testutils

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (r *Repository) Get(ctx context.Context, cheptel *entity.Cheptel) error {
	args := r.Called(*cheptel)
	return args.Error(0)
}

func (r *Repository) QueryByUser(ctx context.Context, user *entity.User, cheptels *[]entity.Cheptel) error {
	args := r.Called(*user, *cheptels)
	return args.Error(0)
}

func (r *Repository) Create(ctx context.Context, cheptel *entity.Cheptel) error {
	args := r.Called(*cheptel)
	return args.Error(0)
}

func (r *Repository) Update(ctx context.Context, cheptel *entity.Cheptel) error {
	args := r.Called(*cheptel)
	return args.Error(0)
}
func (r *Repository) SoftDelete(ctx context.Context, cheptel *entity.Cheptel) error {
	args := r.Called(*cheptel)
	return args.Error(0)
}
