package cheptel

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName = "cheptel.db"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	Repository *Repository
	Cheptel    entity.Cheptel
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = db.NewContextForTest(sqlite.Open(dbName))
	suite.Repository = NewRepository()
}

// this function executes after all tests executed
func (suite *RepositoryTestSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryTestSuite) TestCreate() {
	now := time.Now()
	cheptel, err := suite.Repository.Create(suite.ctx, test.ValidCheptel)
	assert.NoError(suite.T(), err)
	testutils.AssertCheptelCreated(suite.T(), test.ValidCheptel, cheptel, now)
	suite.Cheptel = cheptel
}

func (suite *RepositoryTestSuite) TestInvalidCreate() {
	cheptel, err := suite.Repository.Create(suite.ctx, entity.Cheptel{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptel)
}

func (suite *RepositoryTestSuite) TestDuplicateCreate() {
	cheptel, err := suite.Repository.Create(suite.ctx, test.ValidCheptel)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptel)
}

func (suite *RepositoryTestSuite) TestUpdate() {
	now := time.Now()
	suite.Cheptel.Name = utils.ValidName()
	cheptel, err := suite.Repository.Update(suite.ctx, suite.Cheptel)
	assert.NoError(suite.T(), err)
	testutils.AssertCheptelUpdated(suite.T(), suite.Cheptel, cheptel, now)
}

func (suite *RepositoryTestSuite) TestInvalidUpdate() {
	cheptel, err := suite.Repository.Update(suite.ctx, entity.Cheptel{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptel)
}

func (suite *RepositoryTestSuite) TestGet() {
	cheptel, err := suite.Repository.Get(suite.ctx, suite.Cheptel.ID)
	assert.NoError(suite.T(), err)
	testutils.AssertCheptel(suite.T(), suite.Cheptel, cheptel)
}

func (suite *RepositoryTestSuite) TestNoRowGet() {
	cheptel, err := suite.Repository.Get(suite.ctx, 1000)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
	assert.Empty(suite.T(), cheptel)
}

func (suite *RepositoryTestSuite) TestSoftDelete() {
	err := suite.Repository.SoftDelete(suite.ctx, suite.Cheptel)
	assert.NoError(suite.T(), err)
	cheptel, err := suite.Repository.Get(suite.ctx, suite.Cheptel.ID)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
	assert.Empty(suite.T(), cheptel)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
