package testutils

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (c *Repository) Get(ctx context.Context, cheptelNote *entity.CheptelNote) error {
	args := c.Called(*cheptelNote)
	return args.Error(0)
}

func (c *Repository) QueryByUser(ctx context.Context, user *entity.User, cheptelNotes *[]entity.CheptelNote) error {
	args := c.Called(*user, *cheptelNotes)
	return args.Error(0)
}
func (c *Repository) Create(ctx context.Context, cheptelNote *entity.CheptelNote) error {
	args := c.Called(*cheptelNote)
	return args.Error(0)
}
func (c *Repository) Update(ctx context.Context, cheptelNote *entity.CheptelNote) error {
	args := c.Called(*cheptelNote)
	return args.Error(0)
}
func (c *Repository) SoftDelete(ctx context.Context, cheptelNote *entity.CheptelNote) error {
	args := c.Called(*cheptelNote)
	return args.Error(0)
}
