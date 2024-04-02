package testutils

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type CheptelManager struct {
	mock.Mock
}

func (c *CheptelManager) OnlyMember(ctx context.Context, cheptelID, userID uint) error {
	args := c.Called(cheptelID, userID)
	return args.Error(0)
}

func (c *CheptelManager) Create(ctx context.Context, userID, cheptelID uint) error {
	args := c.Called(userID, cheptelID)
	return args.Error(0)
}

type Repository struct {
	mock.Mock
}

func (r *Repository) FilterByUserID(ctx context.Context, userID uint) (tx *gorm.DB) {
	args := r.Called(userID)
	return args.Get(0).(*gorm.DB)
}

func (r *Repository) Get(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	args := r.Called(*user, *cheptel)
	return args.Error(0)
}

func (r *Repository) QueryByCheptel(ctx context.Context, cheptel *entity.Cheptel, users *[]entity.User) error {
	args := r.Called(*cheptel, *users)
	return args.Error(0)
}

func (r *Repository) Create(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	args := r.Called(*user, *cheptel)
	return args.Error(0)
}

func (r *Repository) Update(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	args := r.Called(*user, *cheptel)
	return args.Error(0)
}
func (r *Repository) SoftDelete(ctx context.Context, user *entity.User, cheptel *entity.Cheptel) error {
	args := r.Called(*user, *cheptel)
	return args.Error(0)
}
