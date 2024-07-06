package service

import (
	"context"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/cheptelalbum/schema"
	"github.com/gaetanDubuc/beeckend/internal/cheptelalbum/testutils"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"

	chepteltestutils "github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	cheptelalbumtestutils "github.com/gaetanDubuc/beeckend/internal/cheptelalbum/testutils"
	cheptelmngtestutils "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/testutils"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx               context.Context
	Service           *Service
	CheptelManager    *cheptelmngtestutils.CheptelManager
	cheptelRepository *chepteltestutils.Repository
	Repository        *cheptelalbumtestutils.Repository
	logger            *log.Logger
	observer          *observer.ObservedLogs
}

// this function executes before the test suite begins execution
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.CheptelManager = &cheptelmngtestutils.CheptelManager{}
	suite.cheptelRepository = &chepteltestutils.Repository{}
	suite.Repository = &cheptelalbumtestutils.Repository{}
	logger, obs, _ := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.Service = NewService(suite.Repository, suite.cheptelRepository, suite.CheptelManager, logger)
}

func (suite *RepositoryTestSuite) TestQueryByUserFail() {
	testcases := test.ServiceTestCases[schema.QueryRequest, []entity.CheptelAlbum]{
		{
			Name: "fail to query cheptels by user",
			Req: schema.QueryRequest{
				UserID: test.ValidUser.ID,
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.cheptelRepository.On("QueryByUser", entity.User{
					Model: gorm.Model{
						ID: test.ValidUser.ID,
					},
				}, []entity.Cheptel{}).Return([]entity.Cheptel{test.ValidCheptel}, test.AnError).Once()
			},
		},
		{
			Name: "fail to query cheptel albums by owner ids",
			Req: schema.QueryRequest{
				UserID: test.ValidUser.ID,
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.cheptelRepository.On("QueryByUser", entity.User{
					Model: gorm.Model{
						ID: test.ValidUser.ID,
					},
				}, []entity.Cheptel{}).Return([]entity.Cheptel{test.ValidCheptel}, nil).Once()
				suite.Repository.On("QueryByOwnerIDs", &[]entity.CheptelAlbum{}, []uint{test.ValidCheptel.ID}).
					Return([]entity.CheptelAlbum{test.ValidCheptelAlbum}, test.AnError).Once()
			},
		},
		{
			Name: "Invalid request",
			Req: schema.QueryRequest{
				UserID: 0,
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.Name, func(t *testing.T) {
			tc.ShouldFailAndBeEmpty(suite.T(), suite.Service.QueryByUser)
		})
	}
}

func (suite *RepositoryTestSuite) TestQueryByUserSuccess() {
	testcases := test.ServiceTestCases[schema.QueryRequest, []entity.CheptelAlbum]{
		{
			Name: "succeed",
			Req: schema.QueryRequest{
				UserID: test.ValidUser.ID,
			},
			RegisterMocks: func() {
				suite.cheptelRepository.On("QueryByUser", entity.User{
					Model: gorm.Model{
						ID: test.ValidUser.ID,
					},
				}, []entity.Cheptel{}).Return([]entity.Cheptel{test.ValidCheptel}, nil).Once()
				suite.Repository.On("QueryByOwnerIDs", &[]entity.CheptelAlbum{}, []uint{test.ValidCheptel.ID}).
					Return([]entity.CheptelAlbum{test.ValidCheptelAlbum}, nil).Once()
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.Name, func(t *testing.T) {
			albums := tc.ShouldSucceedAndNotEmpty(suite.T(), suite.Service.QueryByUser)
			testutils.AssertAlbums(
				suite.T(),
				[]entity.CheptelAlbum{test.ValidCheptelAlbum},
				albums,
			)
		})
	}
}

func (suite *RepositoryTestSuite) TestCreateFail() {
	testcases := test.ServiceTestCases[schema.CreateRequest, entity.CheptelAlbum]{
		{
			Name: "fail to validate request",
			Req:  schema.CreateRequest{},
		},
		{
			Name: "fail to check if user is member of cheptel",
			Req: schema.CreateRequest{
				UserID:    test.ValidUser.ID,
				CheptelID: test.ValidCheptel.ID,
				AlbumID:   test.ValidCheptel.ID,
				Name:      "new album",
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(test.AnError).Once()
			},
		},
		{
			Name: "fail to create cheptel album",
			Req: schema.CreateRequest{
				UserID:    test.ValidUser.ID,
				CheptelID: test.ValidCheptel.ID,
				AlbumID:   test.ValidCheptel.ID,
				Name:      "new album",
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Create", &entity.CheptelAlbum{
					Album: entity.Album{
						Name:    "new album",
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, test.AnError).Once()
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.Name, func(t *testing.T) {
			tc.ShouldFailAndBeEmpty(suite.T(), suite.Service.Create)
		})
	}
}

func (suite *RepositoryTestSuite) TestCreateSuccess() {
	testcases := test.ServiceTestCases[schema.CreateRequest, entity.CheptelAlbum]{
		{
			Name: "CreateSuccess",
			Req: schema.CreateRequest{
				UserID:    test.ValidUser.ID,
				CheptelID: test.ValidCheptel.ID,
				AlbumID:   test.ValidCheptel.ID,
				Name:      "new album",
			},
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Create", &entity.CheptelAlbum{
					Album: entity.Album{
						Name:    "new album",
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, nil).Once()
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.Name, func(t *testing.T) {
			tc.ShouldSucceedAndNotEmpty(suite.T(), suite.Service.Create)
		})
	}
}

func (suite *RepositoryTestSuite) TestUpdateFail() {
	testcases := test.ServiceTestCases[schema.UpdateRequest, entity.CheptelAlbum]{
		{
			Name: "fail to validate request",
			Req:  schema.UpdateRequest{},
		},
		{
			Name: "fail to check if user is member of cheptel",
			Req: schema.UpdateRequest{
				UserID:    test.ValidUser.ID,
				CheptelID: test.ValidCheptel.ID,
				AlbumID:   test.ValidCheptelAlbum.ID,
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(test.AnError).Once()
			},
		},
		{
			Name: "fail to update cheptel album",
			Req: schema.UpdateRequest{
				UserID:    test.ValidUser.ID,
				CheptelID: test.ValidCheptel.ID,
				AlbumID:   test.ValidCheptelAlbum.ID,
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Get", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, test.AnError).Once()
			},
		},
		{
			Name: "fail to check if user is member of new cheptel",
			Req: schema.UpdateRequest{
				UserID:       test.ValidUser.ID,
				CheptelID:    test.ValidCheptel.ID,
				AlbumID:      test.ValidCheptelAlbum.ID,
				NewCheptelID: test.ValidCheptel.ID,
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Get", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, nil).Once()
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(test.AnError).Once()
			},
		},
		{
			Name: "fail to update cheptel album",
			Req: schema.UpdateRequest{
				UserID:       test.ValidUser.ID,
				CheptelID:    test.ValidCheptel.ID,
				AlbumID:      test.ValidCheptelAlbum.ID,
				NewCheptelID: test.ValidCheptel.ID,
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Get", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, nil).Once()
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Update", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, test.AnError).Once()
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.Name, func(t *testing.T) {
			tc.ShouldFailAndBeEmpty(suite.T(), suite.Service.Update)
		})
	}
}

func (suite *RepositoryTestSuite) TestUpdateSuccess() {
	testcases := test.ServiceTestCases[schema.UpdateRequest, entity.CheptelAlbum]{
		{
			Name: "update succeed",
			Req: schema.UpdateRequest{
				UserID:    test.ValidUser.ID,
				CheptelID: test.ValidCheptel.ID,
				AlbumID:   test.ValidCheptelAlbum.ID,
			},
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Get", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, nil).Once()
				suite.Repository.On("Update", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
					},
				}).Return(test.ValidCheptelAlbum, nil).Once()
			},
		},
		{
			Name: "Update succeed with new cheptel",
			Req: schema.UpdateRequest{
				UserID:       test.ValidUser.ID,
				CheptelID:    test.ValidCheptel.ID,
				AlbumID:      test.ValidCheptelAlbum.ID,
				NewCheptelID: test.ValidCheptel.ID,
			},
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Get", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, nil).Once()
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("Update", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.ValidCheptelAlbum, nil).Once()
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.Name, func(t *testing.T) {
			tc.ShouldSucceedAndNotEmpty(suite.T(), suite.Service.Update)
		})
	}
}

func (suite *RepositoryTestSuite) TestDeleteFail() {
	testcases := test.ServiceTestCases[schema.Request, error]{
		{
			Name: "fail to validate request",
			Req:  schema.Request{},
		},
		{
			Name: "fail to check if user is member of cheptel",
			Req: schema.Request{
				UserID:    test.ValidUser.ID,
				CheptelID: test.ValidCheptel.ID,
				AlbumID:   test.ValidCheptelAlbum.ID,
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(test.AnError).Once()
			},
		},
		{
			Name: "fail to delete cheptel album",
			Req: schema.Request{
				UserID:    test.ValidUser.ID,
				CheptelID: test.ValidCheptel.ID,
				AlbumID:   test.ValidCheptelAlbum.ID,
			},
			WantedError: test.AnError,
			RegisterMocks: func() {
				suite.CheptelManager.On(
					"OnlyMember",
					test.ValidCheptel.ID,
					test.ValidUser.ID).Return(nil).Once()
				suite.Repository.On("SoftDelete", &entity.CheptelAlbum{
					Album: entity.Album{
						Model: gorm.Model{
							ID: test.ValidCheptelAlbum.ID,
						},
						OwnerID: test.ValidCheptel.ID,
					},
				}).Return(test.AnError).Once()
			},
		},
	}

	for _, tc := range testcases {
		suite.T().Run(tc.Name, func(t *testing.T) {
			tc.ShouldFail(suite.T(), suite.Service.Delete)
		})
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
