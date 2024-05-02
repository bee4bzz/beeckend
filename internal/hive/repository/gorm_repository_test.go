package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gaetanDubuc/beeckend/internal/db"
	dbx "github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	mock       *sqlmock.Sqlmock
	db         *dbx.DB
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
func (suite *RepositoryTestSuite) TearDownTest() {
	// we make sure that all expectations were met
	if err := (*suite.mock).ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RepositoryTestSuite) Test_A_User_Can_Query_Its_Hives() {
	(*suite.mock).ExpectQuery(
		`SELECT .*
		FROM "hives" INNER JOIN user_cheptels 
		ON user_cheptels\."user_id" = \$1 AND user_cheptels\."cheptel_id" = hives\."cheptel_id" 
		WHERE "hives"\."deleted_at" IS NULL`).WithArgs(uint(1)).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "cheptel_id"}).
			AddRow(test.ValidHive.ID, test.ValidHive.Name, test.ValidHive.CheptelID))

	(*suite.mock).ExpectQuery(
		`SELECT .*
		FROM "cheptels" WHERE "cheptels"\."id" = \$1 AND "cheptels"\."deleted_at" IS NULL`).
		WithArgs(test.ValidHive.CheptelID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).
			AddRow(test.ValidCheptel.ID, test.ValidCheptel.Name))

	(*suite.mock).ExpectQuery(
		`SELECT .*
		FROM "hive_notes" WHERE "hive_notes"\."hive_id" = \$1 AND "hive_notes"\."deleted_at" IS NULL`).
		WithArgs(test.ValidHive.ID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "hive_id"}).
			AddRow(test.ValidHiveNote.ID, test.ValidHiveNote.HiveID))

	hives := []entity.Hive{}
	err := suite.Repository.QueryByUser(suite.ctx, &test.ValidUser, &hives)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), hives, 1)
}

func (suite *RepositoryTestSuite) Test_Return_Empty_List_When_The_User_Does_Not_Have_Any_Hive() {
	(*suite.mock).ExpectQuery(
		`SELECT .*
		FROM "hives" INNER JOIN user_cheptels 
		ON user_cheptels\."user_id" = \$1 AND user_cheptels\."cheptel_id" = hives\."cheptel_id" 
		WHERE "hives"\."deleted_at" IS NULL`).WithArgs(uint(1)).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "cheptel_id"}))

	hives := []entity.Hive{}
	err := suite.Repository.QueryByUser(suite.ctx, &test.ValidUser, &hives)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), hives, 0)
}

func (suite *RepositoryTestSuite) Test_Return_Error_When_The_Repo_Fails_To_Query_Hives() {
	(*suite.mock).ExpectQuery(
		`SELECT .*
		FROM "hives" INNER JOIN user_cheptels 
		ON user_cheptels\."user_id" = \$1 AND user_cheptels\."cheptel_id" = hives\."cheptel_id" 
		WHERE "hives"\."deleted_at" IS NULL`).WithArgs(uint(1)).WillReturnError(
		test.ErrMock,
	)

	hives := []entity.Hive{}
	err := suite.Repository.QueryByUser(suite.ctx, &test.ValidUser, &hives)
	assert.ErrorIs(suite.T(), err, test.ErrMock)
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
