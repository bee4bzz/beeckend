package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/authenticator/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/router"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/gaetanDubuc/beeckend/pkg/crypto"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/postgres"
	l "gorm.io/gorm/logger"
)

type APITestSuite struct {
	suite.Suite
	router *gin.Engine
	mock   *sqlmock.Sqlmock
	db     *db.DB

	logger   *log.Logger
	observer *observer.ObservedLogs
	buffer   *bytes.Buffer

	server *httptest.Server

	service        *testutils.Service
	authMiddleware *testutils.Middleware

	loginRootTest          test.APITestCase[schema.Session]
	refreshSessionRootTest test.APITestCase[schema.Session]
	logoutRootTest         test.APITestCase[any]
	publicKeyRootTest      test.APITestCase[schema.PublicKeyResponse]
}

func (suite *APITestSuite) SetupTest() {
	logger, obs, b := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.buffer = b

	mockDb, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	suite.db = db.NewGorm(postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres"},
	), l.Default)
	suite.router, _ = router.New(
		suite.buffer,
		suite.db,
	)

	suite.service = &testutils.Service{}
	suite.authMiddleware = &testutils.Middleware{}

	ValidRSASigningKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	pb, err := crypto.PublicKeyToPem(ValidRSASigningKey.Public())
	if err != nil {
		panic(err)
	}

	RegisterHandlers(
		suite.router.Group(""),
		suite.service,
		suite.authMiddleware,
		string(pb),
		suite.logger,
	)

	suite.server = httptest.NewUnstartedServer(suite.router)
	suite.server.Config = utils.NewServer("", suite.router)
	suite.server.Start()

	host := strings.Split(suite.server.URL, "://")[1]
	suite.loginRootTest = testutils.LoginRootTest.
		WithHost(host).
		WithRouter(suite.router).
		WithLogger(suite.logger)

	suite.refreshSessionRootTest = testutils.RefreshSessionRootTest.
		WithHost(host).
		WithRouter(suite.router).
		WithLogger(suite.logger)

	suite.logoutRootTest = testutils.LogoutRootTest.
		WithHost(host).
		WithRouter(suite.router).
		WithLogger(suite.logger)

	suite.publicKeyRootTest = testutils.PublicKeyRootTest.
		WithHost(host).
		WithRouter(suite.router).
		WithLogger(suite.logger)
}

func (suite *APITestSuite) TearDownTest() {
	suite.authMiddleware.AssertExpectations(suite.T())
	suite.service.AssertExpectations(suite.T())
	suite.T().Log(suite.buffer)
	suite.buffer.Reset()
	suite.observer.TakeAll()
	suite.server.Close()
}

func (suite *APITestSuite) Test_A_User_Can_Login() {
	req := schema.LoginRequest{
		Username: test.ValidUser.Email,
		Password: test.Password,
	}
	suite.authMiddleware.On("OnlyUnauthenticated").
		Return(nil).Once()
	suite.service.On("Login", req).
		Return(schema.Session{
			JWT:        test.ValidJWT,
			RefreshJWT: test.ValidJWT,
		}, nil).Once()

	suite.loginRootTest.WithRequest(
		req,
	).CheckEndpoint(suite.T())
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
