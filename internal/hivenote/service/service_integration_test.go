package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/hivenote/repository"
	"github.com/gaetanDubuc/beeckend/internal/hivenote/schema"
	"github.com/gaetanDubuc/beeckend/internal/hivenote/testutils"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	cheptelmngtestutils "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/testutils"
	hiverepository "github.com/gaetanDubuc/beeckend/internal/hive/repository"
)

const (
	dbName = "hive.db"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx            context.Context
	db             *db.DB
	Service        *Service
	CheptelManager *cheptelmngtestutils.CheptelManager
	logger         *log.Logger
	observer       *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = db.NewGormForTest(sqlite.Open(dbName))
	suite.CheptelManager = &cheptelmngtestutils.CheptelManager{}
	logger, obs, _ := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(
		repository.NewGormRepository(suite.db),
		hiverepository.NewGormRepository(suite.db),
		suite.CheptelManager,
		suite.logger)
}

// this function executes after all tests executed
func (suite *RepositoryIntegrationSuite) TearDownSuite() {
	suite.CheptelManager.AssertExpectations(suite.T())
	if err := os.Remove(dbName); err != nil {
		suite.T().Fatalf("Error while deleting the database file: %s", err)
	}
}

func (suite *RepositoryIntegrationSuite) SetupTest() {
	db.Seed(suite.T(), suite.db,
		&test.ValidUser,
		&test.ValidCheptel,
		&test.ValidHive,
		&test.ValidHiveNote)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	db.Clean(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
	now := time.Now()

	hive, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{
		UserID:     test.ValidUser.ID,
		CheptelID:  test.ValidHive.CheptelID,
		HiveID:     test.ValidHive.ID,
		HiveNoteID: test.ValidHiveNote.ID,
		NewName:    "new name"})

	assert.NoError(suite.T(), err)
	testutils.AssertHiveNoteUpdated(suite.T(), entity.HiveNote{Model: gorm.Model{
		ID: test.ValidHiveNote.ID,
	},
		HiveID:    test.ValidHive.ID,
		Name:      "new name",
		Operation: test.ValidHiveNote.Operation,
	}, hive, now)
}

func (suite *RepositoryIntegrationSuite) TestUpdateFail() {
	validUpdateReq := schema.UpdateRequest{
		UserID:     test.ValidUser.ID,
		CheptelID:  test.ValidHive.CheptelID,
		HiveID:     test.ValidHive.ID,
		HiveNoteID: test.ValidHiveNote.ID,
	}
	hiveNoteNotFoundReq := validUpdateReq.CopyWith(schema.UpdateRequest{
		HiveNoteID: 100,
	})

	unknownHiveReq := validUpdateReq.CopyWith(schema.UpdateRequest{
		HiveID: 100,
	})

	testcases := []struct {
		name string
		req  schema.UpdateRequest
		err  error
		fn   func()
	}{
		{
			name: "hive note should be not found",
			req:  hiveNoteNotFoundReq,
			err:  gorm.ErrRecordNotFound,
			fn: func() {
				suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
			},
		},
		{
			name: "A known user of the cheptel should not be able to update the hive note by an unkown hive",
			req:  unknownHiveReq,
			err:  gorm.ErrRecordNotFound,
			fn: func() {
				suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
			},
		},
		{
			name: "An unknown user of the current cheptel should not be able to update the hive note",
			req:  validUpdateReq,
			err:  test.ErrMock,
			fn: func() {
				suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()
			},
		},
		{name: "the request should be invalid", req: schema.UpdateRequest{}},
		{name: "the request should be invalid", req: validUpdateReq.CopyWith(schema.UpdateRequest{
			NewNBRisers: 100,
		}),
			fn: func() {
				suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
			}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn()
			}
			hive, err := suite.Service.Update(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, hive)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestCreate() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptel.ID, test.ValidUser.ID).Return(nil).Once()
	now := time.Now()

	hive, err := suite.Service.Create(suite.ctx, schema.CreateRequest{
		UserID:     test.ValidUser.ID,
		CheptelID:  test.ValidCheptel.ID,
		HiveID:     test.ValidHive.ID,
		HiveNoteID: 100,
		Name:       "new name",
		Operation:  "new operation",
	})

	assert.NoError(suite.T(), err)
	testutils.AssertHiveNoteCreated(suite.T(), entity.HiveNote{Model: gorm.Model{
		ID: 100,
	},
		HiveID:    test.ValidHive.ID,
		Name:      "new name",
		Operation: "new operation",
	}, hive, now)
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	validCreateReq := schema.CreateRequest{
		UserID:     test.ValidUser.ID,
		CheptelID:  test.ValidHive.CheptelID,
		HiveID:     test.ValidHive.ID,
		HiveNoteID: test.ValidHiveNote.ID,
		Name:       "new name",
	}

	invalidCreateReq := validCreateReq.CopyWith(
		schema.CreateRequest{
			Name: "",
		},
	)

	testcases := []struct {
		name string
		req  schema.CreateRequest
		fn   func()
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to update the hive note", req: validCreateReq, err: test.ErrMock, fn: func() {
			suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()
		}},
		{name: "Create a hive note without a name should return an error", req: invalidCreateReq, fn: func() {
			suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
		}},
		{
			name: "hive should not be found",
			req: validCreateReq.CopyWith(
				schema.CreateRequest{
					HiveID: 100,
				},
			),
			err: gorm.ErrRecordNotFound,
			fn: func() {
				suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
			},
		},
		{name: "the request should be invalid", req: schema.CreateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			if tc.fn != nil {
				tc.fn()
			}
			hive, err := suite.Service.Create(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Empty(t, hive)
		})
	}
}

func (suite *RepositoryIntegrationSuite) TestQueryByUser() {
	hivenotes, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{
		UserID: test.ValidUser.ID,
	})

	assert.NoError(suite.T(), err)
	testutils.AssertHiveNotes(suite.T(), []entity.HiveNote{test.ValidHiveNote}, hivenotes)

	assert.Equal(suite.T(), 2, suite.observer.Len())
}

func (suite *RepositoryIntegrationSuite) TestQueryByUserFail() {
	hives, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), hives)
}

func (suite *RepositoryIntegrationSuite) TestDelete() {
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()
	req := schema.Request{
		UserID:     test.ValidUser.ID,
		CheptelID:  test.ValidHive.CheptelID,
		HiveID:     test.ValidHive.ID,
		HiveNoteID: test.ValidHiveNote.ID,
	}
	err := suite.Service.SoftDelete(suite.ctx, req)
	assert.NoError(suite.T(), err)
	err = suite.db.DB().Model(&test.ValidHiveNote).First(&entity.HiveNote{}).Error
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestDeleteFail() {
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidHive.CheptelID, test.ValidUser.ID).Return(nil).Once()

	validReq := schema.Request{
		UserID:     test.ValidUser.ID,
		CheptelID:  test.ValidHive.CheptelID,
		HiveID:     test.ValidHive.ID,
		HiveNoteID: test.ValidHiveNote.ID,
	}

	testcases := []struct {
		name string
		req  schema.Request
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to delete the hive note", req: validReq, err: test.ErrMock},
		{name: "hive should not be found", req: validReq.CopyWith(schema.Request{HiveID: 100}), err: gorm.ErrRecordNotFound},
		{name: "the request should be invalid", req: schema.Request{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
			err := suite.Service.SoftDelete(suite.ctx, tc.req)
			if tc.err == nil {
				assert.Error(t, err)
			} else {
				assert.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func TestRepositoryIntegrationSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationSuite))
}
