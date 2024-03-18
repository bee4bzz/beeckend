package service

import (
	"context"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"

	cheptelmngtestutils "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/testutils"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	Service    *Service
	Repository *cheptelmngtestutils.Repository
	logger     *log.Logger
	observer   *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.Repository = &cheptelmngtestutils.Repository{}
	logger, obs := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(suite.Repository, logger)
}

func (suite *RepositoryTestSuite) TestOnlyMember() {
	suite.Repository.On("Get", entity.User{Model: gorm.Model{ID: test.ValidUser.ID}}, entity.Cheptel{Model: gorm.Model{ID: test.ValidCheptel.ID}}).Return(nil)
	err := suite.Service.OnlyMember(suite.ctx, test.ValidCheptel.ID, test.ValidUser.ID)
	assert.NoError(suite.T(), err)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
