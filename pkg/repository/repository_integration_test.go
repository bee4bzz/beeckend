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
	dbName = "Mock.db"
)

var (
	err = fmt.Errorf("")
)

type MockChild struct {
	gorm.Model
	MockID uint
}

type Mock struct {
	gorm.Model
	err        error
	Attr       int
	MockChilds []MockChild
}

func (m Mock) Validate() error {
	return m.err
}

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx        context.Context
	Repository *Repository[Mock]
	db         *gorm.DB
	Mock       Mock
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()
	db, _ := gorm.Open(sqlite.Open(dbName), &gorm.Config{TranslateError: true})
	err := db.AutoMigrate(&Mock{}, &MockChild{})
	assert.NoError(suite.T(), err)
	suite.db = db
	suite.Repository = NewRepository[Mock](db)
}

// this function executes after all tests executed
func (suite *RepositoryIntegrationSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryIntegrationSuite) SetupTest() {
	err := suite.db.Create(&Mock{Model: gorm.Model{ID: 1}, Attr: 1, MockChilds: []MockChild{{Model: gorm.Model{ID: 1}}}}).Error
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM Mocks")
}

func (suite *RepositoryIntegrationSuite) TestCreate() {
	t := suite.T()
	now := time.Now()
	m := Mock{}
	err := suite.Repository.Create(suite.ctx, &m)
	assert.NoError(t, err)
	assert.NotEmpty(t, m.ID)
	assert.GreaterOrEqual(t, m.CreatedAt.Unix(), now.Unix(), "CreatedAt should be greater than now")
	assert.GreaterOrEqual(t, m.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	testcases := []struct {
		name string
		Mock Mock
		err  error
	}{
		{
			name: "should fail when trying to create a duplicate",
			Mock: Mock{Model: gorm.Model{ID: 1}},
			err:  gorm.ErrDuplicatedKey,
		},
		{
			name: "should fail when the model is invalid",
			Mock: Mock{Model: gorm.Model{ID: 1}, err: err},
			err:  err,
		},
	}
	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Create(suite.ctx, &tc.Mock)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	t := suite.T()
	now := time.Now()
	Mock := &Mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.Update(suite.ctx, Mock)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, Mock.UpdatedAt, now, "UpdatedAt should be greater than now")
	assert.Greater(t, Mock.UpdatedAt, Mock.CreatedAt, "UpdatedAt should be greater than CreatedAt")
	assert.Equal(t, 1, Mock.Attr)
}

func (suite *RepositoryIntegrationSuite) TestUpdateFail() {
	testcases := []struct {
		name string
		Mock Mock
		err  error
	}{
		{
			name: "should fail when model is invalid",
			Mock: Mock{Model: gorm.Model{ID: 1}, err: err},
			err:  err,
		},
		{
			name: "should fail when model is empty",
			Mock: Mock{},
			err:  gorm.ErrMissingWhereClause,
		},
	}
	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Update(suite.ctx, &tc.Mock)
			assert.ErrorIs(t, err, tc.err)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestGet() {
	t := suite.T()
	m := Mock{Model: gorm.Model{ID: 1}}
	err := suite.Repository.Get(suite.ctx, &m)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), m.ID)
	assert.Equal(t, 1, m.Attr)
	assert.Len(t, m.MockChilds, 1)
	assert.NotEmpty(t, m.CreatedAt)
	assert.NotEmpty(t, m.UpdatedAt)
}

func (suite *RepositoryIntegrationSuite) TestNoRowGet() {
	t := suite.T()
	m := Mock{Model: gorm.Model{ID: 1000}}
	err := suite.Repository.Get(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestSoftDelete() {
	t := suite.T()
	m := Mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.SoftDelete(suite.ctx, &m)
	assert.NoError(t, err)
	err = suite.Repository.Get(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestInvalidSoftDelete() {
	t := suite.T()
	m := Mock{}
	err := suite.Repository.SoftDelete(suite.ctx, &m)
	assert.ErrorIs(t, err, gorm.ErrMissingWhereClause)
}

func (suite *RepositoryIntegrationSuite) TestHardDelete() {
	t := suite.T()
	m := Mock{
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
