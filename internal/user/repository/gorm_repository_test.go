package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	mock       *sqlmock.Sqlmock
	db         *gorm.DB
	Repository *GormRepository
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	mockDb, mock, _ := sqlmock.New()
	suite.mock = &mock
	suite.db = db.NewGorm(postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres"},
	), logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{},
	))
	suite.Repository = NewGormRepository(suite.db)
}

// this function executes after all tests executed
func (suite *RepositoryTestSuite) TearDownSuite() {
	// we make sure that all expectations were met
	if err := (*suite.mock).ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
