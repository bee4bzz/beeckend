package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/user/testutils"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName = "user.db"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	Repository *GormRepository
	User       entity.User
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
	user, err := suite.Repository.Create(suite.ctx, test.ValidUser)
	assert.NoError(suite.T(), err)
	testutils.AssertUserCreated(suite.T(), test.ValidUser, user, now)
	suite.User = user
}

func (suite *RepositoryTestSuite) TestInvalidCreate() {
	user, err := suite.Repository.Create(suite.ctx, entity.User{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), user)
}

func (suite *RepositoryTestSuite) TestDuplicateCreate() {
	user, err := suite.Repository.Create(suite.ctx, test.ValidUser)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), user)
}

func (suite *RepositoryTestSuite) TestUpdate() {
	now := time.Now()
	suite.User.Name = utils.ValidName()
	user, err := suite.Repository.Update(suite.ctx, suite.User)
	assert.NoError(suite.T(), err)
	testutils.AssertUserUpdated(suite.T(), suite.User, user, now)
}

func (suite *RepositoryTestSuite) TestInvalidUpdate() {
	user, err := suite.Repository.Update(suite.ctx, entity.User{Model: gorm.Model{ID: suite.User.ID()}, Email: "not an email"})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), user)
}

func (suite *RepositoryTestSuite) TestGet() {
	user, err := suite.Repository.Get(suite.ctx, suite.User)
	assert.NoError(suite.T(), err)
	testutils.AssertUser(suite.T(), suite.User, user)
}

func (suite *RepositoryTestSuite) TestNoRowGet() {
	user, err := suite.Repository.Get(suite.ctx, entity.User{Model: gorm.Model{ID: 10}})
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
	assert.Empty(suite.T(), user)
}

func (suite *RepositoryTestSuite) TestSoftDelete() {
	err := suite.Repository.SoftDelete(suite.ctx, suite.User)
	assert.NoError(suite.T(), err)
	user, err := suite.Repository.Get(suite.ctx, suite.User)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
	assert.Empty(suite.T(), user)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
