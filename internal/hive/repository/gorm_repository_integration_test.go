package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/hive/testutils"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName = "hive.db"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *gorm.DB
	Repository *GormRepository
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = db.NewGormForTest(sqlite.Open(dbName))
	suite.Repository = NewGormRepository(suite.db)
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

func (suite *RepositoryTestSuite) TestCreate() {
	hive := entity.Hive{
		Model: gorm.Model{
			ID: 100,
		},
		Name:      "new hive",
		CheptelID: test.ValidCheptel.ID,
	}
	hiveCopy := hive
	now := time.Now()
	err := suite.Repository.Create(suite.ctx, &hive)
	assert.NoError(suite.T(), err)
	testutils.AssertHiveCreated(suite.T(), hiveCopy, hive, now)
}

func (suite *RepositoryTestSuite) TestCreateFail() {
	tc := []entity.Hive{
		test.ValidHive,
		{CheptelID: test.ValidHive.CheptelID, Name: test.ValidHive.Name},
	}
	for _, c := range tc {
		err := suite.Repository.Create(suite.ctx, &c)
		assert.ErrorIs(suite.T(), err, gorm.ErrDuplicatedKey)
	}
}

func (suite *RepositoryTestSuite) TestUpdate() {
	now := time.Now()
	user := entity.Hive{Model: gorm.Model{ID: test.ValidHive.ID}, Name: "new name"}
	err := suite.Repository.Update(suite.ctx, &user)
	assert.NoError(suite.T(), err)
	test.ValidHive.Name = "new name"
	testutils.AssertHiveUpdated(suite.T(), test.ValidHive, user, now)
}

func (suite *RepositoryTestSuite) TestGet() {
	user := entity.Hive{Model: gorm.Model{ID: test.ValidHive.ID}}
	err := suite.Repository.Get(suite.ctx, &user)
	assert.NoError(suite.T(), err)
	testutils.AssertHive(suite.T(), test.ValidHive, user)
}

func (suite *RepositoryTestSuite) TestSoftDelete() {
	user := entity.Hive{Model: gorm.Model{ID: 100}}
	suite.Repository.Create(suite.ctx, &user)
	err := suite.Repository.SoftDelete(suite.ctx, &user)
	assert.NoError(suite.T(), err)
	err = suite.Repository.Get(suite.ctx, &user)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryTestSuite) TestQuery() {
	testcases := []struct {
		entity.User
		len int
	}{
		{test.ValidUser, 1},
		{entity.User{Model: gorm.Model{ID: 100}}, 0},
	}

	for _, tc := range testcases {
		hives := []entity.Hive{}
		err := suite.Repository.QueryByUser(suite.ctx, tc.User, &hives)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), hives, tc.len)
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
