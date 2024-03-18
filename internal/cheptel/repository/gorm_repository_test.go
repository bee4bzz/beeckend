package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
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
	Repository *GormRepository
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	mockDb, mock, _ := sqlmock.New()
	suite.mock = &mock
	suite.db = db.NewGormForTest(postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres"},
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
	(*suite.mock).ExpectQuery(
		`SELECT .* 
		FROM "cheptels" 
		JOIN "user_cheptels" ON "user_cheptels"\."cheptel_id" = "cheptels"\."id" AND "user_cheptels"\."user_id" = \$1 
		WHERE "cheptels"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.ValidCheptel.ID))
	(*suite.mock).ExpectQuery(
		`SELECT .* 
		FROM "albums" 
		WHERE "owner_type" = \$1 AND "albums"\."owner_id" = \$2 AND "albums"\."deleted_at" IS NULL`,
	).WithArgs("cheptels", test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	(*suite.mock).ExpectQuery(
		`SELECT .* FROM "hives" WHERE "hives"\."cheptel_id" = \$1 AND "hives"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	(*suite.mock).ExpectQuery(
		`SELECT .* 
		FROM "cheptel_notes" 
		WHERE "cheptel_notes"\."cheptel_id" = \$1 AND "cheptel_notes"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	(*suite.mock).ExpectQuery(
		`SELECT .* 
		FROM "cheptels" 
		JOIN "user_cheptels" ON "user_cheptels"\."cheptel_id" = "cheptels"\."id" AND "user_cheptels"\."user_id" = \$1 
		WHERE "cheptels"\."deleted_at" IS NULL`,
	).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	testcases := []struct {
		name string
		entity.User
		len int
	}{
		{"should find cheptels", test.ValidUser, 1},
		{"should not find cheptel", entity.User{Model: gorm.Model{ID: 100}}, 0},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			cheptels := []entity.Cheptel{}
			err := suite.Repository.QueryByUser(suite.ctx, &tc.User, &cheptels)
			assert.NoError(t, err)
			assert.Len(t, cheptels, tc.len)
		})
	}
}

func (suite *RepositoryTestSuite) TestQueryByUserFail() {
	(*suite.mock).ExpectQuery(
		`SELECT .* FROM "cheptels" 
		JOIN "user_cheptels" ON "user_cheptels"\."cheptel_id" = "cheptels"\."id" AND "user_cheptels"\."user_id" = \$1
		WHERE "cheptels"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidUser.ID).WillReturnError(sql.ErrTxDone)

	cheptels := []entity.Cheptel{}
	err := suite.Repository.QueryByUser(suite.ctx, &test.ValidUser, &cheptels)
	assert.ErrorIs(suite.T(), err, sql.ErrTxDone)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
