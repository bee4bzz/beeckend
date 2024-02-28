package hive

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
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
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
	Repository *repository.Repository[entity.Hive]
	Hive       entity.Hive
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	db := db.NewGormForTest(sqlite.Open(dbName))
	suite.Repository = NewRepository(db)
}

// this function executes after all tests executed
func (suite *RepositoryTestSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryTestSuite) TestCreate() {
	now := time.Now()
	hive, err := suite.Repository.Create(suite.ctx, test.ValidHive)
	assert.NoError(suite.T(), err)
	testutils.AssertHiveCreated(suite.T(), test.ValidHive, hive, now)
	suite.Hive = hive
}

func (suite *RepositoryTestSuite) TestInvalidCreate() {
	hive, err := suite.Repository.Create(suite.ctx, entity.Hive{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), hive)
}

func (suite *RepositoryTestSuite) TestDuplicateCreate() {
	hive, err := suite.Repository.Create(suite.ctx, test.ValidHive)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), hive)
}

func (suite *RepositoryTestSuite) TestUpdate() {
	now := time.Now()
	suite.Hive.Name = utils.ValidName()
	hive, err := suite.Repository.Update(suite.ctx, suite.Hive)
	assert.NoError(suite.T(), err)
	testutils.AssertHiveUpdated(suite.T(), suite.Hive, hive, now)
}

func (suite *RepositoryTestSuite) TestInvalidUpdate() {
	hive, err := suite.Repository.Update(suite.ctx, entity.Hive{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), hive)
}

func (suite *RepositoryTestSuite) TestGet() {
	hive, err := suite.Repository.Get(suite.ctx, suite.Hive)
	assert.NoError(suite.T(), err)
	testutils.AssertHive(suite.T(), suite.Hive, hive)
}

func (suite *RepositoryTestSuite) TestNoRowGet() {
	hive, err := suite.Repository.Get(suite.ctx, entity.Hive{Model: gorm.Model{ID: 1000}})
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
	assert.Empty(suite.T(), hive)
}

func (suite *RepositoryTestSuite) TestSoftDelete() {
	err := suite.Repository.SoftDelete(suite.ctx, suite.Hive)
	assert.NoError(suite.T(), err)
	hive, err := suite.Repository.Get(suite.ctx, suite.Hive)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
	assert.Empty(suite.T(), hive)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
