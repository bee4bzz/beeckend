package service

import (
	"context"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/hive/schema"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"

	chepteltestutils "github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	hivetestutils "github.com/gaetanDubuc/beeckend/internal/hive/testutils"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx            context.Context
	Service        *Service
	CheptelManager *chepteltestutils.CheptelManager
	Repository     *hivetestutils.Repository
	logger         *log.Logger
	observer       *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.CheptelManager = &chepteltestutils.CheptelManager{}
	suite.Repository = &hivetestutils.Repository{}
	logger, obs := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(suite.Repository, suite.CheptelManager, logger)
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	suite.CheptelManager.AssertExpectations(suite.T())
	suite.Repository.AssertExpectations(suite.T())
}

func (suite *RepositoryTestSuite) TestQueryByUserFail() {
	suite.Repository.On("QueryByUser", entity.User{
		Model: gorm.Model{
			ID: test.ValidUser.ID,
		},
	}, []entity.Hive{}).Return(test.ErrMock).Once()

	hives, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{
		UserID: test.ValidUser.ID,
	})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), hives)
	assert.Equal(suite.T(), 2, suite.observer.Len())
}

func (suite *RepositoryTestSuite) TestDelete() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(nil).Once()
	suite.Repository.On("SoftDelete", entity.Hive{
		Model: gorm.Model{
			ID: test.ValidHive.ID,
		},
		CheptelID: test.ValidCheptel.ID,
	}).Return(nil).Once()

	err := suite.Service.SoftDelete(suite.ctx, schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
		HiveID:    test.ValidHive.ID,
	})
	assert.NoError(suite.T(), err)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
