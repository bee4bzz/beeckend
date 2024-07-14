package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/cheptelalbum/testutils"
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
	db         *db.DB
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
	db.Clean(suite.T(), suite.db)
	db.Seed(suite.T(), suite.db, &test.ValidCheptel)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	db.Clean(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TestCreate() {
	album := entity.CheptelAlbum{
		Model: gorm.Model{
			ID: 100,
		},
		Album: entity.Album{
			Name:    "new album",
			OwnerID: test.ValidCheptel.ID,
		},
	}
	albumCopy := album
	now := time.Now()
	err := suite.Repository.Create(suite.ctx, &album)
	assert.NoError(suite.T(), err)
	testutils.AssertAlbumCreated(suite.T(), albumCopy, album, now)
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	tc := []entity.CheptelAlbum{
		test.ValidCheptelAlbum,
	}
	for _, c := range tc {
		err := suite.Repository.Create(suite.ctx, &c)
		assert.ErrorIs(suite.T(), err, gorm.ErrDuplicatedKey)
	}
}

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	now := time.Now()
	album := entity.CheptelAlbum{
		Model: gorm.Model{ID: test.ValidCheptelAlbum.ID},
		Album: entity.Album{Name: "new name"},
	}
	err := suite.Repository.Update(suite.ctx, &album)
	assert.NoError(suite.T(), err)
	test.ValidCheptelAlbum.Name = "new name"
	testutils.AssertAlbumUpdated(suite.T(), test.ValidCheptelAlbum, album, now)
}

func (suite *RepositoryIntegrationSuite) TestGet() {
	album := entity.CheptelAlbum{
		Model: gorm.Model{ID: test.ValidCheptelAlbum.ID},
	}
	err := suite.Repository.Get(suite.ctx, &album)
	assert.NoError(suite.T(), err)
	testutils.AssertAlbum(suite.T(), test.ValidCheptelAlbum, album)
}

func (suite *RepositoryIntegrationSuite) TestSoftDelete() {
	err := suite.Repository.SoftDelete(suite.ctx, &test.ValidCheptelAlbum)
	assert.NoError(suite.T(), err)
	err = suite.Repository.Get(suite.ctx, &test.ValidCheptelAlbum)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
