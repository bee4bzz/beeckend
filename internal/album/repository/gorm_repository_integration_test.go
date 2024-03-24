package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/album/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName = "album.db"
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
	album := entity.Album{
		Model: gorm.Model{
			ID: 100,
		},
		Name:      "new album",
		OwnerID:   test.ValidCheptel.ID,
		OwnerType: "cheptels",
	}
	albumCopy := album
	now := time.Now()
	err := suite.Repository.Create(suite.ctx, &album)
	assert.NoError(suite.T(), err)
	testutils.AssertAlbumCreated(suite.T(), albumCopy, album, now)
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	tc := []entity.Album{
		test.ValidAlbum,
	}
	for _, c := range tc {
		err := suite.Repository.Create(suite.ctx, &c)
		assert.ErrorIs(suite.T(), err, gorm.ErrDuplicatedKey)
	}
}

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	now := time.Now()
	album := entity.Album{Model: gorm.Model{ID: test.ValidAlbum.ID}, Name: "new name"}
	err := suite.Repository.Update(suite.ctx, &album)
	assert.NoError(suite.T(), err)
	test.ValidAlbum.Name = "new name"
	testutils.AssertAlbumUpdated(suite.T(), test.ValidAlbum, album, now)
}

func (suite *RepositoryIntegrationSuite) TestGet() {
	album := entity.Album{Model: gorm.Model{ID: test.ValidAlbum.ID}}
	err := suite.Repository.Get(suite.ctx, &album)
	assert.NoError(suite.T(), err)
	testutils.AssertAlbum(suite.T(), test.ValidAlbum, album)
}

func (suite *RepositoryIntegrationSuite) TestSoftDelete() {
	err := suite.Repository.SoftDelete(suite.ctx, &test.ValidAlbum)
	assert.NoError(suite.T(), err)
	err = suite.Repository.Get(suite.ctx, &test.ValidAlbum)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestQueryCheptelAlbumsByUser() {
	testcases := []struct {
		entity.User
		len int
	}{
		{test.ValidUser, 1},
		{entity.User{Model: gorm.Model{ID: 100}}, 0},
	}

	for _, tc := range testcases {
		albums := []entity.Album{}
		err := suite.Repository.QueryCheptelAlbumsByUser(suite.ctx, &tc.User, &albums)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), albums, tc.len)
	}
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
