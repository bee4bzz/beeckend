package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/hive/repository"
	"github.com/gaetanDubuc/beeckend/internal/hive/schema"
	"github.com/gaetanDubuc/beeckend/internal/hive/testutils"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	chepteltestutils "github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
)

const (
	dbName = "hive.db"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx            context.Context
	db             *gorm.DB
	Service        *Service
	CheptelManager *chepteltestutils.CheptelManager
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = db.NewGormForTest(sqlite.Open(dbName))
	suite.CheptelManager = &chepteltestutils.CheptelManager{}
	suite.Service = NewService(repository.NewGormRepository(suite.db), suite.CheptelManager)
}

// this function executes after all tests executed
func (suite *RepositoryTestSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryTestSuite) SetupTest() {
	db.Seed(suite.T(), suite.db)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	db.Clean(suite.T(), suite.db)
}

func (suite *RepositoryTestSuite) TestUpdate() {
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID+1, test.ValidUser.ID).Return(nil).Once()
	now := time.Now()

	hive, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{
		UserID:       test.ValidUser.ID,
		CheptelID:    test.ValidHive.CheptelID,
		HiveID:       test.ValidHive.ID,
		NewCheptelID: test.ValidHive.CheptelID + 1,
		NewName:      "new name"})

	assert.NoError(suite.T(), err)
	testutils.AssertHiveUpdated(suite.T(), entity.Hive{Model: gorm.Model{
		ID: test.ValidHive.ID,
	},
		CheptelID: 10,
		Name:      "new name",
	}, hive, now)
}

func (suite *RepositoryTestSuite) TestUpdateFail() {
	// hive should be not found
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()

	// An unknown user of the new cheptel should not be able to update the hive
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID+1, test.ValidUser.ID).Return(test.ErrMock).Once()

	// An unknown user of the current cheptel should not be able to update the hive
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()

	validUpdateReq := schema.UpdateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidHive.CheptelID,
		HiveID:    test.ValidHive.ID,
	}
	hiveNotFoundReq := validUpdateReq.CopyWith(schema.UpdateRequest{
		HiveID: test.ValidHive.ID + 1,
	})

	newCheptelReq := validUpdateReq.CopyWith(schema.UpdateRequest{
		NewCheptelID: test.ValidHive.CheptelID + 1,
	})

	testcases := []struct {
		name string
		req  schema.UpdateRequest
		err  error
	}{
		{name: "hive should be not found", req: hiveNotFoundReq, err: gorm.ErrRecordNotFound},
		{name: "An unknown user of the new cheptel should not be able to update the hive", req: newCheptelReq, err: test.ErrMock},
		{name: "An unknown user of the current cheptel should not be able to update the hive", req: validUpdateReq, err: test.ErrMock},
		{name: "the request should be invalid", req: schema.UpdateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			hive, err := suite.Service.Update(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, hive)
		})
	}
}

func (suite *RepositoryTestSuite) TestCreate() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(nil).Once()
	now := time.Now()

	hive, err := suite.Service.Create(suite.ctx, schema.CreateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
		Name:      "new name"})

	assert.NoError(suite.T(), err)
	testutils.AssertHiveCreated(suite.T(), entity.Hive{Model: gorm.Model{
		ID: test.ValidHive.ID + 1,
	},
		CheptelID: test.ValidCheptel.ID,
		Name:      "new name",
	}, hive, now)
}

func (suite *RepositoryTestSuite) TestCreateFail() {
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()

	validCreateReq := schema.CreateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidHive.CheptelID,
		HiveID:    test.ValidHive.ID,
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
		{name: "An unknown user of the current cheptel should not be able to update the hive", req: validCreateReq, err: test.ErrMock},
		{name: "Create an hive without a name should return an error", req: invalidCreateReq},
		{name: "the request should be invalid", req: schema.CreateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			hive, err := suite.Service.Create(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, hive)
		})
	}
}

func (suite *RepositoryTestSuite) TestQueryByUser() {
	hives, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{
		UserID: test.ValidUser.ID,
	})

	assert.NoError(suite.T(), err)
	testutils.AssertHives(suite.T(), []entity.Hive{test.ValidHive}, hives)
}

func (suite *RepositoryTestSuite) TestDelete() {
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Twice()
	req := schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidHive.CheptelID,
		HiveID:    test.ValidHive.ID,
	}
	err := suite.Service.SoftDelete(suite.ctx, req)
	assert.NoError(suite.T(), err)
	_, err = suite.Service.Get(suite.ctx, req)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryTestSuite) TestDeleteFail() {
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()

	validReq := schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidHive.CheptelID,
		HiveID:    test.ValidHive.ID,
	}

	testcases := []struct {
		name string
		req  schema.Request
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to update the hive", req: validReq, err: test.ErrMock},
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

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
