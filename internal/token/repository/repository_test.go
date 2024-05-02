package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/errors"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/test"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/testutils"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/dbcontext"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	mock       *sqlmock.Sqlmock
	db         *dbcontext.DB
	Repository *Repository
}

// this function executes before the test suite begins execution.
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, _, _ := log.NewForTest()
	mockDb, mock, _ := sqlmock.New()
	suite.mock = &mock
	suite.db = dbcontext.New(dbx.NewFromDB(mockDb, "postgres"))
	suite.Repository = NewRepository(suite.db, logger)
}

// this function executes after all tests executed.
func (suite *RepositoryTestSuite) TearDownSuite() {
	// we make sure that all expectations were met
	if err := (*suite.mock).ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
}

func (suite *RepositoryTestSuite) TestQuery() {
	testcases := []struct {
		name    string
		filters []map[string]any
		len     int
		fn      func()
	}{
		{
			name:    "valid token",
			filters: []map[string]any{{"UUID": test.DeviceSecurityToken.UUID}},
			len:     1,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnRows(
					sqlmock.NewRows([]string{"UUID"}).AddRow(
						test.DeviceSecurityToken.UUID,
					),
				)
			},
		},

		{
			name:    "valid multiple filters",
			filters: []map[string]any{{"UUID": test.DeviceSecurityToken.UUID}, {"owner_UUID": test.DeviceSecurityToken.OwnerUUID}},
			len:     1,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE \("UUID"=\$1\) OR \("owner_UUID"=\$2\)`).WithArgs(
					test.DeviceSecurityToken.UUID, test.DeviceSecurityToken.OwnerUUID).WillReturnRows(
					sqlmock.NewRows([]string{"UUID"}).AddRow(
						test.DeviceSecurityToken.UUID,
					),
				)
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.fn()
			tokens := []entity.Token{}
			err := suite.Repository.Query(suite.ctx, entity.DeviceSecurityType, tc.filters, &tokens)
			assert.NoError(suite.T(), err)
			assert.Len(suite.T(), tokens, tc.len)
		})
	}
}

func (suite *RepositoryTestSuite) TestQueryFail() {
	(*suite.mock).ExpectQuery(
		`SELECT \* FROM "devicesecurity" 
		WHERE "UUID"=\$1`).WithArgs(
		test.DeviceSecurityToken.UUID).WillReturnError(
		testutils.ErrMocked,
	)

	tokens := []entity.Token{}
	err := suite.Repository.Query(suite.ctx, entity.DeviceSecurityType, []map[string]any{{"UUID": test.DeviceSecurityToken.UUID}}, &tokens)
	assert.ErrorIs(suite.T(), err, testutils.ErrMocked)
	assert.Len(suite.T(), tokens, 0)
}

func (suite *RepositoryTestSuite) TestGetBy() {
	(*suite.mock).ExpectQuery(
		`SELECT \* FROM "devicesecurity" 
		WHERE "UUID"=\$1`).WithArgs(
		test.DeviceSecurityToken.UUID).WillReturnRows(
		sqlmock.NewRows([]string{"UUID"}).AddRow(
			test.DeviceSecurityToken.UUID,
		),
	)

	token := entity.Token{
		OwnerType: entity.DeviceSecurityType,
	}
	err := suite.Repository.GetBy(suite.ctx, map[string]any{"UUID": test.DeviceSecurityToken.UUID}, &token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), test.DeviceSecurityToken.UUID, token.UUID)
}

func (suite *RepositoryTestSuite) TestGetByFail() {
	(*suite.mock).ExpectQuery(
		`SELECT \* FROM "devicesecurity" 
		WHERE "UUID"=\$1`).WithArgs(
		test.DeviceSecurityToken.UUID).WillReturnError(
		testutils.ErrMocked,
	)

	token := entity.Token{
		OwnerType: entity.DeviceSecurityType,
	}
	err := suite.Repository.GetBy(suite.ctx, map[string]any{"UUID": test.DeviceSecurityToken.UUID}, &token)
	assert.ErrorIs(suite.T(), err, testutils.ErrMocked)
	assert.Equal(suite.T(), entity.Token{OwnerType: entity.DeviceSecurityType}, token)
}

func (suite *RepositoryTestSuite) TestCreate() {
	(*suite.mock).ExpectQuery(
		`SELECT \* FROM "devicesecurity" 
		WHERE "UUID"=\$1`).WithArgs(
		test.DeviceSecurityToken.UUID).WillReturnError(
		sql.ErrNoRows,
	)
	(*suite.mock).ExpectExec(
		`INSERT INTO "devicesecurity" \("UUID", "created_at", "hashed_token", "owner_UUID", "updated_at"\) 
		VALUES \(\$1, \$2, \$3, \$4, \$5\)`).WithArgs(
		test.DeviceSecurityToken.UUID,
		test.DeviceSecurityToken.CreatedAt,
		test.DeviceSecurityToken.HashedToken,
		test.DeviceSecurityToken.OwnerUUID,
		test.DeviceSecurityToken.UpdatedAt,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)
	(*suite.mock).ExpectQuery(
		`SELECT \* FROM "devicesecurity" 
		WHERE "UUID"=\$1`).WithArgs(
		test.DeviceSecurityToken.UUID).WillReturnRows(
		sqlmock.NewRows([]string{"UUID"}).AddRow(
			test.DeviceSecurityToken.UUID,
		))

	err := suite.Repository.Create(suite.ctx, &test.DeviceSecurityToken)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestCreateFail() {
	testcases := []struct {
		name string
		err  error
		fn   func()
	}{
		{
			name: "should fail to create a token with the same UUID",
			err:  errors.ErrDuplicate,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnRows(
					sqlmock.NewRows([]string{"UUID"}).AddRow(
						test.DeviceSecurityToken.UUID,
					),
				)
			},
		},
		{
			name: "should fail to create a token with an error",
			err:  testutils.ErrMocked,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnError(
					testutils.ErrMocked,
				)
			},
		},
		{
			name: "fail to insert a token",
			err:  testutils.ErrMocked,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnError(
					sql.ErrNoRows,
				)
				(*suite.mock).ExpectExec(
					`INSERT INTO "devicesecurity" \("UUID", "created_at", "hashed_token", "owner_UUID", "updated_at"\) 
					VALUES \(\$1, \$2, \$3, \$4, \$5\)`).WithArgs(
					test.DeviceSecurityToken.UUID,
					test.DeviceSecurityToken.CreatedAt,
					test.DeviceSecurityToken.HashedToken,
					test.DeviceSecurityToken.OwnerUUID,
					test.DeviceSecurityToken.UpdatedAt,
				).WillReturnError(
					testutils.ErrMocked,
				)
			},
		},
		{
			name: "fail to get the created token",
			err:  testutils.ErrMocked,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnError(
					sql.ErrNoRows,
				)
				(*suite.mock).ExpectExec(
					`INSERT INTO "devicesecurity" \("UUID", "created_at", "hashed_token", "owner_UUID", "updated_at"\) 
					VALUES \(\$1, \$2, \$3, \$4, \$5\)`).WithArgs(
					test.DeviceSecurityToken.UUID,
					test.DeviceSecurityToken.CreatedAt,
					test.DeviceSecurityToken.HashedToken,
					test.DeviceSecurityToken.OwnerUUID,
					test.DeviceSecurityToken.UpdatedAt,
				).WillReturnResult(
					sqlmock.NewResult(1, 1),
				)
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnError(
					testutils.ErrMocked,
				)
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.fn()
			err := suite.Repository.Create(suite.ctx, &test.DeviceSecurityToken)
			assert.ErrorIs(suite.T(), err, tc.err)
		})
	}
}

func (suite *RepositoryTestSuite) TestUpdate() {
	(*suite.mock).ExpectExec(
		`UPDATE "devicesecurity" 
		SET "created_at"=\$1, "hashed_token"=\$2, "owner_UUID"=\$3, "updated_at"=\$4 
		WHERE "UUID"=\$5`).WithArgs(
		test.DeviceSecurityToken.CreatedAt,
		test.DeviceSecurityToken.HashedToken,
		test.DeviceSecurityToken.OwnerUUID,
		test.DeviceSecurityToken.UpdatedAt,
		test.DeviceSecurityToken.UUID,
	).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	(*suite.mock).ExpectQuery(
		`SELECT \* FROM "devicesecurity" 
		WHERE "UUID"=\$1`).WithArgs(
		test.DeviceSecurityToken.UUID).WillReturnRows(
		sqlmock.NewRows([]string{"UUID"}).AddRow(
			test.DeviceSecurityToken.UUID,
		))

	err := suite.Repository.Update(suite.ctx, &test.DeviceSecurityToken)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestUpdateFail() {
	testcases := []struct {
		name string
		err  error
		fn   func()
	}{
		{
			name: "should fail to update the token",
			err:  testutils.ErrMocked,
			fn: func() {
				(*suite.mock).ExpectExec(
					`UPDATE "devicesecurity" 
					SET "created_at"=\$1, "hashed_token"=\$2, "owner_UUID"=\$3, "updated_at"=\$4 
					WHERE "UUID"=\$5`).WithArgs(
					test.DeviceSecurityToken.CreatedAt,
					test.DeviceSecurityToken.HashedToken,
					test.DeviceSecurityToken.OwnerUUID,
					test.DeviceSecurityToken.UpdatedAt,
					test.DeviceSecurityToken.UUID,
				).WillReturnError(
					testutils.ErrMocked,
				)
			},
		},
		{
			name: "should fail to get the updated token",
			err:  testutils.ErrMocked,
			fn: func() {
				(*suite.mock).ExpectExec(
					`UPDATE "devicesecurity" 
					SET "created_at"=\$1, "hashed_token"=\$2, "owner_UUID"=\$3, "updated_at"=\$4 
					WHERE "UUID"=\$5`).WithArgs(
					test.DeviceSecurityToken.CreatedAt,
					test.DeviceSecurityToken.HashedToken,
					test.DeviceSecurityToken.OwnerUUID,
					test.DeviceSecurityToken.UpdatedAt,
					test.DeviceSecurityToken.UUID,
				).WillReturnResult(
					sqlmock.NewResult(1, 1),
				)

				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnError(
					testutils.ErrMocked,
				)
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.fn()
			err := suite.Repository.Update(suite.ctx, &test.DeviceSecurityToken)
			assert.ErrorIs(suite.T(), err, testutils.ErrMocked)
		})
	}
}

func (suite *RepositoryTestSuite) TestDelete() {
	(*suite.mock).ExpectQuery(
		`SELECT \* FROM "devicesecurity" 
		WHERE "UUID"=\$1`).WithArgs(
		test.DeviceSecurityToken.UUID).WillReturnRows(
		sqlmock.NewRows([]string{"UUID"}).AddRow(
			test.DeviceSecurityToken.UUID,
		))
	(*suite.mock).ExpectExec(
		`DELETE FROM "devicesecurity" 
		WHERE "UUID"=\$1`).WithArgs(
		test.DeviceSecurityToken.UUID).WillReturnResult(
		sqlmock.NewResult(1, 1),
	)

	err := suite.Repository.Delete(suite.ctx, &test.DeviceSecurityToken)
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestDeleteFail() {
	testcases := []struct {
		name string
		err  error
		fn   func()
	}{
		{
			name: "should fail to delete the token",
			err:  testutils.ErrMocked,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnRows(
					sqlmock.NewRows([]string{"UUID"}).AddRow(
						test.DeviceSecurityToken.UUID,
					))
				(*suite.mock).ExpectExec(
					`DELETE FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnError(
					testutils.ErrMocked,
				)
			},
		},
		{
			name: "should fail to get the deleted token",
			err:  testutils.ErrMocked,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT \* FROM "devicesecurity" 
					WHERE "UUID"=\$1`).WithArgs(
					test.DeviceSecurityToken.UUID).WillReturnError(
					testutils.ErrMocked,
				)
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.fn()
			err := suite.Repository.Delete(suite.ctx, &test.DeviceSecurityToken)
			assert.ErrorIs(suite.T(), err, tc.err)
		})
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
