// Middlewares for authentication. These Middlewares are used to restrict access to
// certain routes or to inject data into the context. They are defined outside of service
// to make a clean separation between the service and the transport layers.
package auth

import (
	"context"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/errors"
	"github.com/gin-gonic/gin"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/golang-jwt/jwt"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/pkg/log"
)

type JWTVerifier[V any] interface {
	VerificationKey() V
	SigninMethod() string
}

func NewMiddleware[V any](jwtVerifier JWTVerifier[V], userRepository UserRepository, GetMQTTHeaderPwd func() string, logger log.Logger) Middleware {
	return middleware[V]{
		jwtVerifier:    jwtVerifier,
		userRepository: userRepository,
		logger:         logger,
	}
}

type middleware[V any] struct {
	jwtVerifier      JWTVerifier[V]
	userRepository   UserRepository
	GetMQTTHeaderPwd func() string
	logger           log.Logger
}

// OnlyUnauthenticated reject all authenticated requests.
func (m middleware[V]) OnlyUnauthenticated(c *routing.Context) error {
	if c.Request.Header.Get("Authorization") != "" {
		// Not explicit to avoid instantly give malicious user that they must
		// not be authenticated to continue with this request.
		// This is by no mean a security feature !
		return errors.ErrForbiddenResp
	}
	return nil
}

func (m middleware[V]) AuthHandler(c *gin.Context) error {
	return fogojwt.JWT(
		m.jwtVerifier.VerificationKey(),
		&schema.JWTClaims{},
		fogojwt.Options[V]{
			TokenHandler:  m.handleToken,
			SigningMethod: m.jwtVerifier.SigninMethod(),
		})(c)
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
func (m middleware[V]) handleToken(c *routing.Context, token *jwt.Token) error {
	claim := token.Claims.(Claim)

	ctx := c.Request.Context()
	user, err := m.userRepository.Get(ctx, claim.GetUUID())

	// If the user is not found, we return an error.
	// else if the user has been updated since the token was issued (e.g. confirmation) we need to return an error
	if err != nil {
		return errors.ErrInternalServerResp.Join(err)
	} else if user.Administrator != claim.GetAdministrator().ElseZero() && user.Confirmed != claim.GetConfirmed().ElseZero() {
		return auth.ErrInvalidSessionResp
	}

	ctx = m.WithAuthenticatedUser(
		ctx,
		user,
	)

	c.Request = c.Request.WithContext(ctx)

	return nil
}

func (m middleware[V]) OnlyExpiredSession(c *routing.Context) error {
	return fogojwt.JWT(
		m.jwtVerifier.VerificationKey(),
		&schema.ExpiredJWTClaims{},
		fogojwt.Options[V]{
			TokenHandler:  m.handleToken,
			SigningMethod: m.jwtVerifier.SigninMethod(),
		})(c)
}

// OnlyRootAdmin uses context identity to check if identity is an admin.
func (m middleware[V]) OnlyRootAdmin(c *routing.Context) error {
	user := m.CurrentAuthenticatedUser(c.Request.Context())

	if !user.IsAdministrator() {
		return errors.ErrForbiddenResp
	}

	return nil
}

func (m middleware[V]) OnlyConfirmedUser(c *routing.Context) error {
	user := m.CurrentAuthenticatedUser(c.Request.Context())

	if !user.IsConfirmed() {
		return auth.ErrUserNotConfirmedResp
	}

	return nil
}

func (m middleware[V]) OnlyUnconfirmedUser(c *routing.Context) error {
	user := m.CurrentAuthenticatedUser(c.Request.Context())

	if user.IsConfirmed() {
		return auth.ErrUserAlreadyConfirmedResp
	}

	return nil
}

type (
	contextKey int
)

const (
	userKey contextKey = iota
)

// WithUser returns a context that contains the user identity from the given JWT.
func (m middleware[V]) WithAuthenticatedUser(ctx context.Context, user entity.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func (m middleware[V]) CurrentAuthenticatedUser(ctx context.Context) entity.User {
	user := ctx.Value(userKey)
	if user != nil {
		if user, ok := user.(entity.User); ok {
			return user
		}
	}
	panic(auth.ErrUserNotFoundResp)
}
