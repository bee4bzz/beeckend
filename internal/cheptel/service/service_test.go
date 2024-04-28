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
	cheptelmngtestutils "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/testutils"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx                      context.Context
	Service                  *Service
	CheptelManager           *cheptelmngtestutils.CheptelManager
	CheptelManagerRepository *cheptelmngtestutils.Repository
	Repository               *chepteltestutils.Repository
	logger                   *log.Logger
	observer                 *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.CheptelManager = &cheptelmngtestutils.CheptelManager{}
	suite.CheptelManagerRepository = &cheptelmngtestutils.Repository{}
	suite.Repository = &chepteltestutils.Repository{}
	logger, obs := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(suite.Repository, suite.CheptelManager, suite.CheptelManagerRepository, logger)
}

func (suite *RepositoryTestSuite) TestQueryByUserFail() {
	suite.Repository.On("QueryByUser", entity.User{
		Model: gorm.Model{
			ID: test.ValidUser.ID,
		},
	}, []entity.Cheptel{}).Return([]entity.Cheptel{}, test.ErrMock).Once()

	cheptels, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{
		UserID: test.ValidUser.ID,
	})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptels)
	assert.Equal(suite.T(), 2, suite.observer.Len())
}

func (suite *RepositoryTestSuite) TestUpdateFail() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(nil).Once()
	suite.Repository.On("Update", entity.Cheptel{Model: gorm.Model{ID: test.ValidCheptel.ID}}).Return(test.ErrMock).Once()

	cheptel, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
	})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptel)
	assert.Equal(suite.T(), 2, suite.observer.Len())
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
