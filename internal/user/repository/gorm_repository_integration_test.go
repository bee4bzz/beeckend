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

type RepositoryIntegrationTestSuite struct {
	suite.Suite
	ctx        context.Context
	Repository *GormRepository
	db         *gorm.DB
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = db.NewGormForTest(sqlite.Open(dbName))
	suite.Repository = NewGormRepository(suite.db)
}

// this function executes after all tests executed
func (suite *RepositoryIntegrationTestSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryIntegrationTestSuite) SetupTest() {
	db.Seed(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationTestSuite) TearDownTest() {
	db.Clean(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationTestSuite) TestCreate() {
	now := time.Now()
	err := suite.Repository.Create(suite.ctx, &entity.User{
		Name:  "new user",
		Email: utils.ValidEmail(),
	})
	assert.NoError(suite.T(), err)
	utils.AssertCreated(suite.T(), test.ValidUser.Model, now)
}

func (suite *RepositoryIntegrationTestSuite) TestUpdate() {
	now := time.Now()
	user := entity.User{Model: gorm.Model{ID: test.ValidUser.ID}, Name: "new name"}
	err := suite.Repository.Update(suite.ctx, &user)
	assert.NoError(suite.T(), err)
	test.ValidUser.Name = "new name"
	testutils.AssertUserUpdated(suite.T(), test.ValidUser, user, now)
}

func (suite *RepositoryIntegrationTestSuite) TestGet() {
	user := entity.User{Model: gorm.Model{ID: test.ValidUser.ID}}
	err := suite.Repository.Get(suite.ctx, &user)
	assert.NoError(suite.T(), err)
	testutils.AssertUser(suite.T(), test.ValidUser, user)
}

func (suite *RepositoryIntegrationTestSuite) TestSoftDelete() {
	err := suite.Repository.SoftDelete(suite.ctx, &test.ValidUser)
	assert.NoError(suite.T(), err)
	err = suite.Repository.Get(suite.ctx, &test.ValidUser)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationTestSuite))
}
