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
	err := suite.db.Create(&test.ValidHive).Error
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM hives")
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
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID+1, test.ValidUser.ID).Return(nil).Once()

	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID+2, test.ValidUser.ID).Return(test.ErrMock).Once()

	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()
	t := suite.T()

	notFoundUpdateReq := schema.UpdateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidHive.CheptelID + 1,
		HiveID:    test.ValidHive.ID,
	}

	validUpdateReq := notFoundUpdateReq.CopyWith(schema.UpdateRequest{
		CheptelID:    test.ValidHive.CheptelID,
		NewCheptelID: test.ValidHive.CheptelID + 2,
	})

	fmt.Print(validUpdateReq)

	testcases := []struct {
		name string
		req  schema.UpdateRequest
		err  error
	}{
		{name: "hive should be not found", req: notFoundUpdateReq, err: gorm.ErrRecordNotFound},
		{name: "An unknown user of the new cheptel should not be able to update the hive", req: validUpdateReq, err: test.ErrMock},
		{name: "An unknown user of the current cheptel should not be able to update the hive", req: validUpdateReq, err: test.ErrMock},
		{name: "the request should be invalid", req: schema.UpdateRequest{}},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
