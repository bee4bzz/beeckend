package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName = "cheptel.db"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx        context.Context
	db         *gorm.DB
	Repository *GormRepository
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = db.NewGormForTest(sqlite.Open(dbName))
	suite.Repository = NewGormRepository(suite.db)
}

// this function executes after all tests executed
func (suite *RepositoryIntegrationSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		suite.T().Errorf("Error while deleting the database file: %s", err)
	}
}

func (suite *RepositoryIntegrationSuite) SetupTest() {
	db.Seed(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	db.Clean(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TestCreate() {
	cheptel := entity.Cheptel{
		Model: gorm.Model{
			ID: 100,
		},
		Name: "new cheptel",
	}
	cheptelCopy := cheptel
	now := time.Now()
	err := suite.Repository.Create(suite.ctx, &cheptel)
	assert.NoError(suite.T(), err)
	testutils.AssertCheptelCreated(suite.T(), cheptelCopy, cheptel, now)
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	tc := []entity.Cheptel{
		test.ValidCheptel,
	}
	for _, c := range tc {
		err := suite.Repository.Create(suite.ctx, &c)
		assert.ErrorIs(suite.T(), err, gorm.ErrDuplicatedKey)
	}
}

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	now := time.Now()
	cheptel := entity.Cheptel{Model: gorm.Model{ID: test.ValidCheptel.ID}, Name: "new name"}
	err := suite.Repository.Update(suite.ctx, &cheptel)
	assert.NoError(suite.T(), err)
	test.ValidCheptel.Name = "new name"
	testutils.AssertCheptelUpdated(suite.T(), test.ValidCheptel, cheptel, now)
}

func (suite *RepositoryIntegrationSuite) TestGet() {
	cheptel := entity.Cheptel{Model: gorm.Model{ID: test.ValidCheptel.ID}}
	err := suite.Repository.Get(suite.ctx, &cheptel)
	assert.NoError(suite.T(), err)
	testutils.AssertCheptel(suite.T(), test.ValidCheptel, cheptel)
}

func (suite *RepositoryIntegrationSuite) TestSoftDelete() {
	cheptel := entity.Cheptel{Model: gorm.Model{ID: 100}}
	suite.Repository.Create(suite.ctx, &cheptel)
	err := suite.Repository.SoftDelete(suite.ctx, &cheptel)
	assert.NoError(suite.T(), err)
	err = suite.Repository.Get(suite.ctx, &cheptel)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestQueryByUser() {
	testcases := []struct {
		entity.User
		len int
	}{
		{test.ValidUser, 1},
		{entity.User{Model: gorm.Model{ID: 100}}, 0},
	}

	for _, tc := range testcases {
		cheptels := []entity.Cheptel{}
		err := suite.Repository.QueryByUser(suite.ctx, tc.User, &cheptels)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), cheptels, tc.len)
	}
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
