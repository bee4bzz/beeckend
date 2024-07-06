package service

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/repository"
	"github.com/gaetanDubuc/beeckend/internal/cheptel/schema"
	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/errors"
	cheptelmngrepo "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/repository"
	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/service"
)

const (
	dbName = "cheptel.db"
)

type ServiceIntegrationSuite struct {
	suite.Suite
	ctx                      context.Context
	db                       *db.DB
	Service                  *Service
	CheptelManager           *service.Service
	CheptelManagerRepository *cheptelmngrepo.GormRepository
	logger                   *log.Logger
	buffer                   *bytes.Buffer
	observer                 *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *ServiceIntegrationSuite) SetupSuite() {
	logger, obs, _ := log.NewForTest()
	suite.logger = logger
	suite.observer = obs

	suite.ctx = context.Background()

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

	suite.CheptelManagerRepository = cheptelmngrepo.NewGormRepository(suite.db)
	suite.CheptelManager = &service.Service{Repository: suite.CheptelManagerRepository}

	pool, err := pgxpool.New(suite.ctx, config.DatabaseURL)
	if err != nil {
		panic(err)
	}

	suite.Service = NewService(
		repository.NewGormRepository(suite.db, pool, suite.logger),
		suite.CheptelManager,
		suite.CheptelManagerRepository,
		suite.logger)
}

func (suite *ServiceIntegrationSuite) SetupTest() {
	db.Clean(suite.T(), suite.db)
	db.Seed(suite.T(), suite.db, &test.ValidUser)
}

func (suite *ServiceIntegrationSuite) TearDownTest() {
	suite.T().Log(suite.buffer)
	suite.observer.TakeAll()
}

func (suite *ServiceIntegrationSuite) TestUpdate() {
	now := time.Now()

	cheptel, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
		NewName:   "new name"})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptelUpdated(suite.T(), entity.Cheptel{Model: gorm.Model{
		ID:        test.ValidCheptel.ID,
		CreatedAt: test.ValidCheptel.CreatedAt,
		UpdatedAt: now,
	},
		Name:   "new name",
		Hives:  test.ValidCheptel.Hives,
		Notes:  test.ValidCheptel.Notes,
		Albums: test.ValidCheptel.Albums,
		Users:  test.ValidCheptel.Users,
	}, cheptel, now)
}

func (suite *ServiceIntegrationSuite) TestUpdateFail() {
	validUpdateReq := schema.UpdateRequest{
		UserID:    100,
		CheptelID: test.ValidCheptel.ID,
	}
	cheptelNotFoundReq := validUpdateReq.CopyWith(schema.UpdateRequest{
		CheptelID: 100,
	})

	testcases := []struct {
		name string
		req  schema.UpdateRequest
		err  error
	}{
		{name: "cheptel should be not found", req: cheptelNotFoundReq, err: errors.ErrNotMember},
		{name: "An unknown user of the current cheptel should not be able to update the cheptel", req: validUpdateReq, err: errors.ErrNotMember},
		{name: "the request should be invalid", req: schema.UpdateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			cheptel, err := suite.Service.Update(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, cheptel)
		})
	}
}

func (suite *ServiceIntegrationSuite) TestCreate() {
	now := time.Now()

	cheptel, err := suite.Service.Create(suite.ctx, schema.CreateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: uint(100),
		Name:      "new name"})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptelCreated(suite.T(), entity.Cheptel{Model: gorm.Model{
		ID: 100,
	},
		Name: "new name",
	}, cheptel, now)
}

func (suite *ServiceIntegrationSuite) TestCreateFail() {
	validCreateReq := schema.CreateRequest{
		UserID:    100,
		CheptelID: 100,
		Name:      "new name",
	}

	invalidCreateReq := validCreateReq.CopyWith(
		schema.CreateRequest{
			Name: "",
		},
	)

	testcases := []struct {
		name string
		req  schema.CreateRequest
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to create the cheptel",
			req: validCreateReq, err: gorm.ErrForeignKeyViolated},
		{name: "Create an cheptel without a name should return an error", req: invalidCreateReq},
		{name: "the request should be invalid", req: schema.CreateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			cheptel, err := suite.Service.Create(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, cheptel)
		})
	}
}

func (suite *ServiceIntegrationSuite) TestQueryByUser() {
	cheptels, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{
		UserID: test.ValidUser.ID,
	})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptels(suite.T(), []entity.Cheptel{test.ValidCheptel, test.ValidCheptel2}, cheptels)

	assert.Equal(suite.T(), 2, suite.observer.Len())
}

func (suite *ServiceIntegrationSuite) TestQueryByUserFail() {
	cheptels, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptels)
}

func (suite *ServiceIntegrationSuite) TestDelete() {
	req := schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
	}
	err := suite.Service.SoftDelete(suite.ctx, req)
	assert.NoError(suite.T(), err)
	err = suite.db.DB().Model(&entity.Cheptel{}).First(&test.ValidCheptel).Error
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *ServiceIntegrationSuite) TestDeleteFail() {
	validReq := schema.Request{
		UserID:    100,
		CheptelID: test.ValidCheptel.ID,
	}

	testcases := []struct {
		name string
		req  schema.Request
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to delete the cheptel", req: validReq, err: errors.ErrNotMember},
		{name: "the request should be invalid", req: schema.Request{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Service.SoftDelete(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func TestServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceIntegrationSuite))
}
