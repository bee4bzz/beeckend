package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
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

func (suite *RepositoryTestSuite) TestQueryByUser() {
	(*suite.mock).ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."deleted_at" IS NULL AND "users"\."id" = \$1`).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name"}).AddRow(test.ValidUser.ID, test.ValidUser.Email, test.ValidUser.Name))
	(*suite.mock).ExpectQuery(`SELECT \* FROM "user_cheptels" WHERE "user_cheptels"\."user_id" = \$1`).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"user_id", "cheptel_id"}).AddRow(test.ValidUser.ID, test.ValidCheptel.ID))
	(*suite.mock).ExpectQuery(`SELECT \* FROM "cheptels" WHERE "cheptels"\."id" = \$1`).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(test.ValidCheptel.ID, test.ValidCheptel.Name))
	(*suite.mock).ExpectQuery(`SELECT \* FROM "hives" WHERE "hives"\."cheptel_id" = \$1`).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "cheptel_id"}).AddRow(test.ValidHive.ID, test.ValidHive.Name, test.ValidCheptel.ID))
	(*suite.mock).ExpectQuery(`SELECT \* FROM "hive_notes" WHERE "hive_notes"\."hive_id" = \$1 AND "hive_notes"\."deleted_at" IS NULL`).WithArgs(test.ValidHive.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "hive_id", "operation"}).AddRow(test.ValidHiveNote.ID, test.ValidHiveNote.Name, test.ValidHiveNote.HiveID, test.ValidHiveNote.Operation))

	(*suite.mock).ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."deleted_at" IS NULL AND "users"\."id" = \$1`).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name"}))
	(*suite.mock).ExpectQuery(`SELECT \* FROM "user_cheptels" WHERE "user_cheptels"\."user_id" = \$1`).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"user_id", "cheptel_id"}))

	testcases := []struct {
		entity.User
		len int
	}{
		{test.ValidUser, 1},
		{entity.User{Model: gorm.Model{ID: 100}}, 0},
	}

	for _, tc := range testcases {
		hiveNotes := []entity.HiveNote{}
		err := suite.Repository.QueryByUser(suite.ctx, &tc.User, &hiveNotes)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), hiveNotes, tc.len)
	}
}

func (suite *RepositoryTestSuite) TestQueryByUserFail() {
	(*suite.mock).ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."deleted_at" IS NULL AND "users"\."id" = \$1`).WithArgs(test.ValidUser.ID).WillReturnError(sql.ErrTxDone)

	hiveNotes := []entity.HiveNote{}
	err := suite.Repository.QueryByUser(suite.ctx, &test.ValidUser, &hiveNotes)
	assert.ErrorIs(suite.T(), err, sql.ErrTxDone)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
