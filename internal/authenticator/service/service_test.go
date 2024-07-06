package auth

import (
	"bytes"
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	auth "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth"
	refreshtokenschema "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/refresh-token/schema"
	rfshtestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/refresh-token/testutils"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/schema"
	authtestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/testutils"
	cryptotestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/crypto/testutils"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/test"
	usertestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/user/testutils/v2"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
	"go.uber.org/zap/zaptest/observer"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx      context.Context
	logger   log.Logger
	observer *observer.ObservedLogs
	buffer   *bytes.Buffer

	service             *Service
	userRepository      *usertestutils.Repository
	refreshTokenService *rfshtestutils.Repository
	jwtGenerator        *authtestutils.JWTGenerator
	hasher              *cryptotestutils.Hasher
}

// this function executes before the test suite begins execution.
func (suite *ServiceTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, obs, b := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.buffer = b

	suite.userRepository = &usertestutils.Repository{}
	suite.refreshTokenService = &rfshtestutils.Repository{}
	suite.jwtGenerator = &authtestutils.JWTGenerator{}
	suite.hasher = &cryptotestutils.Hasher{}

	suite.service = NewService(
		suite.userRepository,
		suite.refreshTokenService,
		suite.jwtGenerator,
		suite.hasher,
		10,
		logger,
	)
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.buffer.Reset()
	suite.observer.TakeAll()
}

func (suite *ServiceTestSuite) TearDownTest() {
	suite.userRepository.AssertExpectations(suite.T())
	suite.refreshTokenService.AssertExpectations(suite.T())
	suite.jwtGenerator.AssertExpectations(suite.T())
	suite.hasher.AssertExpectations(suite.T())
	suite.T().Log(suite.buffer)
}

