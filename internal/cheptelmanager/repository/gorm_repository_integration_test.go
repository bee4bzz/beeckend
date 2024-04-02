package repository

import (
	"context"
	"os"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/service"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	hivetestutils "github.com/gaetanDubuc/beeckend/internal/hive/testutils"
)

const (
	dbName = "cheptelmanager.db"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx        context.Context
	db         *gorm.DB
	Repository service.Repository
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
	testcases := []struct {
		name    string
		cheptel *entity.Cheptel
	}{
		{
			name: "Valid cheptel",
			cheptel: &entity.Cheptel{
				Model: gorm.Model{ID: test.ValidCheptel.ID},
				Name:  test.ValidCheptel.Name,
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Create(suite.ctx, &test.ValidUser, tc.cheptel)
			assert.NoError(t, err)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestFilterByUser() {
	suite.T().Run("user 1 has multiple cheptels", func(t *testing.T) {
		cheptels := []entity.Cheptel{
			{Model: gorm.Model{ID: test.ValidCheptel.ID}},
		}
		err := suite.Repository.FilterByUserID(suite.ctx, test.ValidUser.ID).Table("cheptels AS Cheptel").Preload(clause.Associations).Find(&cheptels).Error
		assert.NoError(suite.T(), err)
		testutils.AssertCheptels(suite.T(), []entity.Cheptel{test.ValidCheptel, test.ValidCheptel2}, cheptels)
	})

	suite.T().Run("user 2 has one cheptel", func(t *testing.T) {
		cheptels := []entity.Cheptel{
			{Model: gorm.Model{ID: test.ValidCheptel2.ID}},
		}
		err := suite.Repository.FilterByUserID(suite.ctx, test.ValidUser2.ID).Table("cheptels AS Cheptel").Preload(clause.Associations).Find(&cheptels).Error
		assert.NoError(suite.T(), err)
		testutils.AssertCheptels(suite.T(), []entity.Cheptel{test.ValidCheptel}, cheptels)
	})

	suite.T().Run("user 2 has one hive", func(t *testing.T) {
		hives := []entity.Hive{
			{Model: gorm.Model{ID: test.ValidHive.ID}},
		}
		err := suite.Repository.FilterByUserID(suite.ctx, test.ValidUser2.ID).Joins("Cheptel").Preload(entity.HiveNotesKey).Find(&hives).Error
		assert.NoError(suite.T(), err)
		hivetestutils.AssertHives(suite.T(), []entity.Hive{test.ValidHive}, hives)
	})
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
