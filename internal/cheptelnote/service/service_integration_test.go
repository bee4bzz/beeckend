package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/cheptelnote/repository"
	"github.com/gaetanDubuc/beeckend/internal/cheptelnote/schema"
	"github.com/gaetanDubuc/beeckend/internal/cheptelnote/testutils"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	chepteltestutils "github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
)

const (
	dbName = "hive.db"
)

type RepositoryIntegrationSuite struct {
	suite.Suite
	ctx            context.Context
	db             *gorm.DB
	Service        *Service
	CheptelManager *chepteltestutils.CheptelManager
	logger         *log.Logger
	observer       *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.db = db.NewGormForTest(sqlite.Open(dbName))
	suite.CheptelManager = &chepteltestutils.CheptelManager{}
	logger, obs := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(repository.NewGormRepository(suite.db), suite.CheptelManager, suite.logger)
}

// this function executes after all tests executed
func (suite *RepositoryIntegrationSuite) TearDownSuite() {
	if err := os.Remove(dbName); err != nil {
		suite.T().Fatalf("Error while deleting the database file: %s", err)
	}
}

func (suite *RepositoryIntegrationSuite) SetupTest() {
	db.Seed(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TearDownTest() {
	db.Clean(suite.T(), suite.db)
}

func (suite *RepositoryIntegrationSuite) TestUpdate() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID, test.ValidUser.ID).Return(nil).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID+1, test.ValidUser.ID).Return(nil).Once()
	now := time.Now()

	hive, err := suite.Service.Update(suite.ctx, schema.UpdateRequest{
		UserID:       test.ValidUser.ID,
		CheptelID:    test.ValidCheptelNote.CheptelID,
		NoteID:       test.ValidCheptelNote.ID,
		NewCheptelID: test.ValidCheptelNote.CheptelID + 1,
		NewName:      "new name"})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptelNoteUpdated(suite.T(), entity.CheptelNote{Model: gorm.Model{
		ID: test.ValidCheptelNote.ID,
	},
		CheptelID: 10,
		Name:      "new name",
		Weather:   entity.CLOUDY,
		Flora:     "new flora",
	}, hive, now)
}

func (suite *RepositoryIntegrationSuite) TestUpdateFail() {
	// hive should be not found
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID, test.ValidUser.ID).Return(nil).Once()

	// An unknown user of the new cheptel should not be able to update the hive
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID, test.ValidUser.ID).Return(nil).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID+1, test.ValidUser.ID).Return(test.ErrMock).Once()

	// An unknown user of the current cheptel should not be able to update the hive
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()

	validUpdateReq := schema.UpdateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptelNote.CheptelID,
		NoteID:    test.ValidCheptelNote.ID,
	}
	hiveNotFoundReq := validUpdateReq.CopyWith(schema.UpdateRequest{
		NoteID: 100,
	})

	newCheptelReq := validUpdateReq.CopyWith(schema.UpdateRequest{
		NewCheptelID: test.ValidCheptelNote.CheptelID + 1,
	})

	testcases := []struct {
		name string
		req  schema.UpdateRequest
		err  error
	}{
		{name: "hive should be not found", req: hiveNotFoundReq, err: gorm.ErrRecordNotFound},
		{name: "An unknown user of the new cheptel should not be able to update the hive", req: newCheptelReq, err: test.ErrMock},
		{name: "An unknown user of the current cheptel should not be able to update the hive", req: validUpdateReq, err: test.ErrMock},
		{name: "the request should be invalid", req: schema.UpdateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
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
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptel.ID,
		NoteID:    100,
		Name:      "new name",
		Weather:   entity.CLOUDY,
		Flora:     "new flora",
	})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptelNoteCreated(suite.T(), entity.CheptelNote{Model: gorm.Model{
		ID: 100,
	},
		CheptelID: test.ValidCheptel.ID,
		Name:      "new name",
		Weather:   entity.CLOUDY,
		Flora:     "new flora",
	}, hive, now)
}

func (suite *RepositoryIntegrationSuite) TestCreateFail() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID, test.ValidUser.ID).Return(nil).Once()

	validCreateReq := schema.CreateRequest{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptelNote.CheptelID,
		NoteID:    test.ValidCheptelNote.ID,
		Name:      "new name",
	}

	invalidCreateReq := validCreateReq.CopyWith(
		schema.CreateRequest{
			Name: "",
		},
	)

	testcases := []struct {
		name string
		req  schema.CreateRequest
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to update the hive", req: validCreateReq, err: test.ErrMock},
		{name: "Create an hive without a name should return an error", req: invalidCreateReq},
		{name: "the request should be invalid", req: schema.CreateRequest{}},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.name, func(t *testing.T) {
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
	cheptelNotes, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{
		UserID: test.ValidUser.ID,
	})

	assert.NoError(suite.T(), err)
	testutils.AssertCheptelNotes(suite.T(), []entity.CheptelNote{test.ValidCheptelNote}, cheptelNotes)
	assert.Equal(suite.T(), 2, suite.observer.Len())
}

func (suite *RepositoryIntegrationSuite) TestQueryByUserFail() {
	cheptelNotes, err := suite.Service.QueryByUser(suite.ctx, schema.QueryRequest{})
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), cheptelNotes)
}

func (suite *RepositoryIntegrationSuite) TestDelete() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID, test.ValidUser.ID).Return(nil).Once()
	req := schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptelNote.CheptelID,
		NoteID:    test.ValidCheptelNote.ID,
	}
	err := suite.Service.SoftDelete(suite.ctx, req)
	assert.NoError(suite.T(), err)
	err = suite.db.Model(&entity.CheptelNote{}).First(test.ValidCheptelNote).Error
	assert.ErrorIs(suite.T(), err, gorm.ErrRecordNotFound)
}

func (suite *RepositoryIntegrationSuite) TestDeleteFail() {
	suite.CheptelManager.On("OnlyMember", test.ValidCheptelNote.CheptelID, test.ValidUser.ID).Return(test.ErrMock).Once()

	validReq := schema.Request{
		UserID:    test.ValidUser.ID,
		CheptelID: test.ValidCheptelNote.CheptelID,
		NoteID:    test.ValidCheptelNote.ID,
	}

	testcases := []struct {
		name string
		req  schema.Request
		err  error
	}{
		{name: "An unknown user of the current cheptel should not be able to delete the hive", req: validReq, err: test.ErrMock},
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
