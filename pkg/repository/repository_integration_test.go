package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName = "mock.db"
)

type mock struct {
	gorm.Model
	err error
}

func (m mock) Validate() error {
	return m.err
}

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	Repository *Repository[mock]
	db         *gorm.DB
	mock       mock
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	db, _ := gorm.Open(sqlite.Open(dbName))
	db.AutoMigrate(&mock{})
	suite.db = db
	suite.Repository = NewRepository[mock](db)
	suite.mock = mock{Model: gorm.Model{ID: 1}}
}

// this function executes after all tests executed
func (suite *RepositoryTestSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryTestSuite) SetupTest() {
	err := suite.db.Create(&suite.mock).Error
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM mocks")
}

func (suite *RepositoryTestSuite) TestCreate() {
	t := suite.T()
	now := time.Now()
	m := mock{}
	err := suite.Repository.Create(suite.ctx, &m)
	assert.NoError(t, err)
	assert.NotEmpty(t, m.ID)
	assert.GreaterOrEqual(t, m.CreatedAt.Unix(), now.Unix(), "CreatedAt should be greater than now")
	assert.GreaterOrEqual(t, m.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}

func (suite *RepositoryTestSuite) TestInvalidCreate() {
	t := suite.T()
	m := mock{err: fmt.Errorf("error")}
	err := suite.Repository.Create(suite.ctx, &m)
	assert.Error(t, err)
	assert.Empty(t, m.ID)
}

func (suite *RepositoryTestSuite) TestDuplicateCreate() {
	t := suite.T()
	m := mock{}
	err := suite.Repository.Create(suite.ctx, &m)
	assert.NoError(t, err)
	err = suite.Repository.Create(suite.ctx, &m)
	assert.Error(t, err)
}

func (suite *RepositoryTestSuite) TestUpdate() {
	t := suite.T()
	now := time.Now()
	values := &mock{
		Model: gorm.Model{ID: suite.mock.ID},
		err:   fmt.Errorf("error"),
	}
	err := suite.Repository.Update(suite.ctx, values)
	assert.Error(t, err)
	values.err = nil

	err = suite.Repository.Update(suite.ctx, values)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, values.UpdatedAt, now, "UpdatedAt should be greater than now")
	assert.Greater(t, values.UpdatedAt, values.CreatedAt, "UpdatedAt should be greater than CreatedAt")
}

func (suite *RepositoryTestSuite) TestInvalidUpdate() {
	t := suite.T()
	m := mock{}
	err := suite.Repository.Update(suite.ctx, &m)
	assert.Error(t, err)
}

func (suite *RepositoryTestSuite) TestGet() {
	t := suite.T()
	m := mock{Model: gorm.Model{ID: 1}}
	err := suite.Repository.Get(suite.ctx, &m)
	assert.NoError(t, err)
	assert.Equal(t, suite.mock.ID, m.ID)
}

func (suite *RepositoryTestSuite) TestNoRowGet() {
	t := suite.T()
	m := mock{Model: gorm.Model{ID: 1000}}
	err := suite.Repository.Get(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryTestSuite) TestSoftDelete() {
	t := suite.T()
	m := mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.SoftDelete(suite.ctx, &m)
	assert.NoError(t, err)
	err = suite.Repository.Get(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryTestSuite) TestInvalidSoftDelete() {
	t := suite.T()
	m := mock{}
	err := suite.Repository.SoftDelete(suite.ctx, &m)
	assert.Error(t, err)
}

func (suite *RepositoryTestSuite) TestHardDelete() {
	t := suite.T()
	m := mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.HardDelete(suite.ctx, &m)
	assert.NoError(t, err)
	err = suite.Repository.Get(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
