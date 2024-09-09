package api

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/router"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
)

type APITestSuite struct {
	suite.Suite

	logger   *log.Logger
	observer *observer.ObservedLogs
	buffer   *bytes.Buffer

	service        *testutils.Service
	authMiddleware *testutils.Middleware
	router         *gin.Engine
}

func (suite *APITestSuite) SetupTest() {
	logger, obs, b := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.buffer = b

	suite.router, _ = router.New(
		suite.buffer,
		&db.DB{},
	)

	suite.service = &testutils.Service{}
	suite.authMiddleware = &testutils.Middleware{}

	RegisterHandlers(
		suite.router.Group(""),
		suite.service,
		suite.authMiddleware,
		"pbk",
		suite.logger,
	)
}

func (suite *APITestSuite) TearDownTest() {
	suite.authMiddleware.AssertExpectations(suite.T())
	suite.service.AssertExpectations(suite.T())
	suite.T().Log(suite.buffer)
	suite.buffer.Reset()
	suite.observer.TakeAll()
}

func (suite *APITestSuite) Test_A_User_Can_Login() {
	req := schema.LoginRequest{
		Username: test.ValidUser.Email,
		Password: test.Password,
	}
	suite.authMiddleware.On("OnlyUnauthenticated").
		Return().Once()
	suite.service.On("Login", req).
		Return(schema.Session{
			JWT:        "JWT",
			RefreshJWT: "RefreshJWT",
		}, nil).Once()

	testutils.LoginRootTest.
		WithLogger(suite.logger).
		WithRouter(suite.router).
		WithRequest(
			req,
		).CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_Return_An_Error_If_A_User_Login_With_An_Invalid_Username() {
	suite.authMiddleware.On("OnlyUnauthenticated").
		Return().Once()

	testutils.LoginRootTest.
		WithLogger(suite.logger).
		WithRouter(suite.router).
		WithWantStatus(http.StatusBadRequest).
		WithWantResponse("").
		WithBody(
			"",
		).CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_Return_An_Error_If_The_Service_Fails_To_Login() {
	req := schema.LoginRequest{
		Username: test.ValidUser.Email,
		Password: test.Password,
	}
	suite.authMiddleware.On("OnlyUnauthenticated").
		Return().Once()
	suite.service.On("Login", req).
		Return(schema.Session{}, test.AnError).Once()

	testutils.LoginRootTest.
		WithLogger(suite.logger).
		WithRouter(suite.router).
		WithRequest(
			req,
		).WithWantStatus(http.StatusUnauthorized).
		WithWantResponse("").
		CheckEndpoint(suite.T())
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
