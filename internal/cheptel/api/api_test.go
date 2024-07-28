package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	authtestutils "github.com/gaetanDubuc/beeckend/internal/authenticator/testutils"
	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/router"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	authMiddleware *authtestutils.Middleware
	upgrader       *testutils.Upgrader
	conn           *testutils.WSConn
	resource       *Resource[*testutils.WSConn]

	QueryRootTest test.APITestCase[[]entity.Cheptel]
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
	suite.authMiddleware = &authtestutils.Middleware{}

	upgrader := &websocket.Upgrader{}
	RegisterHandlers(
		suite.router.Group(""),
		suite.service,
		suite.authMiddleware,
		upgrader,
		suite.logger,
	)

	suite.server = httptest.NewUnstartedServer(suite.router)
	suite.server.Config = utils.NewServer("", suite.router)
	suite.server.Start()

	suite.QueryRootTest = testutils.QueryRootTest.
		WithHost(strings.Split(suite.server.URL, "://")[1]).
		WithLogger(suite.logger).
		WithRouter(suite.router)

	suite.upgrader = &testutils.Upgrader{}
	suite.conn = &testutils.WSConn{}

	suite.resource = &Resource[*testutils.WSConn]{
		upgrader:       suite.upgrader,
		service:        suite.service,
		authMiddleware: suite.authMiddleware,
		logger:         suite.logger,
	}

}

func (suite *APITestSuite) TearDownTest() {
	suite.service.AssertExpectations(suite.T())
	suite.T().Log(suite.buffer)
	suite.buffer.Reset()
	suite.observer.TakeAll()
	suite.server.Close()
}

func (suite *APITestSuite) Test_User_Can_Subscribe_To_Cheptels_Modifications() {
	expectedStates := []*[]entity.Cheptel{
		&test.ValidUser.Cheptels,
	}
	suite.service.On("Subscribe", &entity.User{}).
		Return(
			expectedStates,
			nil).Once()

	c := suite.QueryRootTest.
		Dial(suite.T())

	defer c.Close()

	cheptels := &[]entity.Cheptel{}
	actualStates := []*[]entity.Cheptel{}
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = json.Unmarshal(message, &cheptels)
		if err != nil {
			assert.NoError(suite.T(), err)
			suite.T().FailNow()
		}
		actualStates = append(actualStates, cheptels)
	}
	assert.Equal(suite.T(), expectedStates, actualStates)
}

func (suite *APITestSuite) Test_Return_An_Error_When_The_Request_Can_Not_Be_Upgraded() {
	suite.upgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).
		Return(&testutils.WSConn{}, test.AnError).Once()

	assert.PanicsWithError(suite.T(), test.AnError.Error(), func() {
		suite.resource.query(&gin.Context{})
	})
}

func (suite *APITestSuite) Test_Panic_When_The_Service_Can_Not_Subscribe_To_Modifications() {
	suite.upgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).
		Return(suite.conn, nil).Once()

	suite.conn.On("Close").Return(nil).Once()

	suite.service.On("Subscribe", &entity.User{}).
		Return(
			nil,
			test.AnError).Once()

	assert.PanicsWithError(suite.T(), test.AnError.Error(), func() {
		suite.resource.query(&gin.Context{
			Request: &http.Request{},
		})
	})
}

func (suite *APITestSuite) Test_Panic_When_The_Resource_Can_Not_Write_Message_To_The_Websocket() {
	suite.upgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).
		Return(suite.conn, nil).Once()

	expectedStream := []*[]entity.Cheptel{
		&test.ValidUser.Cheptels,
	}
	suite.service.On("Subscribe", &entity.User{}).
		Return(
			expectedStream,
			nil).Once()

	suite.conn.On("Close").Return(test.AnError).Once()
	suite.conn.On("WriteMessage", mock.Anything, mock.Anything).
		Return(test.AnError).Once()

	assert.PanicsWithError(suite.T(), test.AnError.Error(), func() {
		suite.resource.query(&gin.Context{
			Request: &http.Request{},
		})
	})
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
