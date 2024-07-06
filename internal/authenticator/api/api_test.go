package auth

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/stretchr/testify/suite"
	refreshtokenschema "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/refresh-token/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/testutils"
	authtestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/testutils"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/errors"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/test/v2"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/utils"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
	"go.uber.org/zap/zaptest/observer"
)

type APITestSuite struct {
	suite.Suite
	ctx            context.Context
	router         *routing.Router
	service        *testutils.Service
	authMiddleware *authtestutils.Middleware
	baseURL        string
	logger         log.Logger
	observer       *observer.ObservedLogs
	buffer         *bytes.Buffer

	loginRootTest          test.APITestCase[schema.Session]
	refreshSessionRootTest test.APITestCase[schema.Session]
	publicKeyRootTest      test.APITestCase[schema.PublicKeyResponse]
	logoutRootTest         test.APITestCase[any]
}

// this function executes before the test suite begins execution.
func (suite *APITestSuite) SetupSuite() {
	suite.ctx = context.Background()
	logger, obs, b := log.NewForTest()
	suite.logger = logger
	suite.observer = obs
	suite.buffer = b

	suite.router = utils.NewRouter(logger)
	suite.service = &testutils.Service{}
	suite.authMiddleware = &testutils.Middleware{}
	RegisterHandlers(
		suite.router.Group(""),
		suite.service,
		suite.authMiddleware,
		string(test.ValidRSAPublicKeyPEM),
		logger)

	suite.loginRootTest = testutils.LoginRootTest.
		WithRouter(suite.router).
		WithLogger(suite.logger)

	suite.refreshSessionRootTest = testutils.RefreshSessionRootTest.
		WithRouter(suite.router).
		WithLogger(suite.logger)

	suite.logoutRootTest = testutils.LogoutRootTest.
		WithRouter(suite.router).
		WithLogger(suite.logger)

	suite.publicKeyRootTest = testutils.PublicKeyRootTest.
		WithRouter(suite.router).
		WithLogger(suite.logger)
}

func (suite *APITestSuite) SetupTest() {
	suite.buffer.Reset()
	suite.observer.TakeAll()
}

func (suite *APITestSuite) TearDownTest() {
	suite.service.AssertExpectations(suite.T())
	suite.authMiddleware.AssertExpectations(suite.T())
	suite.T().Log(suite.buffer)
}

