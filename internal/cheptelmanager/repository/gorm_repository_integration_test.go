package repository

import (
	"bytes"
	"context"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbName = "cheptelmanager.db"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx        context.Context
	db         *db.DB
	Repository *GormRepository
	buffer     *bytes.Buffer
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, _, b := log.NewForTest()
	suite.buffer = b

	config, err := utils.LoadConfig("../../../")
	if err != nil {
		logger.Fatal("cannot load config:", err)
		panic(err)
	}
	suite.db = db.NewGormWithMigrate(
		postgres.Open(config.DBSource),
		"file://../../../migrations",
		config.DatabaseURL,
		logger)

	suite.Repository = NewGormRepository(suite.db)
}

func (suite *RepositoryIntegrationSuite) SetupTest() {
	db.Clean(suite.T(), suite.db)
	db.Seed(suite.T(), suite.db, &test.ValidUser)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	suite.T().Log(suite.buffer)
	suite.buffer.Reset()
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

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
