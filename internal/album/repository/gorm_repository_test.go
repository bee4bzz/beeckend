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

func (suite *RepositoryTestSuite) TestQueryCheptelAlbumsByUser() {
	testcases := []struct {
		name string
		user entity.User
		len  int
		fn   func()
	}{
		{
			name: "should find albums",
			user: test.ValidUser,
			len:  1,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "users"
					WHERE "users"."deleted_at" IS NULL AND "users"\."id" = \$1`,
				).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name"}).AddRow(test.ValidUser.ID, test.ValidUser.Email, test.ValidUser.Name))
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "user_cheptels"
					WHERE "user_cheptels"\."user_id" = \$1`,
				).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"cheptel_id", "user_id"}).AddRow(test.ValidCheptel.ID, test.ValidUser.ID))
				(*suite.mock).ExpectQuery(
					`SELECT .* FROM "cheptels" WHERE "cheptels"\."id" = \$1 AND "cheptels"\."deleted_at" IS NULL`,
				).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(test.ValidCheptel.ID, test.ValidCheptel.Name))
				(*suite.mock).ExpectQuery(
					`SELECT .* 
					FROM "albums" 
					WHERE "owner_type" = \$1 AND "albums"\."owner_id" = \$2 AND "albums"\."deleted_at" IS NULL`,
				).WithArgs("cheptels", test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "owner_type", "owner_id"}).AddRow(test.ValidAlbum.ID, test.ValidAlbum.OwnerType, test.ValidAlbum.OwnerID))
			},
		},
		{
			name: "should not find album",
			user: entity.User{Model: gorm.Model{ID: 100}},
			len:  0,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "users"
					WHERE "users"."deleted_at" IS NULL AND "users"\."id" = \$1`,
				).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name"}))
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "user_cheptels"
					WHERE "user_cheptels"\."user_id" = \$1`,
				).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"cheptel_id", "user_id"}))
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.fn()
			albums := []entity.Album{}
			err := suite.Repository.QueryCheptelAlbumsByUser(suite.ctx, &tc.user, &albums)
			assert.NoError(t, err)
			assert.Len(t, albums, tc.len)
		})
	}
}

func (suite *RepositoryTestSuite) TestQueryCheptelAlbumsByUserFail() {
	(*suite.mock).ExpectQuery(
		`SELECT .*
		FROM "users"
		WHERE "users"."deleted_at" IS NULL AND "users"\."id" = \$1`,
	).WithArgs(test.ValidUser.ID).WillReturnError(sql.ErrTxDone)

	albums := []entity.Album{}
	err := suite.Repository.QueryCheptelAlbumsByUser(suite.ctx, &test.ValidUser, &albums)
	assert.ErrorIs(suite.T(), err, sql.ErrTxDone)
}

func (suite *RepositoryTestSuite) TestQueryHiveAlbumsByUser() {
	testcases := []struct {
		name string
		user entity.User
		len  int
		fn   func()
	}{
		{
			name: "should find albums",
			user: test.ValidUser,
			len:  1,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "users"
					WHERE "users"."deleted_at" IS NULL AND "users"\."id" = \$1`,
				).WithArgs(test.ValidUser.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "email", "name"}).AddRow(
						test.ValidUser.ID, test.ValidUser.Email, test.ValidUser.Name))
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "user_cheptels"
					WHERE "user_cheptels"\."user_id" = \$1`,
				).WithArgs(test.ValidUser.ID).WillReturnRows(
					sqlmock.NewRows([]string{"cheptel_id", "user_id"}).AddRow(
						test.ValidCheptel.ID, test.ValidUser.ID))
				(*suite.mock).ExpectQuery(
					`SELECT .* FROM "cheptels" WHERE "cheptels"\."id" = \$1 AND "cheptels"\."deleted_at" IS NULL`,
				).WithArgs(test.ValidCheptel.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "name"}).AddRow(
						test.ValidCheptel.ID, test.ValidCheptel.Name))
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "hives"
					WHERE "hives"\."cheptel_id" = \$1 AND "hives"\."deleted_at" IS NULL`,
				).WithArgs(test.ValidCheptel.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "cheptel_id"}).AddRow(
						test.ValidHive.ID, test.ValidCheptel.ID))
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "hive_notes"
					WHERE "hive_notes"\."hive_id" = \$1 AND "hive_notes"\."deleted_at" IS NULL`,
				).WithArgs(test.ValidHive.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "hive_id"}).AddRow(
						test.ValidHiveNote.ID, test.ValidHive.ID))
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "albums"
					WHERE "owner_type" = \$1 AND "albums"\."owner_id" = \$2 AND "albums"\."deleted_at" IS NULL`,
				).WithArgs("hive_notes", test.ValidHiveNote.ID).WillReturnRows(
					sqlmock.NewRows([]string{"id", "owner_type", "owner_id"}).AddRow(
						test.ValidAlbum.ID, test.ValidAlbum.OwnerType, test.ValidAlbum.OwnerID))
			},
		},
		{
			name: "should not find album",
			user: entity.User{Model: gorm.Model{ID: 100}},
			len:  0,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "users"
					WHERE "users"."deleted_at" IS NULL AND "users"\."id" = \$1`,
				).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name"}))
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "user_cheptels"
					WHERE "user_cheptels"\."user_id" = \$1`,
				).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"cheptel_id", "user_id"}))
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.fn()
			albums := []entity.Album{}
			err := suite.Repository.QueryHiveNoteAlbumsByUser(suite.ctx, &tc.user, &albums)
			assert.NoError(t, err)
			assert.Len(t, albums, tc.len)
		})
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
