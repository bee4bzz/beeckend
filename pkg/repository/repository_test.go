package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	mock       *sqlmock.Sqlmock
	db         *gorm.DB
	Repository *Repository[mock]
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	mockDb, m, _ := sqlmock.New()
	suite.mock = &m
	db, _ := gorm.Open(postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres"},
	))
	suite.db = db
	suite.Repository = NewRepository[mock](db)
}

// this function executes after all tests executed
func (suite *RepositoryTestSuite) TearDownSuite() {
	// we make sure that all expectations were met
	if err := (*suite.mock).ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RepositoryTestSuite) TestDB() {
	assert.Equal(suite.T(), suite.db, suite.Repository.DB())
}

func (suite *RepositoryTestSuite) TestGet() {
	(*suite.mock).ExpectQuery(`SELECT \* FROM "mocks" WHERE "mocks"\."id" = \$1 AND "mocks"\."deleted_at" IS NULL AND "mocks"\."id" = \$2 ORDER BY "mocks"\."id" LIMIT \$3`).WithArgs(1, 1, 1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	t := suite.T()
	mock := &mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.Get(suite.ctx, mock)
	assert.NoError(t, err)
}

func (suite *RepositoryTestSuite) TestCreate() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectQuery(`INSERT INTO "mocks" \("created_at","updated_at","deleted_at","attr"\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING "id"`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 0).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	(*suite.mock).ExpectCommit()

	t := suite.T()
	now := time.Now()
	mock := &mock{}
	err := suite.Repository.Create(suite.ctx, mock)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, mock.CreatedAt, now, "CreatedAt should be greater than now")
	assert.GreaterOrEqual(t, mock.UpdatedAt, now, "UpdatedAt should be greater than now")
}

func (suite *RepositoryTestSuite) TestUpdate() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(`UPDATE "mocks" SET "updated_at"=\$1 WHERE "mocks"\."deleted_at" IS NULL AND "id" = \$2`).WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
	(*suite.mock).ExpectQuery(`SELECT \* FROM "mocks" WHERE \("mocks"\."id" = \$1 AND "mocks"\."updated_at" = \$2\) AND "mocks"\."deleted_at" IS NULL AND "mocks"\."id" = \$3 ORDER BY "mocks"\."id" LIMIT \$4`).WithArgs(1, sqlmock.AnyArg(), 1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "attr"}).AddRow(1, 1))
	(*suite.mock).ExpectCommit()

	t := suite.T()
	now := time.Now()
	mock := &mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.Update(suite.ctx, mock)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, mock.UpdatedAt, now, "UpdatedAt should be greater than now")
	assert.Greater(t, mock.UpdatedAt, mock.CreatedAt, "UpdatedAt should be greater than CreatedAt")
	assert.Equal(t, 1, mock.Attr)
}

func (suite *RepositoryTestSuite) TestUpdateFail() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(`UPDATE "mocks" SET "updated_at"=\$1 WHERE "mocks"\."deleted_at" IS NULL AND "id" = \$2`).WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
	(*suite.mock).ExpectQuery(`SELECT \* FROM "mocks" WHERE \("mocks"\."id" = \$1 AND "mocks"\."updated_at" = \$2\) AND "mocks"\."deleted_at" IS NULL AND "mocks"\."id" = \$3 ORDER BY "mocks"\."id" LIMIT \$4`).WithArgs(1, sqlmock.AnyArg(), 1, 1).WillReturnError(sql.ErrNoRows)
	(*suite.mock).ExpectRollback()

	mock := &mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.Update(suite.ctx, mock)
	assert.ErrorIs(suite.T(), err, sql.ErrNoRows)
}

func (suite *RepositoryTestSuite) TestSoftDelete() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(`UPDATE "mocks" SET "deleted_at"=\$1 WHERE "mocks"\."id" = \$2 AND "mocks"\."deleted_at" IS NULL`).WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
	(*suite.mock).ExpectCommit()

	t := suite.T()
	m := mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.SoftDelete(suite.ctx, &m)
	assert.NoError(t, err)
}

func (suite *RepositoryTestSuite) TestSoftDeleteFail() {
	(*suite.mock).ExpectBegin()
	err := suite.Repository.SoftDelete(suite.ctx, &mock{})
	assert.ErrorContains(suite.T(), err, gorm.ErrMissingWhereClause.Error())
}

func (suite *RepositoryTestSuite) TestHardDelete() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(`DELETE FROM "mocks" WHERE "mocks"\."id" = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	(*suite.mock).ExpectCommit()

	t := suite.T()
	m := mock{
		Model: gorm.Model{ID: 1},
	}
	err := suite.Repository.HardDelete(suite.ctx, &m)
	assert.NoError(t, err)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