func (suite *ServiceTestSuite) Test_A_User_Can_Login() {
	// Given
	user := test.UserNoAdminConfirmed
	suite.userRepository.On(
		"GetFromUsername",
		user.Email).Return(user, nil).Once()
	suite.hasher.On(
		"AreSameHash",
		test.ValidPassword,
		user.Password).Return(true).Once()
	suite.jwtGenerator.On(
		"GenerateJWT",
		user.UUID,
	).Return("jwt", nil).Once()
	suite.refreshTokenService.On(
		"Create",
		refreshtokenschema.CreateRequest{
			UserUUID: user.UUID,
		}).Return("refresh_token", nil).Once()

	// When
	session, err := suite.service.Login(suite.ctx, schema.LoginRequest{
		Username: user.Email,
		Password: test.ValidPassword,
	})

	// Then
	suite.NoError(err)
	suite.Equal("jwt", session.Token)
	suite.Equal("refresh_token", session.RefreshToken)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_If_A_User_Login_With_Wrong_Password() {
	// Given
	user := test.UserNoAdminConfirmed
	suite.userRepository.On(
		"GetFromUsername",
		user.Email).Return(user, nil).Once()
	suite.hasher.On(
		"AreSameHash",
		test.InvalidPassword,
		user.Password).Return(false).Once()

	// When
	session, err := suite.service.Login(suite.ctx, schema.LoginRequest{
		Username: user.Email,
		Password: test.InvalidPassword,
	})

	// Then
	suite.Empty(session)
	suite.ErrorIs(err, auth.ErrWrongPassword)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_If_A_User_Login_With_Invalid_Username() {
	// Given
	suite.userRepository.On(
		"GetFromUsername",
		test.UserNoAdminConfirmed.Email).Return(entity.User{}, sql.ErrNoRows).Once()

	// When
	session, err := suite.service.Login(suite.ctx, schema.LoginRequest{
		Username: test.UserNoAdminConfirmed.Email,
		Password: test.ValidPassword,
	})

	// Then
	suite.Empty(session)
	suite.ErrorIs(err, sql.ErrNoRows)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_When_The_JWT_Generation_Fails() {
	// Given
	user := test.UserNoAdminConfirmed2
	suite.userRepository.On(
		"GetFromUsername",
		user.Email).Return(user, nil).Once()
	suite.hasher.On(
		"AreSameHash",
		test.ValidPassword,
		user.Password).Return(true).Once()
	suite.jwtGenerator.On(
		"GenerateJWT",
		user.UUID,
	).Return("", test.ErrMocked).Once()

	// When
	session, err := suite.service.Login(suite.ctx, schema.LoginRequest{
		Username: user.Email,
		Password: test.ValidPassword,
	})

	// Then
	suite.Empty(session)
	suite.ErrorIs(err, test.ErrMocked)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_When_The_Refresh_Token_Creation_Fails() {
	// Given
	user := test.UserNoAdminConfirmed
	suite.userRepository.On(
		"GetFromUsername",
		user.Email).Return(user, nil).Once()
	suite.hasher.On(
		"AreSameHash",
		test.ValidPassword,
		user.Password).Return(true).Once()
	suite.jwtGenerator.On(
		"GenerateJWT",
		user.UUID,
	).Return("jwt", nil).Once()
	suite.refreshTokenService.On(
		"Create",
		refreshtokenschema.CreateRequest{
			UserUUID: user.UUID,
		}).Return("", test.ErrMocked).Once()

	// When
	session, err := suite.service.Login(suite.ctx, schema.LoginRequest{
		Username: user.Email,
		Password: test.ValidPassword,
	})

	// Then
	suite.Empty(session)
	suite.ErrorIs(err, test.ErrMocked)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_When_The_Login_Request_Is_Invalid() {
	// When
	session, err := suite.service.Login(suite.ctx, schema.LoginRequest{
		Username: "",
		Password: test.ValidPassword,
	})

	// Then
	suite.Empty(session)
	suite.Error(err)
}

func (suite *ServiceTestSuite) Test_A_User_Can_Refresh_A_Token() {
	// Given
	UUID := uuid.New()
	suite.refreshTokenService.On(
		"Refresh",
		refreshtokenschema.RefreshRequest{
			UserUUID:  test.UserNoAdminConfirmed.UUID,
			TokenUUID: UUID,
			Token:     "refresh_token",
		}).Return("new_refresh_token", nil).Once()

	suite.userRepository.On(
		"Get",
		test.UserNoAdminConfirmed.UUID).
		Return(test.UserNoAdminConfirmed, nil).Once()

	suite.jwtGenerator.On(
		"GenerateJWT",
		test.UserNoAdminConfirmed.UUID,
	).Return("jwt", nil).Once()

	// When
	session, err := suite.service.RefreshSession(suite.ctx,
		refreshtokenschema.RefreshRequest{
			UserUUID:  test.UserNoAdminConfirmed.UUID,
			TokenUUID: UUID,
			Token:     "refresh_token",
		},
	)

	// Then
	suite.NoError(err)
	suite.Equal("jwt", session.Token)
	suite.Equal("new_refresh_token", session.RefreshToken)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_If_A_User_Tries_To_Refresh_A_Token_With_Invalid_User() {
	// Given
	req := refreshtokenschema.RefreshRequest{
		UserUUID:  uuid.New(),
		TokenUUID: uuid.New(),
		Token:     "refresh_token",
	}
	suite.refreshTokenService.On(
		"Refresh",
		req,
	).Return("new_refresh_token", nil).Once()

	suite.userRepository.On(
		"Get",
		req.UserUUID).
		Return(entity.User{}, sql.ErrNoRows).Once()

	// When
	session, err := suite.service.RefreshSession(
		suite.ctx,
		req,
	)

	// Then
	suite.Empty(session)
	suite.ErrorIs(err, sql.ErrNoRows)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_If_A_User_Tries_To_Refresh_A_Token_With_Invalid_Refresh_Token() {
	// Given
	req := refreshtokenschema.RefreshRequest{
		UserUUID:  test.UserNoAdminConfirmed.UUID,
		TokenUUID: uuid.New(),
		Token:     "refresh_token",
	}

	suite.refreshTokenService.On(
		"Refresh",
		req,
	).Return("", sql.ErrNoRows).Once()

	// When
	session, err := suite.service.RefreshSession(
		suite.ctx,
		req,
	)

	// Then
	suite.Empty(session)
	suite.ErrorIs(err, sql.ErrNoRows)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_When_A_User_Tries_To_Refresh_A_Token_And_The_Generator_Fails() {
	// Given
	UUID := uuid.New()
	suite.refreshTokenService.On(
		"Refresh",
		refreshtokenschema.RefreshRequest{
			UserUUID:  test.UserNoAdminConfirmed.UUID,
			TokenUUID: UUID,
			Token:     "refresh_token",
		}).Return("new_refresh_token", nil).Once()

	suite.userRepository.On(
		"Get",
		test.UserNoAdminConfirmed.UUID).
		Return(test.UserNoAdminConfirmed, nil).Once()

	suite.jwtGenerator.On(
		"GenerateJWT",
		test.UserNoAdminConfirmed.UUID,
	).Return("", test.ErrMocked).Once()

	// When
	session, err := suite.service.RefreshSession(suite.ctx,
		refreshtokenschema.RefreshRequest{
			UserUUID:  test.UserNoAdminConfirmed.UUID,
			TokenUUID: UUID,
			Token:     "refresh_token",
		},
	)

	// Then
	suite.Empty(session)
	suite.ErrorIs(err, test.ErrMocked)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_When_A_User_Tries_To_Refresh_A_Token_With_Invalid_Refresh_Token() {
	// When
	session, err := suite.service.RefreshSession(
		suite.ctx,
		refreshtokenschema.RefreshRequest{},
	)

	// Then
	suite.Empty(session)
	suite.Error(err)
}

func (suite *ServiceTestSuite) Test_A_User_Can_Logout() {
	// Given
	suite.refreshTokenService.On(
		"DeleteFromUser",
		refreshtokenschema.DeleteFromUserRequest{
			UserUUID: test.UserNoAdminConfirmed.UUID,
		}).Return(nil).Once()

	// When
	err := suite.service.Logout(suite.ctx, schema.LogoutRequest{
		UserUUID: test.UserNoAdminConfirmed.UUID,
	})

	// Then
	suite.NoError(err)
}

func (suite *ServiceTestSuite) Test_Return_An_Error_When_A_User_Tries_To_Logout_With_Invalid_Request() {
	// When
	err := suite.service.Logout(suite.ctx, schema.LogoutRequest{})

	// Then
	suite.Error(err)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

/*
var (
	MockLogger, _, _ = log.NewForTest()
)

func InitService() (s Service, JWTManager *testutils.JWTManager[any], userRepo *usertestutils.Repository, hasher *cryptotestutils.Hasher) {
	JWTManager = &testutils.JWTManager[any]{
		GenerateJWTReturn: test.ValidJWT,
	}
	userRepo = &usertestutils.Repository{
		GetFromUsernameReturn: test.UserAdminConfirmed,
	}
	hasher = &cryptotestutils.Hasher{
		AreSameHashReturn: true,
	}
	s = NewService[any](userRepo, JWTManager, hasher, test.ValidTokenExpiration, MockLogger)
	return
}

func TestService_Login(t *testing.T) {
	t.Run("Login success", func(t *testing.T) {
		// prepare
		s, _, _, _ := InitService()

		// test
		jwt, token, err := s.Login(context.Background(), test.UserAdminConfirmed.Email, test.ValidPassword)

		// assert
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		AssertJWT(t, jwt, test.ValidJWT)
	})

	t.Run("fail to login due to wrong password", func(t *testing.T) {
		s, _, _, hasher := InitService()
		hasher.AreSameHashReturn = false

		jwt, token, err := s.Login(context.Background(), test.UserAdminConfirmed.Email, test.InvalidPassword)

		assert.Empty(t, token)
		assert.Empty(t, jwt)
		assert.ErrorIs(t, err, auth.ErrWrongPassword)
	})

	t.Run("fail due to internal error", func(t *testing.T) {
		s, _, userRepo, _ := InitService()
		userRepo.GetFromUsernameReturn = entity.User{}
		userRepo.GetFromUsernameError = testutils.ErrMocked

		jwt, token, err := s.Login(context.Background(), test.UserAdminConfirmed.Email, test.ValidPassword)

		assert.Empty(t, token)
		assert.ErrorIs(t, err, ErrGetUser)
		assert.ErrorIs(t, err, testutils.ErrMocked)
		assert.Empty(t, jwt)
	})

	t.Run("fail due to generate token", func(t *testing.T) {
		s, authRepo, _, _ := InitService()
		authRepo.GenerateJWTError = testutils.ErrMocked

		jwt, token, err := s.Login(context.Background(), test.UserAdminConfirmed.Email, test.ValidPassword)

		assert.Empty(t, token)
		assert.ErrorIs(t, err, errors.ErrInternalServer)
		assert.ErrorIs(t, err, testutils.ErrMocked)
		assert.Empty(t, jwt)
	})
}

func TestService_GenerateJWT(t *testing.T) {
	t.Run("generate jwt succeed", func(t *testing.T) {
		// prepare
		s, _, _, _ := InitService()

		// test
		token, err := s.GenerateJWT(context.Background(), test.UserAdminConfirmed)

		// assert
		assert.Nil(t, err)
		AssertJWT(t, token, test.ValidJWT)
	})

	t.Run("make user claim succeed", func(t *testing.T) {
		// test
		claim := MakeUserClaim(test.UserAdminConfirmed, test.ValidTokenExpiration)

		// assert
		AssertUserClaim(t, claim, test.UserAdminConfirmed)
	})

	t.Run("fail to generate a jwt", func(t *testing.T) {
		s, authRepo, _, _ := InitService()
		authRepo.GenerateJWTError = testutils.ErrMocked

		token, err := s.GenerateJWT(context.Background(), test.UserAdminConfirmed)

		assert.Empty(t, token)
		assert.ErrorIs(t, err, testutils.ErrMocked)
	})
}

func AssertJWT(t *testing.T, jwt string, expected string) {
	assert.NotEmpty(t, jwt)
	assert.Equal(t, expected, jwt)
}

func AssertUserClaim(t *testing.T, claim *JWTClaims, expected entity.User) {
	assert.NotEmpty(t, claim)
	assert.Equal(t, expected.UUID.String(), claim.UUID.String())
	assert.Equal(t, expected.Administrator, claim.Administrator.ElseZero())
	assert.Equal(t, expected.Confirmed, claim.Confirmed.ElseZero())
	assert.Less(t, expected.CreatedAt.UnixMilli()/1000, claim.Exp)
}
*/
