package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/config"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/crypto"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/test/v2"
	errors "gitlab.com/fogo-dev/infrastructure/web-api/internal/token"
	repository "gitlab.com/fogo-dev/infrastructure/web-api/internal/token/repository"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/token/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/token/testutils"
	db "gitlab.com/fogo-dev/infrastructure/web-api/internal/utils"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/dbcontext"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
	"go.uber.org/zap/zaptest/observer"
)

type ServiceIntegrationSuite struct {
	suite.Suite
	ctx        context.Context
	db         *dbcontext.DB
	Service    *Service
	repository *repository.Repository
	logger     log.Logger
	observer   *observer.ObservedLogs
}

// this function executes before the test suite begins execution.
func (suite *ServiceIntegrationSuite) SetupSuite() {
	logger, obs, _ := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.ctx = context.Background()

	cfg, err := config.GetConfig(logger)
	if err != nil {
		panic(err)
	}

	suite.db = db.NewDB(cfg.DSN, "file://../../../migrations", suite.logger)

	suite.repository = repository.NewRepository(suite.db, suite.logger)
	suite.Service = NewService(suite.repository, crypto.NewHasher(32), suite.logger)
}

func (suite *ServiceIntegrationSuite) SetupTest() {
	test.Clean(suite.T(), suite.db)
	test.Seed(suite.T(), suite.db, []any{&test.HomeWithManyUsersManySpaces}, []any{&test.Space1UnderHomeWithManyUsersManySpaces}, []any{&test.Device1InSpace1}, []any{&test.DeviceSecurityToken})
}

func (suite *ServiceIntegrationSuite) TearDownTest() {
	test.Clean(suite.T(), suite.db)
}

func (suite *ServiceIntegrationSuite) TestCreate() {
	validCreateRequest := schema.CreateRequest{
		OwnerUUID: test.Device1InSpace1.UUID,
		OwnerType: entity.DeviceSecurityType,
	}

	now := time.Now().UTC()

	token, err := suite.Service.Create(suite.ctx, validCreateRequest)
	suite.NoError(err)
	testutils.AssertTokenCreated(suite.T(), entity.Token{
		Base: entity.Base{
			UUID: token.UUID,
		},
		OwnerUUID: validCreateRequest.OwnerUUID,
		Token:     token.Token,
	}, token, now)
}

func (suite *ServiceIntegrationSuite) TestConfirmToken() {
	validConfirmRequest := schema.ConfirmRequest{
		OwnerUUID:  test.Device1InSpace1.UUID,
		OwnerType:  entity.DeviceSecurityType,
		Expiration: 100,
		Token:      test.DeviceSecurityToken.Token, // token hashed sha256
		TokenUUID:  test.DeviceSecurityToken.UUID,
	}

	err := suite.Service.ConfirmToken(suite.ctx, validConfirmRequest)
	suite.NoError(err)
	err = suite.repository.Get(suite.ctx, &test.DeviceSecurityToken)
	suite.ErrorIs(err, sql.ErrNoRows)
}

func (suite *ServiceIntegrationSuite) TestConfirmTokenFail() {
	validConfirmRequest := schema.ConfirmRequest{
		OwnerUUID:  test.Device1InSpace1.UUID,
		OwnerType:  entity.DeviceSecurityType,
		Expiration: 0,
		Token:      test.DeviceSecurityToken.Token, // token hashed sha256
		TokenUUID:  test.DeviceSecurityToken.UUID,
	}

	err := suite.Service.ConfirmToken(suite.ctx, validConfirmRequest)
	suite.ErrorIs(err, errors.ErrTokenExpired)
	err = suite.repository.Get(suite.ctx, &test.DeviceSecurityToken)
	suite.ErrorIs(err, sql.ErrNoRows)
}

func TestServiceIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServiceIntegrationSuite))
}
