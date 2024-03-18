package service

import (
	"context"
	"os"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/errors"
	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/repository"
	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/schema"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName = "cheptelmanager.db"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx      context.Context
	db       *gorm.DB
	Service  *Service
	logger   *log.Logger
	observer *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = db.NewGormForTest(sqlite.Open(dbName))
	logger, obs := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(repository.NewGormRepository(suite.db), suite.logger)
}

// this function executes after all tests executed
func (suite *RepositoryIntegrationSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		suite.T().Fatalf("Error while deleting the database file: %s", err)
	}
}

func (suite *RepositoryIntegrationSuite) SetupTest() {
	db.Seed(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	db.Clean(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TestCreate() {
	err := suite.Service.Create(suite.ctx, schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
		MemberID:  test.ValidUser.ID,
	})

	assert.NoError(suite.T(), err)
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	validCreateReq := schema.Request{
		UserID:    100,
		CheptelID: test.ValidCheptel.ID,
		MemberID:  100,
	}

	invalidCreateReq := validCreateReq.CopyWith(
		schema.Request{
			MemberID: 0,
		},
	)

	testcases := []struct {
		name string
		req  schema.Request
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to create a new association", req: validCreateReq, err: errors.ErrNotMember},
		{name: "Create an assoication without a name should return an error", req: invalidCreateReq},
		{name: "the request should be invalid", req: schema.Request{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Service.Create(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestDelete() {
	req := schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
		MemberID:  test.ValidUser.ID,
	}
	err := suite.Service.SoftDelete(suite.ctx, req)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryIntegrationSuite) TestDeleteFail() {
	validReq := schema.Request{
		UserID:    100,
		CheptelID: test.ValidCheptel.ID,
		MemberID:  test.ValidUser.ID,
	}

	testcases := []struct {
		name string
		req  schema.Request
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to delete another member", req: validReq, err: errors.ErrNotMember},
		{name: "the request should be invalid", req: schema.Request{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Service.SoftDelete(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func TestRepositoryIntegrationSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
