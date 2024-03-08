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

var (
	err = fmt.Errorf("")
)

type mock struct {
	gorm.Model
	err  error
	Attr int
}

func (m mock) Validate() error {
	return m.err
}

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx        context.Context
	Repository *Repository[mock]
	db         *gorm.DB
	mock       mock
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()
	db, _ := gorm.Open(sqlite.Open(dbName), &gorm.Config{TranslateError: true})
	db.AutoMigrate(&mock{})
	suite.db = db
	suite.Repository = NewRepository[mock](db)
}

// this function executes after all tests executed
func (suite *RepositoryIntegrationSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryIntegrationSuite) SetupTest() {
	err := suite.db.Create(&mock{Model: gorm.Model{ID: 1}, Attr: 1}).Error
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM mocks")
}

func (suite *RepositoryIntegrationSuite) TestCreate() {
	t := suite.T()
	now := time.Now()
	m := mock{}
	err := suite.Repository.Create(suite.ctx, &m)
	assert.NoError(t, err)
	assert.NotEmpty(t, m.ID)
	assert.GreaterOrEqual(t, m.CreatedAt.Unix(), now.Unix(), "CreatedAt should be greater than now")
	assert.GreaterOrEqual(t, m.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	testcases := []struct {
		name string
		mock mock
		err  error
	}{
		{
			name: "should fail when trying to create a duplicate",
			mock: mock{Model: gorm.Model{ID: 1}},
			err:  gorm.ErrDuplicatedKey,
		},
		{
			name: "should fail when the model is invalid",
			mock: mock{Model: gorm.Model{ID: 1}, err: err},
			err:  err,
		},
	}
	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Create(suite.ctx, &tc.mock)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	t := suite.T()
	now := time.Now()
	mock := &mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.Update(suite.ctx, mock)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, mock.UpdatedAt, now, "UpdatedAt should be greater than now")
	assert.Greater(t, mock.UpdatedAt, mock.CreatedAt, "UpdatedAt should be greater than CreatedAt")
	assert.Equal(t, 1, mock.Attr)
}

func (suite *RepositoryIntegrationSuite) TestUpdateFail() {
	testcases := []struct {
		name string
		mock mock
		err  error
	}{
		{
			name: "should fail when model is invalid",
			mock: mock{Model: gorm.Model{ID: 1}, err: err},
			err:  err,
		},
		{
			name: "should fail when model is empty",
			mock: mock{},
			err:  gorm.ErrMissingWhereClause,
		},
	}
	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Update(suite.ctx, &tc.mock)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestGet() {
	t := suite.T()
	m := mock{Model: gorm.Model{ID: 1}}
	err := suite.Repository.Get(suite.ctx, &m)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), m.ID)
	assert.NotEmpty(t, m.CreatedAt)
	assert.NotEmpty(t, m.UpdatedAt)
}

func (suite *RepositoryIntegrationSuite) TestNoRowGet() {
	t := suite.T()
	m := mock{Model: gorm.Model{ID: 1000}}
	err := suite.Repository.Get(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestSoftDelete() {
	t := suite.T()
	m := mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.SoftDelete(suite.ctx, &m)
	assert.NoError(t, err)
	err = suite.Repository.Get(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestInvalidSoftDelete() {
	t := suite.T()
	m := mock{}
	err := suite.Repository.SoftDelete(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrMissingWhereClause)
}

func (suite *RepositoryIntegrationSuite) TestHardDelete() {
	t := suite.T()
	m := mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.HardDelete(suite.ctx, &m)
	assert.NoError(t, err)
	err = suite.Repository.Get(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestRepositoryIntegrationSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
