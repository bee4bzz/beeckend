package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/cheptelnote/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName = "cheptelNote.db"
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
		panic(fmt.Errorf("Error while deleting the database file: %s", err))
	}
}

func (suite *RepositoryIntegrationSuite) SetupTest() {
	db.Seed(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	db.Clean(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TestCreate() {
	cheptelNote := entity.CheptelNote{
		Model: gorm.Model{
			ID: 100,
		},
		Name:      "new cheptel note",
		CheptelID: test.ValidCheptel.ID,
		Weather:   entity.CLOUDY,
		Flora:     "new flora",
	}
	cheptelNoteCopy := cheptelNote
	now := time.Now()
	err := suite.Repository.Create(suite.ctx, &cheptelNote)
	assert.NoError(suite.T(), err)
	testutils.AssertCheptelNoteCreated(suite.T(), cheptelNoteCopy, cheptelNote, now)
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	tc := []entity.CheptelNote{
		test.ValidCheptelNote,
	}
	for _, c := range tc {
		err := suite.Repository.Create(suite.ctx, &c)
		assert.ErrorIs(suite.T(), err, gorm.ErrDuplicatedKey)
	}
}

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	now := time.Now()
	cheptelNote := entity.CheptelNote{Model: gorm.Model{ID: test.ValidCheptelNote.ID}, Name: "new name"}
	err := suite.Repository.Update(suite.ctx, &cheptelNote)
	assert.NoError(suite.T(), err)
	test.ValidCheptelNote.Name = "new name"
	testutils.AssertCheptelNoteUpdated(suite.T(), test.ValidCheptelNote, cheptelNote, now)
}

func (suite *RepositoryIntegrationSuite) TestGet() {
	cheptelNote := entity.CheptelNote{Model: gorm.Model{ID: test.ValidCheptelNote.ID}}
	err := suite.Repository.Get(suite.ctx, &cheptelNote)
	assert.NoError(suite.T(), err)
	testutils.AssertCheptelNote(suite.T(), test.ValidCheptelNote, cheptelNote)
}

func (suite *RepositoryIntegrationSuite) TestSoftDelete() {
	err := suite.Repository.SoftDelete(suite.ctx, &test.ValidCheptelNote)
	assert.NoError(suite.T(), err)
	err = suite.Repository.Get(suite.ctx, &test.ValidCheptelNote)
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
		cheptelNotes := []entity.CheptelNote{}
		err := suite.Repository.QueryByUser(suite.ctx, &tc.User, &cheptelNotes)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), cheptelNotes, tc.len)
	}
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
