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
	testcases := []struct {
		name string
		user entity.User
		len  int
		fn   func()
	}{
		{
			name: "valid user",
			user: test.ValidUser,
			len:  1,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "users" WHERE "users"\."deleted_at" IS NULL AND "users"\."id" = \$1`,
				).WithArgs(test.ValidUser.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "email", "name"}).AddRow(
						test.ValidUser.ID, test.ValidUser.Email, test.ValidUser.Name),
				)
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "user_cheptels" WHERE "user_cheptels"\."user_id" = \$1`,
				).WithArgs(test.ValidUser.ID).WillReturnRows(
					sqlmock.NewRows([]string{"user_id", "cheptel_id"}).AddRow(
						test.ValidUser.ID, test.ValidCheptel.ID),
				)
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "cheptels" WHERE "cheptels"\."id" = \$1`,
				).WithArgs(test.ValidCheptel.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "name"}).AddRow(
						test.ValidCheptel.ID, test.ValidCheptel.Name),
				)
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "cheptel_notes" WHERE "cheptel_notes"\."cheptel_id" = \$1`,
				).WithArgs(test.ValidCheptel.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "name", "cheptel_id"}).AddRow(
						test.ValidCheptelNote.ID, test.ValidCheptelNote.Name, test.ValidCheptel.ID),
				)
			},
		},
		{
			name: "no cheptel note",
			user: entity.User{Model: gorm.Model{ID: 100}},
			len:  0,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "users" WHERE "users"\."deleted_at" IS NULL AND "users"\."id" = \$1`,
				).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name"}))
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "user_cheptels" WHERE "user_cheptels"\."user_id" = \$1`,
				).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"user_id", "cheptel_id"}))
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.fn()
			cheptelNotes := []entity.CheptelNote{}
			err := suite.Repository.QueryByUser(suite.ctx, &tc.user, &cheptelNotes)
			assert.NoError(suite.T(), err)
			assert.Len(suite.T(), cheptelNotes, tc.len)
		})
	}
}

func (suite *RepositoryTestSuite) TestQueryByUserFail() {
	(*suite.mock).ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."deleted_at" IS NULL AND "users"\."id" = \$1`).WithArgs(test.ValidUser.ID).WillReturnError(sql.ErrTxDone)

	cheptelNotes := []entity.CheptelNote{}
	err := suite.Repository.QueryByUser(suite.ctx, &test.ValidUser, &cheptelNotes)
	assert.ErrorIs(suite.T(), err, sql.ErrTxDone)
}

func (suite *RepositoryTestSuite) TestSoftDelete() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(
		`UPDATE "cheptel_notes" SET "deleted_at"=\$1 WHERE "cheptel_notes"\."id" = \$2 AND "cheptel_notes"\."deleted_at" IS NULL`,
	).WithArgs(sqlmock.AnyArg(), test.ValidCheptelNote.ID).WillReturnResult(sqlmock.NewResult(0, 1))
	(*suite.mock).ExpectCommit()

	err := suite.Repository.SoftDelete(suite.ctx, &entity.CheptelNote{Model: gorm.Model{ID: test.ValidCheptelNote.ID}})
	assert.NoError(suite.T(), err)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