func (suite *APITestSuite) Test_A_User_Can_Login() {
	req := schema.LoginRequest{
		Username: test.UserAdminConfirmed.Email,
		Password: test.ValidPassword,
	}
	suite.authMiddleware.On("OnlyUnauthenticated").
		Return(nil).Once()
	suite.service.On("Login", req).
		Return(schema.Session{
			Token:        test.ValidJWT,
			RefreshToken: test.ValidJWT,
		}, nil).Once()

	suite.loginRootTest.WithRequest(
		req,
	).CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_Return_Error_When_A_User_Login_Being_Already_Authenticated() {
	suite.authMiddleware.On("OnlyUnauthenticated").
		Return(errors.ErrUnauthorizedResp).Once()

	suite.loginRootTest.
		WithWantStatus(http.StatusUnauthorized).
		WithWantResponse(errors.ErrUnauthorizedResp.Message).
		WithRequest(
			schema.LoginRequest{
				Username: test.UserAdminConfirmed.Email,
				Password: test.ValidPassword,
			},
		).CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_Return_Error_When_A_User_Login_With_Wrong_Credentials() {
	req := schema.LoginRequest{
		Username: test.UserAdminConfirmed.Email,
		Password: test.InvalidPassword,
	}
	suite.authMiddleware.On("OnlyUnauthenticated").
		Return(nil).Once()
	suite.service.On("Login", req).
		Return(schema.Session{}, errors.ErrForbiddenResp).Once()

	suite.loginRootTest.
		WithWantStatus(http.StatusForbidden).
		WithWantResponse(errors.ErrForbiddenResp.Message).
		WithRequest(req).
		CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_A_User_Can_Logout() {
	suite.authMiddleware.On("AuthHandler").
		Return(nil).Once()

	suite.authMiddleware.On("CurrentAuthenticatedUser").
		Return(test.UserAdminUnconfirmed).Once()

	suite.service.On("Logout", schema.LogoutRequest{
		UserUUID: test.UserAdminUnconfirmed.UUID,
	}).
		Return(nil).Once()

	suite.logoutRootTest.CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_Return_Error_When_A_User_Logout_Being_Not_Authenticated() {
	suite.authMiddleware.On("AuthHandler").
		Return(errors.ErrUnauthorizedResp).Once()

	suite.logoutRootTest.
		WithWantStatus(http.StatusUnauthorized).
		WithWantResponse(errors.ErrUnauthorizedResp.Message).
		CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_Return_Error_When_The_Service_Can_Not_Logout() {
	suite.authMiddleware.On("AuthHandler").
		Return(nil).Once()

	suite.authMiddleware.On("CurrentAuthenticatedUser").
		Return(test.UserAdminUnconfirmed).Once()

	suite.service.On("Logout", schema.LogoutRequest{
		UserUUID: test.UserAdminUnconfirmed.UUID,
	}).
		Return(errors.ErrInternalResp).Once()

	suite.logoutRootTest.
		WithWantStatus(http.StatusInternalServerError).
		WithWantResponse(errors.ErrInternalResp.Message).
		CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_A_User_Can_Get_Public_Key() {
	suite.publicKeyRootTest.CheckEndpoint(suite.T())
}

func (suite *APITestSuite) Test_A_User_Can_Refresh_Session() {
	req := refreshtokenschema.RefreshRequest{
		UserUUID:  test.UserAdminUnconfirmed.UUID,
		TokenUUID: test.TokenAdminUnconfirmed.UUID,
		Token:     test.ValidJWT,
	}
	suite.authMiddleware.On("OnlyExpiredSession").
		Return(nil).Once()
	suite.authMiddleware.On("CurrentAuthenticatedUser").
		Return(test.UserAdminUnconfirmed).Once()
	suite.service.On("RefreshSession", req).
		Return(schema.Session{
			Token:        test.ValidJWT,
			RefreshToken: test.ValidJWT,
		}, nil).Once()

	suite.refreshSessionRootTest.
		WithRequest(
			req,
		).CheckEndpoint(suite.T())
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

/*
func Init(prepare func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware)) *routing.Router {
	logger, _, _ := log.NewForTest()

	s := &authtestutils.Service{
		LoginReturn:       test.ValidJWT,
		GenerateJWTReturn: test.ValidJWT,
	}
	refreshTokenService := &refreshtestutils.Service{
		RefreshTokenReturn: test.ValidJWT,
		GenerateJWTReturn:  test.ValidJWT,
	}
	authMid := &authtestutils.Middleware{}
	if prepare != nil {
		prepare(s, refreshTokenService, authMid)
	}

	router := testutils.Router(logger)
	RegisterHandlers(router.Group(""), s, refreshTokenService, authMid, string(test.ValidRSAPublicKeyPEM), logger)

	return router
}

func TestAPI_Login(t *testing.T) {
	ValidResponse := fmt.Sprintf(`{"%s":"%s".*"%s":"%s".*`, TokenJSONKey, test.ValidJWT, RefreshTokenJSONKey, test.ValidJWT)
	// root test common to all subtests
	rootTest := authtestutils.LoginRootTest("").CopyWith(
		test.OpAPITestCase{
			Router:       optional.Of(Init(nil)),
			Body:         optional.Of(FullLoginRequest),
			WantResponse: optional.Of(ValidResponse),
		},
	)

	t.Run("success", func(t *testing.T) {
		rootTest.CheckEndpoint(t)
		authtestutils.TestWithBasicAuth(rootTest, test.ValidEmail, test.ValidPassword).CheckEndpoint(t)
	})
	t.Run("login query is validated and fail", func(t *testing.T) {
		test := rootTest.CopyWith(test.OpAPITestCase{
			Body:         optional.Of("{}"),
			WantStatus:   optional.Of(http.StatusBadRequest),
			WantResponse: optional.Of(""),
		})
		test.CheckEndpoint(t)
	})
	t.Run("login fail because of internal error", func(t *testing.T) {
		prepares := []func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware){
			func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware) {
				s.LoginError = testutils.ErrMockedResp
			},
			func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware) {
				refreshTokenService.RefreshTokenError = testutils.ErrMockedResp
			},
		}

		for _, prepare := range prepares {
			t.Run(FunctionName(prepare), func(t *testing.T) {
				AssertFailedResponseRouter(
					t,
					rootTest,
					Init(prepare),
				)
			})
		}
	})
}

func TestAPI_Logout(t *testing.T) {
	// root test common to all subtests
	rootTest := authtestutils.LogoutRootTest("", nil).CopyWith(
		test.OpAPITestCase{
			Router: optional.Of(Init(nil)),
		},
	)

	t.Run("success", func(t *testing.T) {
		rootTest.CheckEndpoint(t)
	})
	t.Run("delete fail because of internal error", func(t *testing.T) {
		prepares := []func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware){
			func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware) {
				refreshTokenService.DeleteFromUserError = testutils.ErrMockedResp
			},
		}

		for _, prepare := range prepares {
			t.Run(FunctionName(prepare), func(t *testing.T) {
				AssertFailedResponseRouter(
					t,
					rootTest,
					Init(prepare),
				)
			})
		}
	})
}

func TestAPI_PublicKey(t *testing.T) {
	ValidResponse := `"key":.*`
	// root test common to all subtests
	rootTest := authtestutils.PublicKeyRootTest("").CopyWith(
		test.OpAPITestCase{
			Router:       optional.Of(Init(nil)),
			WantResponse: optional.Of(ValidResponse),
		},
	)

	t.Run("success", func(t *testing.T) {
		rootTest.CheckEndpoint(t)
	})
}

func TestAPI_RefreshJWT(t *testing.T) {
	ValidResponse := fmt.Sprintf(`.*"%s":"%s".*"%s":"%s".*`, TokenJSONKey, test.ValidJWT, RefreshTokenJSONKey, test.ValidJWT)
	rootTest := refreshtestutils.RefreshJWTRootTest(
		"",
		nil,
		test.ValidToken,
		test.TokenAdminUnconfirmed.OwnerUUID.String(),
	)
	t.Run("success", func(t *testing.T) {
		rootTest.CopyWith(test.OpAPITestCase{
			Router:       optional.Of(Init(nil)),
			WantResponse: optional.Of(ValidResponse),
		}).CheckEndpoint(t)
	})

	t.Run("fail due to internal error or middleware", func(t *testing.T) {
		prepares := []func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware){
			func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware) {
				authMid.OnlyExpiredSessionError = testutils.ErrMockedResp
			},
			func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware) {
				refreshTokenService.ConfirmTokenError = testutils.ErrMockedResp
			},
			func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware) {
				refreshTokenService.GetError = testutils.ErrMockedResp
			},
			func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware) {
				refreshTokenService.RefreshTokenError = testutils.ErrMockedResp
			},
			func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware) {
				s.GenerateJWTError = testutils.ErrMockedResp
			},
		}

		for _, prepare := range prepares {
			t.Run(FunctionName(prepare), func(t *testing.T) {
				tokentestutils.AssertFailedResponseRouter(
					t,
					rootTest,
					Init(prepare),
				)
			})
		}
	})
}

func AssertFailedResponseRouter(t *testing.T, tc test.APITestCase, router *routing.Router) {
	test.AssertFailedResponseRouter[entity.User](t, tc, router, http.StatusSeeOther)
}

func FunctionName(f func(s *authtestutils.Service, refreshTokenService *refreshtestutils.Service, authMid *authtestutils.Middleware)) string {
	return test.FunctionName(f)
}
*/
