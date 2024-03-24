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
	(*suite.mock).ExpectQuery(
		`SELECT \* FROM "hive_notes" WHERE "hive_notes"\."hive_id" = \$1 AND "hive_notes"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidHive.ID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "hive_id"}).AddRow(test.ValidHiveNote.ID, test.ValidHive.ID),
	)

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
		hives := []entity.Hive{}
		err := suite.Repository.QueryByUser(suite.ctx, &tc.User, &hives)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), hives, tc.len)
	}
}

func (suite *RepositoryTestSuite) TestQueryByUserFail() {
	(*suite.mock).ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."deleted_at" IS NULL AND "users"\."id" = \$1`).WithArgs(test.ValidUser.ID).WillReturnError(sql.ErrTxDone)

	hives := []entity.Hive{}
	err := suite.Repository.QueryByUser(suite.ctx, &test.ValidUser, &hives)
	assert.ErrorIs(suite.T(), err, sql.ErrTxDone)
}

func (suite *RepositoryTestSuite) TestSoftDelete() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(
		`UPDATE "hives" SET "deleted_at"=\$1 WHERE \("hives"\."id" = \$2 AND "hives"\."name" = \$3 AND "hives"\."cheptel_id" = \$4\) AND "hives"\."deleted_at" IS NULL`,
	).WithArgs(sqlmock.AnyArg(), test.ValidHive.ID, test.ValidHive.Name, test.ValidHive.CheptelID).WillReturnResult(sqlmock.NewResult(0, 1))
	(*suite.mock).ExpectCommit()

	err := suite.Repository.SoftDelete(suite.ctx, &test.ValidHive)
	assert.NoError(suite.T(), err)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
