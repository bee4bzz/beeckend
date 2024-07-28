package middleware

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/authenticator/testutils"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"
)

type MiddlewareTestSuite struct {
	suite.Suite

	logger   *log.Logger
	observer *observer.ObservedLogs
	buffer   *bytes.Buffer

	c        *gin.Context
	recorder *httptest.ResponseRecorder

	UserRepository *testutils.UserRepository
	middleware     *Middleware
}

func (suite *MiddlewareTestSuite) SetupTest() {
	logger, obs, b := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.buffer = b

	suite.c, suite.recorder = test.NewContext()

	suite.UserRepository = &testutils.UserRepository{}
	suite.middleware = New(
		"HS256",
		func(t *jwt.Token) (interface{}, error) { return []byte("secret"), nil },
		suite.UserRepository,
		suite.logger)
}

func (suite *MiddlewareTestSuite) TearDownTest() {
	suite.UserRepository.AssertExpectations(suite.T())
	suite.T().Log(suite.buffer)
	suite.buffer.Reset()
	suite.observer.TakeAll()
}

func (suite *MiddlewareTestSuite) Test_Authentication_Is_Handled_Without_Error() {
	suite.UserRepository.On("Get", &entity.User{
		Model: gorm.Model{ID: test.ValidUser.ID},
	}).Return(test.ValidUser, nil).Once()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, schema.MakeUserClaim(
		test.ValidUser.ID,
		1,
	))
	key, err := suite.middleware.keyfunc(token)
	assert.NoError(suite.T(), err)
	tokenString, err := token.SignedString(key)

	suite.c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	suite.middleware.AuthHandler(suite.c)

	// Assert
	assert.False(suite.T(), suite.c.IsAborted())
	assert.Equal(suite.T(), suite.recorder.Code, http.StatusOK)

	user := suite.middleware.CurrentAuthenticatedUser(suite.c.Request.Context())
	assert.Equal(suite.T(), test.ValidUser, user)
}

func (suite *MiddlewareTestSuite) Test_Abort_When_The_Token_Is_Absent() {
	suite.middleware.AuthHandler(suite.c)

	// Assert
	assert.True(suite.T(), suite.c.IsAborted())
	assert.Equal(suite.T(), http.StatusUnauthorized, suite.recorder.Code)

	assert.PanicsWithError(suite.T(), ErrUserNotAuthenticated.Error(), func() {
		suite.middleware.CurrentAuthenticatedUser(suite.c.Request.Context())
	})
}

func (suite *MiddlewareTestSuite) Test_Abort_When_The_Repository_Could_Not_Find_The_User() {
	suite.UserRepository.On("Get", &entity.User{
		Model: gorm.Model{ID: test.ValidUser.ID},
	}).Return(entity.User{}, sql.ErrNoRows).Once()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, schema.MakeUserClaim(
		test.ValidUser.ID,
		1,
	))
	key, err := suite.middleware.keyfunc(token)
	assert.NoError(suite.T(), err)
	tokenString, err := token.SignedString(key)

	suite.c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	suite.middleware.AuthHandler(suite.c)

	// Assert
	assert.True(suite.T(), suite.c.IsAborted())
	assert.Equal(suite.T(), http.StatusUnauthorized, suite.recorder.Code)
	assert.PanicsWithError(suite.T(), ErrUserNotAuthenticated.Error(), func() {
		suite.middleware.CurrentAuthenticatedUser(suite.c.Request.Context())
	})
}

func (suite *MiddlewareTestSuite) Test_Abort_When_The_JWT_Is_Incomplete() {
	now := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &schema.Claims{
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "beeckend",
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(1) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	})
	key, err := suite.middleware.keyfunc(token)
	assert.NoError(suite.T(), err)
	tokenString, err := token.SignedString(key)

	suite.c.Request.Header.Set("Authorization", "Bearer "+tokenString)

	suite.middleware.AuthHandler(suite.c)

	// Assert
	assert.True(suite.T(), suite.c.IsAborted())
	assert.Equal(suite.T(), http.StatusUnauthorized, suite.recorder.Code)
	assert.PanicsWithError(suite.T(), ErrUserNotAuthenticated.Error(), func() {
		suite.middleware.CurrentAuthenticatedUser(suite.c.Request.Context())
	})
}

func (suite *MiddlewareTestSuite) Test_Only_Unauthenticated_User_Is_Accepted() {
	suite.middleware.OnlyUnauthenticated(suite.c)

	// Assert
	assert.False(suite.T(), suite.c.IsAborted())
	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
}

func (suite *MiddlewareTestSuite) Test_Abort_When_The_Only_Unauthenticated_User_Is_Accepted() {
	suite.c.Request.Header.Set("Authorization", "Bearer token")
	suite.middleware.OnlyUnauthenticated(suite.c)

	// Assert
	assert.True(suite.T(), suite.c.IsAborted())
	assert.Equal(suite.T(), http.StatusForbidden, suite.recorder.Code)
}

func TestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
