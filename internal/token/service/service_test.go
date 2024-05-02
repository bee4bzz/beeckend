package service

import (
	"context"
	"testing"
	"time"

	"4d63.com/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	cryptotestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/crypto/testutils"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/test"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/testutils"
	errors "gitlab.com/fogo-dev/infrastructure/web-api/internal/token"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/token/schema"
	tokentestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/token/testutils"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
	"go.uber.org/zap/zaptest/observer"
)

var (
	token = test.DeviceSecurityToken
)

type ServiceTestSuite struct {
	suite.Suite
	ctx        context.Context
	service    *Service
	hasher     *cryptotestutils.Hasher
	repository *tokentestutils.Repository
	logger     *log.Logger
	observer   *observer.ObservedLogs
}

// this function executes before the test suite begins execution.
func (suite *ServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, obs, _ := log.NewForTest()
	suite.logger = &logger
	suite.observer = obs

	suite.hasher = &cryptotestutils.Hasher{}
	suite.repository = &tokentestutils.Repository{}
	suite.service = NewService(suite.repository, suite.hasher, logger)
}

func (suite *ServiceTestSuite) TearDownSuite() {
	suite.hasher.AssertExpectations(suite.T())
	suite.repository.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestCreate() {
	validCreateReq := schema.CreateRequest{
		UUID:      optional.Of(token.UUID),
		OwnerUUID: token.OwnerUUID,
		OwnerType: entity.DeviceSecurityType,
	}

	suite.hasher.On("GenerateToken").Return(token.Token).Once()
	suite.hasher.On("Hash", token.Token).Return(token.HashedToken, nil).Once()
	suite.repository.On(
		"Create",
		token.UUID,
		token.OwnerUUID,
		entity.DeviceSecurityType,
		token.HashedToken,
		token.Token,
	).Return(token, nil).Once()

	t, err := suite.service.Create(suite.ctx, validCreateReq)
	assert.NoError(suite.T(), err)
	tokentestutils.AssertToken(suite.T(), token, t)
}

func (suite *ServiceTestSuite) TestCreateFail() {
	validCreateReq := schema.CreateRequest{
		UUID:      optional.Of(token.UUID),
		OwnerUUID: token.OwnerUUID,
		OwnerType: entity.DeviceSecurityType,
	}

	testcases := []struct {
		name string
		req  schema.CreateRequest
		err  error
		fn   func()
	}{
		{
			name: "fail to create a token when hashing the token",
			req:  validCreateReq,
			err:  testutils.ErrMocked,
			fn: func() {
				suite.hasher.On("GenerateToken").Return(token.Token).Once()
				suite.hasher.On("Hash", token.Token).Return("", testutils.ErrMocked).Once()
			},
		},
		{
			name: "fail to create an existing token",
			req:  validCreateReq,
			err:  testutils.ErrMocked,
			fn: func() {
				suite.hasher.On("GenerateToken").Return(token.Token).Once()
				suite.hasher.On("Hash", token.Token).Return(token.HashedToken, nil).Once()
				suite.repository.On(
					"Create",
					token.UUID,
					token.OwnerUUID,
					entity.DeviceSecurityType,
					token.HashedToken,
					token.Token).Return(token, testutils.ErrMocked).Once()
			},
		},
		{
			name: "fail to create a token when the request is invalid",
			req:  schema.CreateRequest{},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn()
			}
			token, err := suite.service.Create(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, token)
		})
	}
}

