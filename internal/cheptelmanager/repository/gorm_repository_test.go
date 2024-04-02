package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gaetanDubuc/beeckend/internal/cheptelmanager/service"
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
	Repository service.Repository
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

func (suite *RepositoryTestSuite) TestGet() {
	(*suite.mock).ExpectQuery(
		`SELECT 
		\"cheptels\"\.\"id\",\"cheptels\"\.\"created_at\",\"cheptels\"\.\"updated_at\",\"cheptels\"\.\"deleted_at\",\"cheptels\"\.\"name\" 
		FROM \"cheptels\" 
		JOIN \"user_cheptels\" ON \"user_cheptels\"\.\"cheptel_id\" = \"cheptels\"\.\"id\" AND \"user_cheptels\"\.\"user_id\" = \$1 
		WHERE \"cheptels\"\.\"deleted_at\" IS NULL AND \"cheptels\"\.\"id\" = \$2`,
	).WithArgs(test.ValidUser.ID, test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(time.Now()))

	testcases := []struct {
		name    string
		cheptel *entity.Cheptel
	}{
		{
			name: "Valid cheptel",
			cheptel: &entity.Cheptel{
				Model: gorm.Model{ID: test.ValidCheptel.ID},
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Get(suite.ctx, &test.ValidUser, tc.cheptel)
			assert.NoError(t, err)
			assert.Equal(t, test.ValidCheptel.ID, tc.cheptel.ID)
		})
	}
}

func (suite *RepositoryTestSuite) TestGetFail() {
	(*suite.mock).ExpectQuery(
		`SELECT \"cheptels\"\.\"id\",\"cheptels\"\.\"created_at\",\"cheptels\"\.\"updated_at\",\"cheptels\"\.\"deleted_at\",\"cheptels\"\.\"name\" 
		FROM \"cheptels\" 
		JOIN \"user_cheptels\" ON \"user_cheptels\"\.\"cheptel_id\" = \"cheptels\"\.\"id\" AND \"user_cheptels\"\.\"user_id\" = \$1 
		WHERE \"cheptels\"\.\"deleted_at\" IS NULL`,
	).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	(*suite.mock).ExpectQuery(
		`SELECT \"cheptels\"\.\"id\",\"cheptels\"\.\"created_at\",\"cheptels\"\.\"updated_at\",\"cheptels\"\.\"deleted_at\",\"cheptels\"\.\"name\" 
		FROM \"cheptels\" 
		JOIN \"user_cheptels\" ON \"user_cheptels\"\.\"cheptel_id\" = \"cheptels\"\.\"id\" AND \"user_cheptels\"\.\"user_id\" = \$1 
		WHERE \"cheptels\"\.\"deleted_at\" IS NULL`,
	).WithArgs(test.ValidUser.ID).WillReturnError(sql.ErrTxDone)

	testcases := []struct {
		name    string
		cheptel *entity.Cheptel
		err     error
	}{
		{
			name:    "unknown cheptel",
			cheptel: &entity.Cheptel{},
			err:     gorm.ErrRecordNotFound,
		},
		{
			name:    "transaction error",
			cheptel: &entity.Cheptel{},
			err:     sql.ErrTxDone,
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Get(suite.ctx, &test.ValidUser, &entity.Cheptel{})
			assert.ErrorIs(suite.T(), err, tc.err)
		})
	}
}

