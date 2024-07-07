// Middlewares for authentication. These Middlewares are used to restrict access to
// certain routes or to inject data into the context. They are defined outside of service
// to make a clean separation between the service and the transport layers.
package middleware

import (
	"context"
	e "errors"
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/errors"
	myjwt "github.com/gaetanDubuc/beeckend/pkg/jwt"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUserNotConfirmed     = e.New("user not confirmed")
	ErrUserAlreadyConfirmed = e.New("user already confirmed")
	ErrUserNotFound         = e.New("user not authenticated")
)

type UserRepository interface {
	Get(ctx context.Context, id uint) (entity.User, error)
}

func NewMiddleware[V any](signingMethod string, keyfunc jwt.Keyfunc, userRepository UserRepository, logger log.Logger) Middleware {
	return Middleware{
		signingMethod:  signingMethod,
		keyfunc:        keyfunc,
		userRepository: userRepository,
		logger:         logger,
	}
}

type Middleware struct {
	signingMethod  string
	keyfunc        jwt.Keyfunc
	userRepository UserRepository
	logger         log.Logger
}

// OnlyUnauthenticated reject all authenticated requests.
func (m Middleware) OnlyUnauthenticated(c *gin.Context) error {
	if c.Request.Header.Get("Authorization") != "" {
		// Not explicit to avoid instantly give malicious user that they must
		// not be authenticated to continue with this request.
		// This is by no mean a security feature !
		return errors.ErrForbiddenResp
	}
	return nil
}

func (m Middleware) AuthHandler(c *gin.Context) error {
	return myjwt.JWT(
		&schema.Claims{},
		myjwt.Params[*gin.Context]{
			TokenHandler:  m.handleToken,
			SigningMethod: m.signingMethod,
			Keyfunc:       m.keyfunc,
		})(c)
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
func (m Middleware) handleToken(c *gin.Context, token *jwt.Token) error {
	claim := token.Claims.(schema.Claims)

	ctx := c.Request.Context()
	user, err := m.userRepository.Get(ctx, claim.ID)

	// If the user is not found, we return an error.
	// else if the user has been updated since the token was issued (e.g. confirmation) we need to return an error
	if err != nil {
		return err
	}

	ctx = m.WithAuthenticatedUser(
		ctx,
		user,
	)

	c.Request = c.Request.WithContext(ctx)

	return nil
}

func (m Middleware) OnlyExpiredSession(c *gin.Context) {
	err := myjwt.JWT(
		&schema.ExpiredClaims{},
		myjwt.Params[*gin.Context]{
			TokenHandler:  m.handleToken,
			SigningMethod: m.signingMethod,
		})(c)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
	}
}

type (
	contextKey int
)

const (
	userKey contextKey = iota
)

// WithUser returns a context that contains the user identity from the given JWT.
func (m Middleware) WithAuthenticatedUser(ctx context.Context, user entity.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func (m Middleware) CurrentAuthenticatedUser(ctx context.Context) entity.User {
	user := ctx.Value(userKey)
	if user != nil {
		if user, ok := user.(entity.User); ok {
			return user
		}
	}
	panic(ErrUserNotFound)
}
