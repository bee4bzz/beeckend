package repository

import (
	"context"
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

func (suite *RepositoryTestSuite) TestGet() {
	(*suite.mock).ExpectQuery(
		`SELECT \"cheptels\".\"id\",\"cheptels\".\"created_at\",\"cheptels\".\"updated_at\",\"cheptels\".\"deleted_at\",\"cheptels\".\"name\" FROM \"cheptels\" JOIN \"user_cheptels\" ON \"user_cheptels\".\"cheptel_id\" = \"cheptels\".\"id\" AND \"user_cheptels\".\"user_id\" = \$1 WHERE \"cheptels\".\"deleted_at\" IS NULL AND "cheptels"."id" = \$2`,
	).WithArgs(test.ValidUser.ID, test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.ValidCheptel.ID))

	(*suite.mock).ExpectQuery(
		`SELECT \"cheptels\".\"id\",\"cheptels\".\"created_at\",\"cheptels\".\"updated_at\",\"cheptels\".\"deleted_at\",\"cheptels\".\"name\" FROM \"cheptels\" JOIN \"user_cheptels\" ON \"user_cheptels\".\"cheptel_id\" = \"cheptels\".\"id\" AND \"user_cheptels\".\"user_id\" = \$1 WHERE \"cheptels\".\"deleted_at\" IS NULL`,
	).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.ValidCheptel.ID))

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
		{
			name:    "Empty cheptel",
			cheptel: &entity.Cheptel{},
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

func (suite *RepositoryTestSuite) TestQuery() {
	(*suite.mock).ExpectQuery(
		`SELECT \"cheptels\".\"id\",\"cheptels\".\"created_at\",\"cheptels\".\"updated_at\",\"cheptels\".\"deleted_at\",\"cheptels\".\"name\" FROM \"cheptels\" JOIN \"user_cheptels\" ON \"user_cheptels\".\"cheptel_id\" = \"cheptels\".\"id\" AND \"user_cheptels\".\"user_id\" = \$1 WHERE \"cheptels\".\"deleted_at\" IS NULL`,
	).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.ValidCheptel.ID))

	(*suite.mock).ExpectQuery(
		`SELECT \"cheptels\".\"id\",\"cheptels\".\"created_at\",\"cheptels\".\"updated_at\",\"cheptels\".\"deleted_at\",\"cheptels\".\"name\" FROM \"cheptels\" JOIN \"user_cheptels\" ON \"user_cheptels\".\"cheptel_id\" = \"cheptels\".\"id\" AND \"user_cheptels\".\"user_id\" = \$1 WHERE \"cheptels\".\"deleted_at\" IS NULL`,
	).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.ValidCheptel.ID))

	testcases := []struct {
		name     string
		cheptels *[]entity.Cheptel
	}{
		{
			name: "Valid cheptel",
			cheptels: &[]entity.Cheptel{
				{Model: gorm.Model{ID: test.ValidCheptel.ID}},
			},
		},
		{
			name:     "Empty cheptels",
			cheptels: &[]entity.Cheptel{},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Query(suite.ctx, &test.ValidUser, tc.cheptels)
			assert.NoError(t, err)
			assert.Equal(t, test.ValidCheptel.ID, (*tc.cheptels)[0].ID)
		})
	}
}

func (suite *RepositoryTestSuite) TestCreate() {
	(*suite.mock).ExpectExec(
		`INSERT INTO "user_cheptels" \("user_id","cheptel_id"\) VALUES \(\$1,\$2\)`,
	).WithArgs(test.ValidUser.ID, test.ValidCheptel.ID).WillReturnResult(sqlmock.NewResult(1, 1))

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
		{
			name:    "Empty cheptel",
			cheptel: &entity.Cheptel{},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Repository.Create(suite.ctx, &test.ValidUser, tc.cheptel)
			assert.NoError(t, err)
		})
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