func (suite *RepositoryTestSuite) TestFilterByUser() {
	(*suite.mock).ExpectQuery(
		`SELECT \"hives\".\"id\",\"hives\".\"created_at\",\"hives\".\"updated_at\",\"hives\".\"deleted_at\",\"hives\".\"name\",\"hives\".\"cheptel_id\",\"Cheptel\".\"id\" AS \"Cheptel__id\",\"Cheptel\".\"created_at\" AS \"Cheptel__created_at\",\"Cheptel\".\"updated_at\" AS \"Cheptel__updated_at\",\"Cheptel\".\"deleted_at\" AS \"Cheptel__deleted_at\",\"Cheptel\".\"name\" AS \"Cheptel__name\"
		FROM \"hives\" 
		INNER JOIN user_cheptels ON user_cheptels.\"user_id\" = \$1 AND user_cheptels.\"cheptel_id\" = Cheptel.\"id\"
		LEFT JOIN \"cheptels\" \"Cheptel\" ON \"hives\".\"cheptel_id\" = \"Cheptel\".\"id\" AND \(\"Cheptel\".\"deleted_at\" IS NULL AND \(\"Cheptel\".\"id\" = \$2 AND \"Cheptel\".\"name\" = \$3\)\)
		WHERE \"hives\".\"deleted_at\" IS NULL`,
	).WithArgs(test.ValidUser.ID, test.ValidCheptel.ID, test.ValidCheptel.Name).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	err := suite.Repository.FilterByUserID(suite.ctx, test.ValidUser.ID).Joins("Cheptel", suite.db.Where(test.ValidCheptel)).Find(&entity.Hive{}).Error
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestUpdate() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(
		`UPDATE "users" SET "updated_at"=\$1 WHERE "users"\."deleted_at" IS NULL AND "id" = \$2`,
	).WithArgs(sqlmock.AnyArg(), test.ValidUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	(*suite.mock).ExpectExec(
		`INSERT INTO "user_cheptels" \("user_id","cheptel_id"\) VALUES .* ON CONFLICT DO NOTHING`,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	(*suite.mock).ExpectCommit()

	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(
		`DELETE FROM "user_cheptels" WHERE "user_cheptels"\."user_id" = \$1 AND "user_cheptels"\."cheptel_id" <> \$2`,
	).WithArgs(test.ValidUser.ID, test.ValidCheptel.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	(*suite.mock).ExpectCommit()

	testcases := []struct {
		name    string
		cheptel *entity.Cheptel
	}{
		{
			name: "Valid cheptel",
			cheptel: &entity.Cheptel{
				Model: gorm.Model{ID: test.ValidCheptel.ID},
				Name:  test.ValidCheptel.Name,
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Update(suite.ctx, &test.ValidUser, tc.cheptel)
			assert.NoError(t, err)
		})
	}
}

func (suite *RepositoryTestSuite) TestUpdateFail() {
	testcases := []struct {
		name    string
		cheptel *entity.Cheptel
	}{
		{
			name:    "empty cheptel",
			cheptel: &entity.Cheptel{},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Update(suite.ctx, &test.ValidUser, tc.cheptel)
			assert.Error(t, err)
		})
	}

}

func (suite *RepositoryTestSuite) TestCreate() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(
		`UPDATE "users" SET "updated_at"=\$1 WHERE "users"\."deleted_at" IS NULL AND "id" = \$2`,
	).WithArgs(sqlmock.AnyArg(), test.ValidUser.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	(*suite.mock).ExpectExec(
		`INSERT INTO "user_cheptels" \("user_id","cheptel_id"\) VALUES .* ON CONFLICT DO NOTHING`,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	(*suite.mock).ExpectCommit()

	testcases := []struct {
		name    string
		cheptel *entity.Cheptel
	}{
		{
			name: "Valid cheptel",
			cheptel: &entity.Cheptel{
				Model: gorm.Model{ID: test.ValidCheptel.ID},
				Name:  test.ValidCheptel.Name,
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Create(suite.ctx, &test.ValidUser, tc.cheptel)
			assert.NoError(t, err)
		})
	}
}

func (suite *RepositoryTestSuite) TestCreateFail() {
	testcases := []struct {
		name    string
		cheptel *entity.Cheptel
	}{
		{
			name:    "empty cheptel",
			cheptel: &entity.Cheptel{},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Create(suite.ctx, &test.ValidUser, tc.cheptel)
			assert.Error(t, err)
		})
	}
}

func (suite *RepositoryTestSuite) TestDelete() {
	(*suite.mock).ExpectBegin()
	(*suite.mock).ExpectExec(
		`DELETE FROM "user_cheptels" WHERE "user_cheptels"\."user_id" = \$1 AND "user_cheptels"\."cheptel_id" = \$2`,
	).WithArgs(test.ValidUser.ID, test.ValidCheptel.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	(*suite.mock).ExpectCommit()

	testcases := []struct {
		name    string
		cheptel *entity.Cheptel
	}{
		{
			name: "Valid cheptel",
			cheptel: &entity.Cheptel{
				Model: gorm.Model{ID: test.ValidCheptel.ID},
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.SoftDelete(suite.ctx, &test.ValidUser, tc.cheptel)
			assert.NoError(t, err)
		})
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
