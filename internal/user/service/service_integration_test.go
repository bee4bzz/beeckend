package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	schema "github.com/gaetanDubuc/beeckend/internal/user"
	"github.com/gaetanDubuc/beeckend/internal/user/repository"
	"github.com/gaetanDubuc/beeckend/internal/user/testutils"
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
	ctx     context.Context
	Service *Service
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	db := db.NewGormForTest(sqlite.Open(dbName))
	suite.Service = NewService(repository.NewGormRepository(db))
	db.Create(&test.ValidUser)
}

// this function executes after all tests executed
func (suite *RepositoryTestSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryTestSuite) TestUpdate() {
	now := time.Now()
	user, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{UserID: test.ValidUser.ID(), Name: "new name"})
	assert.NoError(suite.T(), err)
	testutils.AssertUserUpdated(suite.T(), entity.User{Model: gorm.Model{
		ID: test.ValidUser.ID(),
	},
		Email: test.ValidUser.Email,
		Name:  "new name",
	}, user, now)
}

func (suite *RepositoryTestSuite) TestUpdateNotFound() {
	user, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{UserID: 3, Name: "new name"})
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
	assert.Empty(suite.T(), user)
}

func (suite *RepositoryTestSuite) TestUpdateInvalid() {
	user, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{})
	assert.ErrorIs(suite.T(), err, gorm.ErrMissingWhereClause)
	assert.Empty(suite.T(), user)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
