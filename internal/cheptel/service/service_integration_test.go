package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/repository"
	"github.com/gaetanDubuc/beeckend/internal/cheptel/schema"
	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	cheptelmngtestutils "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/testutils"
)

const (
	dbName = "cheptel.db"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx            context.Context
	db             *gorm.DB
	Service        *Service
	CheptelManager *cheptelmngtestutils.CheptelManager
	logger         *log.Logger
	observer       *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = db.NewGormForTest(sqlite.Open(dbName))
	suite.CheptelManager = &cheptelmngtestutils.CheptelManager{}
	logger, obs := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(repository.NewGormRepository(suite.db), suite.CheptelManager, suite.logger)
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

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(nil).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID+1, test.ValidUser.ID).Return(nil).Once()
	now := time.Now()

	cheptel, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
		NewName:   "new name"})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptelUpdated(suite.T(), entity.Cheptel{Model: gorm.Model{
		ID: test.ValidCheptel.ID,
	},
		Name: "new name",
	}, cheptel, now)
}

func (suite *RepositoryIntegrationSuite) TestUpdateFail() {
	// cheptel should be not found
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID+1, test.ValidUser.ID).Return(nil).Once()

	// An unknown user of the current cheptel should not be able to update the cheptel
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(test.ErrMock).Once()

	validUpdateReq := schema.UpdateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
	}
	cheptelNotFoundReq := validUpdateReq.CopyWith(schema.UpdateRequest{
		CheptelID: test.ValidCheptel.ID + 1,
	})

	testcases := []struct {
		name string
		req  schema.UpdateRequest
		err  error
	}{
		{name: "cheptel should be not found", req: cheptelNotFoundReq, err: gorm.ErrRecordNotFound},
		{name: "An unknown user of the current cheptel should not be able to update the cheptel", req: validUpdateReq, err: test.ErrMock},
		{name: "the request should be invalid", req: schema.UpdateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			cheptel, err := suite.Service.Update(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, cheptel)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestCreate() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID+1, test.ValidUser.ID).Return(nil).Once()
	now := time.Now()

	cheptel, err := suite.Service.Create(suite.ctx, schema.CreateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID + 1,
		Name:      "new name"})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptelCreated(suite.T(), entity.Cheptel{Model: gorm.Model{
		ID: test.ValidCheptel.ID + 1,
	},
		Name: "new name",
	}, cheptel, now)
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(test.ErrMock).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(nil).Once()

	validCreateReq := schema.CreateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
		Name:      "new name",
	}

	invalidCreateReq := validCreateReq.CopyWith(
		schema.CreateRequest{
			Name: "",
		},
	)

	testcases := []struct {
		name string
		req  schema.CreateRequest
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to update the cheptel", req: validCreateReq, err: test.ErrMock},
		{name: "Create an cheptel without a name should return an error", req: invalidCreateReq},
		{name: "the request should be invalid", req: schema.CreateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			cheptel, err := suite.Service.Create(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, cheptel)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestQueryByUser() {
	cheptels, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{
		UserID: test.ValidUser.ID,
	})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptels(suite.T(), []entity.Cheptel{test.ValidCheptel}, cheptels)

	assert.Equal(suite.T(), 2, suite.observer.Len())
}

func (suite *RepositoryIntegrationSuite) TestQueryByUserFail() {
	cheptels, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptels)
}

func (suite *RepositoryIntegrationSuite) TestDelete() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(nil).Once()
	req := schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
	}
	err := suite.Service.SoftDelete(suite.ctx, req)
	assert.NoError(suite.T(), err)
	err = suite.db.Model(&test.ValidCheptel).First(&entity.Cheptel{}).Error
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestDeleteFail() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(test.ErrMock).Once()

	validReq := schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
	}

	testcases := []struct {
		name string
		req  schema.Request
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to delete the cheptel", req: validReq, err: test.ErrMock},
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
