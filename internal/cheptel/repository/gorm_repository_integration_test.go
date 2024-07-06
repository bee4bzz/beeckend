package repository

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dbName = "cheptel.db"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx        context.Context
	db         *db.DB
	Repository *GormRepository[*pgxpool.Conn, *pgx.Conn]
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

	pool, err := pgxpool.New(suite.ctx, config.DatabaseURL)
	if err != nil {
		logger.Error("Unable to connect to database:", err)
		panic(err)
	}

	suite.Repository = NewGormRepository(suite.db, pool, logger)
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
	err := suite.Repository.SoftDelete(suite.ctx, &test.ValidCheptel)
	assert.NoError(suite.T(), err)
	err = suite.Repository.Get(suite.ctx, &test.ValidCheptel)
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestQueryByUser() {
	testcases := []struct {
		entity.User
		len int
	}{
		{test.ValidUser, 2},
		{entity.User{Model: gorm.Model{ID: 100}}, 0},
	}

	for _, tc := range testcases {
		cheptels := []entity.Cheptel{}
		err := suite.Repository.QueryByUser(suite.ctx, &tc.User, &cheptels)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), cheptels, tc.len)
	}
}

func (suite *RepositoryIntegrationSuite) Test_A_User_Can_Subscribe_To_Cheptels_Modifications() {
	cheptels := make(chan *[]entity.Cheptel)

	ctx, cancel := context.WithCancel(suite.ctx)
	err := suite.Repository.Subscribe(ctx, &test.ValidUser, cheptels)
	assert.NoError(suite.T(), err)
	testutils.AssertCheptels(suite.T(), test.ValidUser.Cheptels, *<-cheptels)

	err = suite.Repository.SoftDelete(ctx, &entity.Cheptel{Model: gorm.Model{ID: test.ValidUser.Cheptels[0].ID}})
	assert.NoError(suite.T(), err)

	testutils.AssertCheptels(suite.T(), test.ValidUser.Cheptels[1:], *<-cheptels)

	cancel()

	v, ok := <-cheptels
	assert.Nil(suite.T(), v)
	assert.False(suite.T(), ok)
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
