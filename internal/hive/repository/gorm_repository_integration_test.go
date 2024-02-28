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
	"github.com/gaetanDubuc/beeckend/internal/utils"
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
	Repository *GormRepository
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	db := db.NewGormForTest(sqlite.Open(dbName))
	suite.Repository = NewGormRepository(db)
}

// this function executes after all tests executed
func (suite *RepositoryTestSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryTestSuite) TestCreate() {
	now := time.Now()
	err := suite.Repository.Create(suite.ctx, &test.ValidHive)
	assert.NoError(suite.T(), err)
	utils.AssertCreated(suite.T(), test.ValidHive.Model, now)
}

func (suite *RepositoryTestSuite) TestUpdate() {
	now := time.Now()
	user := entity.Hive{Model: gorm.Model{ID: test.ValidHive.ID()}, Name: "new name"}
	err := suite.Repository.Update(suite.ctx, &user)
	assert.NoError(suite.T(), err)
	test.ValidHive.Name = "new name"
	testutils.AssertHiveUpdated(suite.T(), test.ValidHive, user, now)
}

func (suite *RepositoryTestSuite) TestGet() {
	user := entity.Hive{Model: gorm.Model{ID: test.ValidHive.ID()}}
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

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
