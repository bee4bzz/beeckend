package testutils

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/stretchr/testify/mock"
)

type Repository struct {
	mock.Mock
}

func (c *Repository) Get(ctx context.Context, hive *entity.Hive) error {
	args := c.Called(*hive)
	return args.Error(0)
}

func (c *Repository) QueryByUser(ctx context.Context, user entity.User, hives *[]entity.Hive) error {
	args := c.Called(user, *hives)
	return args.Error(0)
}
func (c *Repository) Create(ctx context.Context, hive *entity.Hive) error {
	args := c.Called(*hive)
	return args.Error(0)
}
func (c *Repository) Update(ctx context.Context, hive *entity.Hive) error {
	args := c.Called(*hive)
	return args.Error(0)
}
func (c *Repository) SoftDelete(ctx context.Context, hive *entity.Hive) error {
	args := c.Called(*hive)
	return args.Error(0)
}
