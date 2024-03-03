package testutils

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type CheptelManager struct {
	mock.Mock
}

func (c *CheptelManager) OnlyMember(ctx context.Context, cheptelID, userID uint) error {
	args := c.Called(cheptelID, userID)
	return args.Error(0)
}
