package repository

import (
	"bytes"
	"context"
	"database/sql"
	l "log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	dbx "github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	mock       *sqlmock.Sqlmock
	db         *dbx.DB
	pool       *testutils.Pool
	Repository *GormRepository[*testutils.ConnPool, *testutils.Conn]

	logger   *log.Logger
	observer *observer.ObservedLogs
	buffer   *bytes.Buffer
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
		l.New(os.Stdout, "\r\n", l.LstdFlags), // io writer
		logger.Config{},
	))

	logger, obs, b := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.buffer = b

	suite.pool = &testutils.Pool{}
	suite.Repository = NewGormRepository(suite.db, suite.pool, suite.logger)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	// we make sure that all expectations were met
	if err := (*suite.mock).ExpectationsWereMet(); err != nil {
		suite.T().Errorf("there were unfulfilled expectations: %s", err)
	}
	suite.pool.AssertExpectations(suite.T())
	suite.T().Log(suite.buffer)

	suite.buffer.Reset()
	suite.observer.TakeAll()
}

func (suite *RepositoryTestSuite) TestQueryByUser() {
	testcases := []struct {
		name string
		user entity.User
		len  int
		fn   func()
	}{
		{
			name: "should find cheptels",
			user: test.ValidUser,
			len:  1,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT .* 
					FROM "cheptels" 
					JOIN "user_cheptels" ON "user_cheptels"\."cheptel_id" = "cheptels"\."id" AND "user_cheptels"\."user_id" = \$1 
					WHERE "cheptels"\."deleted_at" IS NULL`,
				).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.ValidCheptel.ID))
				(*suite.mock).ExpectQuery(
					`SELECT .* 
					FROM "cheptel_albums" 
					WHERE "cheptel_albums"\."cheptel_id" = \$1 AND "cheptel_albums"\."deleted_at" IS NULL`,
				).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))
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
					FROM "user_cheptels"
					WHERE "user_cheptels"\."cheptel_id" = \$1`,
				).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"cheptel_id", "user_id"}).AddRow(test.ValidCheptel.ID, test.ValidUser.ID))
				(*suite.mock).ExpectQuery(
					`SELECT .*
					FROM "users"
					WHERE "users"\."id" = \$1`,
				).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name"}).AddRow(test.ValidUser.ID, test.ValidUser.Email, test.ValidUser.Name))
			},
		},
		{
			name: "should not find cheptel",
			user: entity.User{Model: gorm.Model{ID: 100}},
			len:  0,
			fn: func() {
				(*suite.mock).ExpectQuery(
					`SELECT .* 
						FROM "cheptels" 
						JOIN "user_cheptels" ON "user_cheptels"\."cheptel_id" = "cheptels"\."id" AND "user_cheptels"\."user_id" = \$1 
						WHERE "cheptels"\."deleted_at" IS NULL`,
				).WithArgs(100).WillReturnRows(sqlmock.NewRows([]string{"id"}))

			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc.fn()
			cheptels := []entity.Cheptel{}
			err := suite.Repository.QueryByUser(suite.ctx, &tc.user, &cheptels)
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

func (suite *RepositoryTestSuite) Test_A_User_Can_Subscribe_To_Cheptels_Modifications() {
	connPool := &testutils.ConnPool{}
	suite.pool.On("Acquire").Return(connPool, nil).Once()
	connPool.On("Exec", "LISTEN Cheptels1").
		Return(pgconn.CommandTag{}, nil).Once()
	connPool.On("Exec", "UNLISTEN Cheptels1").
		Return(pgconn.CommandTag{}, nil).Once()
	connPool.On("Exec", "DROP TRIGGER trigger_Cheptels_user_1 ON cheptels").
		Return(pgconn.CommandTag{}, nil).Once()
	connPool.On("Release").Once()

	connPool.On("Exec", "\n\t\t\tCREATE OR REPLACE TRIGGER trigger_Cheptels_user_1\n\t\t\tAFTER UPDATE OR DELETE ON cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (OLD.id IN ('1'))\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\t\t\t\n\t\t\tCREATE OR REPLACE TRIGGER trigger_UserCheptels_user_1\n\t\t\tAFTER UPDATE OR DELETE ON user_cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (OLD.cheptel_id IN ('1')\n\t\t\tOR OLD.user_id = 1)\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\n\t\t\tCREATE OR REPLACE TRIGGER trigger_UserCheptels_user_1_insert\n\t\t\tAFTER INSERT ON user_cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (NEW.cheptel_id IN ('1')\n\t\t\tOR NEW.user_id = 1)\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\t\t\t").Return(
		pgconn.CommandTag{}, nil).Once()

	conn := &testutils.Conn{}
	connPool.On(
		"Conn",
	).Return(conn).Once()

	conn.On("WaitForNotification").Return(
		&pgconn.Notification{Payload: "'1'"}, nil).Once()

	conn.On("WaitForNotification").Return(
		&pgconn.Notification{}, test.AnError).Once()

	RegisterExpectexQueryByUser(suite.mock)
	RegisterExpectexQueryByUser(suite.mock)

	cheptels := make(chan *[]entity.Cheptel)

	err := suite.Repository.Subscribe(suite.ctx, &test.ValidUser, cheptels)

	ok := assert.NoError(suite.T(), err)
	if ok {
		assert.Len(suite.T(), *<-cheptels, 1)
		assert.Len(suite.T(), *<-cheptels, 1)
		v, ok := <-cheptels
		assert.Nil(suite.T(), v)
		assert.False(suite.T(), ok)
	}
}

func (suite *RepositoryTestSuite) Test_Return_Error_When_A_User_Can_Not_Subscribe_To_Modifications() {
	suite.T().Run("Can_Not_Query_The_First_Cheptels_To_Create_Triggers", func(t *testing.T) {
		connPool := &testutils.ConnPool{}
		suite.pool.On("Acquire").Return(connPool, nil).Once()
		connPool.On("Exec", "LISTEN Cheptels1").
			Return(pgconn.CommandTag{}, nil).Once()

		(*suite.mock).ExpectQuery(
			`SELECT .* 
			FROM "cheptels" 
			JOIN "user_cheptels" ON "user_cheptels"\."cheptel_id" = "cheptels"\."id" AND "user_cheptels"\."user_id" = \$1 
			WHERE "cheptels"\."deleted_at" IS NULL`,
		).WithArgs(test.ValidUser.ID).
			WillReturnError(sql.ErrNoRows)

		cheptels := make(chan *[]entity.Cheptel)

		err := suite.Repository.Subscribe(suite.ctx, &test.ValidUser, cheptels)

		assert.ErrorIs(suite.T(), err, sql.ErrNoRows)
	})

	suite.T().Run("Can_Not_Create_Triggers", func(t *testing.T) {
		connPool := &testutils.ConnPool{}
		suite.pool.On("Acquire").Return(connPool, nil).Once()
		connPool.On("Exec", "LISTEN Cheptels1").
			Return(pgconn.CommandTag{}, nil).Once()

		RegisterExpectexQueryByUser(suite.mock)

		connPool.On("Exec", "\n\t\t\tCREATE OR REPLACE TRIGGER trigger_Cheptels_user_1\n\t\t\tAFTER UPDATE OR DELETE ON cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (OLD.id IN ('1'))\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\t\t\t\n\t\t\tCREATE OR REPLACE TRIGGER trigger_UserCheptels_user_1\n\t\t\tAFTER UPDATE OR DELETE ON user_cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (OLD.cheptel_id IN ('1')\n\t\t\tOR OLD.user_id = 1)\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\n\t\t\tCREATE OR REPLACE TRIGGER trigger_UserCheptels_user_1_insert\n\t\t\tAFTER INSERT ON user_cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (NEW.cheptel_id IN ('1')\n\t\t\tOR NEW.user_id = 1)\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\t\t\t").Return(
			pgconn.CommandTag{}, test.AnError).Once()

		cheptels := make(chan *[]entity.Cheptel)

		err := suite.Repository.Subscribe(suite.ctx, &test.ValidUser, cheptels)

		assert.ErrorIs(suite.T(), err, test.AnError)
	})

	suite.T().Run("Can_Not_Acquire_Connection", func(t *testing.T) {
		connPool := &testutils.ConnPool{}
		suite.pool.On("Acquire").Return(connPool, test.AnError).Once()

		cheptels := make(chan *[]entity.Cheptel)

		err := suite.Repository.Subscribe(suite.ctx, &test.ValidUser, cheptels)

		assert.ErrorIs(suite.T(), err, test.AnError)
	})

	suite.T().Run("Can_Not_Listen_To_Modifications", func(t *testing.T) {
		connPool := &testutils.ConnPool{}
		suite.pool.On("Acquire").Return(connPool, nil).Once()
		connPool.On("Exec", "LISTEN Cheptels1").
			Return(pgconn.CommandTag{}, test.AnError).Once()

		cheptels := make(chan *[]entity.Cheptel)

		err := suite.Repository.Subscribe(suite.ctx, &test.ValidUser, cheptels)

		assert.ErrorIs(suite.T(), err, test.AnError)
	})
}

func (suite *RepositoryTestSuite) Test_Return_No_Error_When_The_Stream_Fails() {
	suite.T().Run("Can_Not_Query_Cheptels", func(t *testing.T) {
		connPool := &testutils.ConnPool{}
		suite.pool.On("Acquire").Return(connPool, nil).Once()
		connPool.On("Exec", "LISTEN Cheptels1").
			Return(pgconn.CommandTag{}, nil).Once()
		connPool.On("Exec", "UNLISTEN Cheptels1").
			Return(pgconn.CommandTag{}, nil).Once()
		connPool.On("Exec", "DROP TRIGGER trigger_Cheptels_user_1 ON cheptels").
			Return(pgconn.CommandTag{}, nil).Once()
		connPool.On("Release").Once()

		connPool.On("Exec", "\n\t\t\tCREATE OR REPLACE TRIGGER trigger_Cheptels_user_1\n\t\t\tAFTER UPDATE OR DELETE ON cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (OLD.id IN ('1'))\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\t\t\t\n\t\t\tCREATE OR REPLACE TRIGGER trigger_UserCheptels_user_1\n\t\t\tAFTER UPDATE OR DELETE ON user_cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (OLD.cheptel_id IN ('1')\n\t\t\tOR OLD.user_id = 1)\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\n\t\t\tCREATE OR REPLACE TRIGGER trigger_UserCheptels_user_1_insert\n\t\t\tAFTER INSERT ON user_cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (NEW.cheptel_id IN ('1')\n\t\t\tOR NEW.user_id = 1)\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\t\t\t").Return(
			pgconn.CommandTag{}, nil).Once()

		conn := &testutils.Conn{}
		connPool.On(
			"Conn",
		).Return(conn).Once()

		conn.On("WaitForNotification").Return(
			&pgconn.Notification{Payload: "'1'"}, nil).Once()

		RegisterExpectexQueryByUser(suite.mock)

		(*suite.mock).ExpectQuery(
			`SELECT .* 
			FROM "cheptels" 
			JOIN "user_cheptels" ON "user_cheptels"\."cheptel_id" = "cheptels"\."id" AND "user_cheptels"\."user_id" = \$1 
			WHERE "cheptels"\."deleted_at" IS NULL`,
		).WithArgs(test.ValidUser.ID).
			WillReturnError(sql.ErrNoRows)

		cheptels := make(chan *[]entity.Cheptel)

		err := suite.Repository.Subscribe(suite.ctx, &test.ValidUser, cheptels)

		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), *<-cheptels, 1)
		v, ok := <-cheptels
		assert.Nil(suite.T(), v)
		assert.False(suite.T(), ok)
	})

	suite.T().Run("Can_Not_Unlisten_Modifications_Or_Drop_Triggers", func(t *testing.T) {
		connPool := &testutils.ConnPool{}
		suite.pool.On("Acquire").Return(connPool, nil).Once()
		connPool.On("Exec", "LISTEN Cheptels1").
			Return(pgconn.CommandTag{}, nil).Once()
		connPool.On("Exec", "UNLISTEN Cheptels1").
			Return(pgconn.CommandTag{}, test.AnError).Once()
		connPool.On("Exec", "DROP TRIGGER trigger_Cheptels_user_1 ON cheptels").
			Return(pgconn.CommandTag{}, test.AnError).Once()
		connPool.On("Release").Once()

		RegisterExpectexQueryByUser(suite.mock)

		connPool.On("Exec", "\n\t\t\tCREATE OR REPLACE TRIGGER trigger_Cheptels_user_1\n\t\t\tAFTER UPDATE OR DELETE ON cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (OLD.id IN ('1'))\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\t\t\t\n\t\t\tCREATE OR REPLACE TRIGGER trigger_UserCheptels_user_1\n\t\t\tAFTER UPDATE OR DELETE ON user_cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (OLD.cheptel_id IN ('1')\n\t\t\tOR OLD.user_id = 1)\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\n\t\t\tCREATE OR REPLACE TRIGGER trigger_UserCheptels_user_1_insert\n\t\t\tAFTER INSERT ON user_cheptels\n\t\t\tFOR EACH ROW\n\t\t\tWHEN (NEW.cheptel_id IN ('1')\n\t\t\tOR NEW.user_id = 1)\n\t\t\tEXECUTE FUNCTION notify_changes(Cheptels1);\n\t\t\t").Return(
			pgconn.CommandTag{}, nil).Once()

		conn := &testutils.Conn{}
		connPool.On(
			"Conn",
		).Return(conn).Once()

		conn.On("WaitForNotification").Return(
			&pgconn.Notification{}, test.AnError).Once()

		cheptels := make(chan *[]entity.Cheptel)

		err := suite.Repository.Subscribe(suite.ctx, &test.ValidUser, cheptels)

		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), *<-cheptels, 1)
		v, ok := <-cheptels
		assert.Nil(suite.T(), v)
		assert.False(suite.T(), ok)
	})
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func RegisterExpectexQueryByUser(mock *sqlmock.Sqlmock) {
	(*mock).ExpectQuery(
		`SELECT .* 
	FROM "cheptels" 
	JOIN "user_cheptels" ON "user_cheptels"\."cheptel_id" = "cheptels"\."id" AND "user_cheptels"\."user_id" = \$1 
	WHERE "cheptels"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(test.ValidCheptel.ID))
	(*mock).ExpectQuery(
		`SELECT .* 
	FROM "cheptel_albums" 
	WHERE "cheptel_albums"\."cheptel_id" = \$1 AND "cheptel_albums"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	(*mock).ExpectQuery(
		`SELECT .* FROM "hives" WHERE "hives"\."cheptel_id" = \$1 AND "hives"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	(*mock).ExpectQuery(
		`SELECT .* 
	FROM "cheptel_notes" 
	WHERE "cheptel_notes"\."cheptel_id" = \$1 AND "cheptel_notes"\."deleted_at" IS NULL`,
	).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"id"}))
	(*mock).ExpectQuery(
		`SELECT .*
	FROM "user_cheptels"
	WHERE "user_cheptels"\."cheptel_id" = \$1`,
	).WithArgs(test.ValidCheptel.ID).WillReturnRows(sqlmock.NewRows([]string{"cheptel_id", "user_id"}).AddRow(test.ValidCheptel.ID, test.ValidUser.ID))
	(*mock).ExpectQuery(
		`SELECT .*
	FROM "users"
	WHERE "users"\."id" = \$1`,
	).WithArgs(test.ValidUser.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name"}).AddRow(test.ValidUser.ID, test.ValidUser.Email, test.ValidUser.Name))
}
