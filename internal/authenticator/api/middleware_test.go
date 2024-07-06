package auth

import (
	"context"
	"testing"

	"4d63.com/optional"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	auth "gitlab.com/fogo-dev/infrastructure/web-api/internal/auth"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/auth/schema"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/errors"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/test"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/testutils"
	authtestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/testutils/auth"
	usertestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/testutils/user"
	fogojwt "gitlab.com/fogo-dev/infrastructure/web-api/pkg/jwt"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
)

func InitMid() (mid Middleware, authenticationRepository *testutils.JWTManager[[]byte], authorizationRepository *authtestutils.Repository, userRepo *usertestutils.Repository, ctx *routing.Context) {
	ctx, _ = testutils.RoutingContext()
	logger, _, _ := log.NewForTest()

	authenticationRepository = &testutils.JWTManager[[]byte]{
		GenerateJWTReturn:     test.ValidToken,
		SigninMethodReturn:    fogojwt.HS256,
		VerificationKeyReturn: []byte(test.ValidHSSigningKey),
	}
	authorizationRepository = &authtestutils.Repository{}

	userRepo = &usertestutils.Repository{
		Repository: testutils.Repository[entity.User]{
			GetReturn: test.UserAdminConfirmed},
	}

	mid = NewMiddleware(
		authenticationRepository,
		userRepo,
		func() string {
			return "test"
		},
		logger,
	)

	ctx.Request = ctx.Request.WithContext(mid.WithAuthenticatedUser(ctx.Request.Context(), test.UserAdminConfirmed))
	return
}

func TestUserMiddleware_OnlyConfirmedUser(t *testing.T) {
	t.Run("user is confirmed, should not raise an error.", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		err := mid.OnlyConfirmedUser(ctx)

		assert.NoError(t, err)
	})
	t.Run("user is not confirmed, should raise an error.", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()
		ctx.Request = ctx.Request.WithContext(mid.WithAuthenticatedUser(ctx.Request.Context(), test.UserAdminUnconfirmed))

		err := mid.OnlyConfirmedUser(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, auth.ErrUserNotConfirmed)
	})
}

func TestUserMiddleware_IngressHeaderAuthentication(t *testing.T) {
	t.Run("request has X-FOGO-MQTT-INGRESS header", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()
		ctx.Request.Header.Set("X-FOGO-MQTT-INGRESS", "test")

		err := mid.IngressHeaderAuthentication(ctx)

		assert.NoError(t, err)
	})

	t.Run("request has no X-FOGO-MQTT-INGRESS header", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		err := mid.IngressHeaderAuthentication(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrForbiddenResp)
	})
}

func TestUserMiddleware_OnlyUnconfirmedUser(t *testing.T) {
	t.Run("user is confirmed, should raise an error.", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		err := mid.OnlyUnconfirmedUser(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, auth.ErrUserAlreadyConfirmed)
	})

	t.Run("user is not confirmed, should not raise an error.", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()
		ctx.Request = ctx.Request.WithContext(mid.WithAuthenticatedUser(ctx.Request.Context(), test.UserAdminUnconfirmed))

		err := mid.OnlyUnconfirmedUser(ctx)

		assert.NoError(t, err)
	})
}

func TestUserMiddleware_OnlyRootAdmin(t *testing.T) {
	t.Run("user is an administrator, should not raise any errors.", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		assert.NoError(t, mid.OnlyRootAdmin(ctx))
	})
	t.Run("user is not an administrator, should raise an error.", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()
		ctx.Request = ctx.Request.WithContext(mid.WithAuthenticatedUser(ctx.Request.Context(), test.UserNoAdminConfirmed))

		err := mid.OnlyRootAdmin(ctx)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrForbiddenResp)
	})
}

func TestAuthMiddleware_CurrentAuthenticatedUser(t *testing.T) {
	mid, _, _, _, c := InitMid()

	ctx := context.Background()
	assert.Panics(t, func() { mid.CurrentAuthenticatedUser(ctx) })

	identity := mid.CurrentAuthenticatedUser(c.Request.Context())
	assert.Equal(t, identity, test.UserAdminConfirmed)
}

func Test_handleToken(t *testing.T) {
	t.Run("user is authenticated but the claim is not valid, should raise an error", func(t *testing.T) {
		mid, _, _, _, c := InitMid()

		err := mid.(middleware[[]byte]).handleToken(c, &jwt.Token{
			Claims: &schema.JWTClaims{
				UUID:          entity.GenerateUUID(),
				Administrator: optional.Of(true),
				Confirmed:     optional.Of(true),
				Exp:           0,
			},
		})

		assert.NoError(t, err)
	})

	t.Run("user is authenticated but an error occur while getting its information, should raise an error", func(t *testing.T) {
		mid, _, _, userRepo, c := InitMid()
		userRepo.GetError = testutils.ErrMocked

		err := mid.(middleware[[]byte]).handleToken(c, &jwt.Token{
			Claims: &schema.JWTClaims{},
		})

		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrInternalServer)
		assert.ErrorIs(t, err, testutils.ErrMocked)
	})

	t.Run("user is authenticated but the claim is not the same than the db, should raise an error", func(t *testing.T) {
		mid, _, _, _, c := InitMid()

		err := mid.(middleware[[]byte]).handleToken(c, &jwt.Token{
			Claims: &schema.JWTClaims{},
		})

		assert.Error(t, err)
		assert.ErrorIs(t, err, auth.ErrInvalidSessionResp)
	})
}

func TestMiddleware_VerifyJWT(t *testing.T) {
	t.Run("succeed to verify JWT", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer "+test.ValidJWT)

		err := mid.AuthHandler(ctx)
		assert.NoError(t, err)
	})

	t.Run("Bearer is ill formed, should raise an error", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer ")

		err := mid.AuthHandler(ctx)
		assert.Error(t, err)
	})

	t.Run("the jwt is expired, should raise an error", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer "+test.ExpiredJWT)

		err := mid.AuthHandler(ctx)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Token is expired")
	})

	t.Run("the jwt is malformed, should raise an error", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer "+test.InvalidJWT)

		err := mid.AuthHandler(ctx)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "confirmed: cannot be blank")
	})
}

func TestMiddleware_OnlyExpiredSession(t *testing.T) {
	t.Run("succeed to verify expired JWT", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer "+test.ExpiredJWT)

		err := mid.OnlyExpiredSession(ctx)
		assert.NoError(t, err)
	})

	t.Run("Bearer is ill formed, should raise an error", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer ")

		err := mid.OnlyExpiredSession(ctx)
		assert.Error(t, err)
	})

	t.Run("the jwt is not expired, should raise an error", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer "+test.ValidJWT)

		err := mid.OnlyExpiredSession(ctx)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "Token is not expired")
	})
	t.Run("the jwt is malformed, should raise an error", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer "+test.InvalidJWT)

		err := mid.AuthHandler(ctx)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "confirmed: cannot be blank")
	})
}

func TestMiddleware_OnlyUnauthenticated(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		err := mid.OnlyUnauthenticated(ctx)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mid, _, _, _, ctx := InitMid()

		ctx.Request.Header.Set("Authorization", "Bearer")

		err := mid.OnlyUnauthenticated(ctx)
		assert.ErrorIs(t, errors.ErrForbiddenResp, err)
	})
}