func (suite *ServiceTestSuite) TestConfirmToken() {
	validConfirmReq := schema.ConfirmRequest{
		OwnerUUID:  token.OwnerUUID,
		OwnerType:  entity.DeviceSecurityType,
		Expiration: test.ValidTokenExpiration,
		TokenUUID:  token.UUID,
		Token:      token.Token,
	}

	suite.T().Run("succeed to confirm an existing token", func(t *testing.T) {
		suite.repository.On(
			"GetBy",
			map[string]any{"UUID": token.UUID, "owner_UUID": token.OwnerUUID},
			&entity.Token{OwnerType: entity.DeviceSecurityType}).Return(token, nil).Once()
		suite.hasher.On("AreSameHash", token.Token, token.HashedToken).Return(true).Once()
		suite.repository.On("Delete", &token).Return(nil).Once()

		err := suite.service.ConfirmToken(suite.ctx, validConfirmReq)

		assert.NoError(t, err)
	})
	suite.T().Run("succeed to confirm an existing token if the updated time is not expired", func(t *testing.T) {
		tokenCopy := token
		tokenCopy.CreatedAt = time.Now().UTC().Add(-(test.ValidTokenExpiration + 1) * time.Second)

		suite.repository.On(
			"GetBy",
			map[string]any{"UUID": tokenCopy.UUID, "owner_UUID": tokenCopy.OwnerUUID},
			&entity.Token{OwnerType: entity.DeviceSecurityType}).Return(tokenCopy, nil).Once()
		suite.hasher.On("AreSameHash", tokenCopy.Token, tokenCopy.HashedToken).Return(true).Once()
		suite.repository.On("Delete", &tokenCopy).Return(nil).Once()

		err := suite.service.ConfirmToken(suite.ctx, validConfirmReq)

		assert.NoError(t, err, errors.ErrTokenExpired)
	})
}

func (suite *ServiceTestSuite) TestConfirmTokenFail() {
	validConfirmReq := schema.ConfirmRequest{
		OwnerUUID:  token.OwnerUUID,
		OwnerType:  entity.DeviceSecurityType,
		Expiration: test.ValidTokenExpiration,
		TokenUUID:  token.UUID,
		Token:      token.Token,
	}

	testcases := []struct {
		name string
		req  schema.ConfirmRequest
		err  error
		fn   func()
	}{
		{
			name: "fail to confirm an existing token when the request is invalid",
			req:  schema.ConfirmRequest{},
		},
		{
			name: "fail to confirm an existing token because the token is not the same",
			req:  validConfirmReq,
			err:  errors.ErrTokenInvalid,
			fn: func() {
				suite.repository.On(
					"GetBy",
					map[string]any{"UUID": token.UUID, "owner_UUID": token.OwnerUUID},
					&entity.Token{OwnerType: entity.DeviceSecurityType}).Return(token, nil).Once()
				suite.hasher.On("AreSameHash", token.Token, token.HashedToken).Return(false).Once()
			},
		},
		{
			name: "fail to confirm an existing token because the token is expired",
			req:  validConfirmReq,
			err:  errors.ErrTokenExpired,
			fn: func() {
				tokenCopy := token
				tokenCopy.CreatedAt = time.Now().UTC().Add(-(test.ValidTokenExpiration + 1) * time.Second)
				tokenCopy.UpdatedAt = time.Now().UTC().Add(-(test.ValidTokenExpiration + 1) * time.Second)

				suite.repository.On(
					"GetBy",
					map[string]any{"UUID": tokenCopy.UUID, "owner_UUID": tokenCopy.OwnerUUID},
					&entity.Token{OwnerType: entity.DeviceSecurityType},
				).Return(tokenCopy, nil).Once()
				suite.hasher.On("AreSameHash", tokenCopy.Token, tokenCopy.HashedToken).Return(true).Once()
				suite.repository.On("Delete", &tokenCopy).Return(nil).Once()
			},
		},
		{
			name: "fail to confirm an existing token can't be deleted",
			req:  validConfirmReq,
			err:  testutils.ErrMocked,
			fn: func() {
				suite.repository.On(
					"GetBy",
					map[string]any{"UUID": token.UUID, "owner_UUID": token.OwnerUUID},
					&entity.Token{OwnerType: entity.DeviceSecurityType},
				).Return(token, nil).Once()
				suite.hasher.On("AreSameHash", token.Token, token.HashedToken).Return(true).Once()
				suite.repository.On("Delete", &token).Return(testutils.ErrMocked).Once()
			},
		},
		{
			name: "fail to confirm an existing token because the token is not found",
			req:  validConfirmReq,
			err:  testutils.ErrMocked,
			fn: func() {
				suite.repository.On(
					"GetBy",
					map[string]any{"UUID": token.UUID, "owner_UUID": token.OwnerUUID},
					&entity.Token{OwnerType: entity.DeviceSecurityType},
				).Return(token, testutils.ErrMocked).Once()
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn()
			}
			err := suite.service.ConfirmToken(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			suite.hasher.AssertExpectations(t)
			suite.repository.AssertExpectations(t)
		})
	}
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
