package service

import (
	"context"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"

	chepteltestutils "github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx            context.Context
	Service        *Service
	CheptelManager *chepteltestutils.CheptelManager
	Repository     *chepteltestutils.Repository
	logger         *log.Logger
	observer       *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.CheptelManager = &chepteltestutils.CheptelManager{}
	suite.Repository = &chepteltestutils.Repository{}
	logger, obs := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(suite.Repository, suite.CheptelManager, logger)
}

func (suite *RepositoryTestSuite) TestQueryByUserFail() {
	suite.Repository.On("QueryByUser", entity.User{
		Model: gorm.Model{
			ID: test.ValidUser.ID,
		},
	}, []entity.Cheptel{}).Return(test.ErrMock).Once()

	cheptels, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{
		UserID: test.ValidUser.ID,
	})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptels)
	assert.Equal(suite.T(), 2, suite.observer.Len())
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
