package api

import (
	"bytes"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	testutils "github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/router"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/gaetanDubuc/beeckend/pkg/pagination"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	QueryRootTest test.APITestCase[pagination.Pages[[]entity.Cheptel]]
}

func (suite *APITestSuite) SetupSuite() {
	logger, obs, b := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.buffer = b

	mockDb, mock, _ := sqlmock.New()
	suite.mock = &mock
	suite.db = db.NewGorm(postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres"},
	), l.Default)
	suite.router, _ = router.New(
		suite.db,
	)

	upgrader := &websocket.Upgrader{}
	RegisterHandlers(suite.router.Group(""), upgrader, suite.logger)

	suite.QueryRootTest = testutils.QueryRootTest.
		WithLogger(suite.logger).
		WithRouter(suite.router)

}

func (suite *APITestSuite) SetupTest() {
	suite.buffer.Reset()
	suite.observer.TakeAll()
}

func (suite *APITestSuite) TearDownSuite() {
	if err := (*suite.mock).ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
	suite.T().Log(suite.buffer)
}

func (suite *APITestSuite) Test_User_Can_Query_Cheptels() {
	suite.QueryRootTest.
		CheckEndpoint(suite.T())
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
