package auth

import (
	"context"
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/path"
	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/pkg/log"
	refreshtokenschema "github.com/gaetanDubuc/beeckend/internal/refresh-token/schema"
	"github.com/gin-gonic/gin"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	Middleware interface {
		AuthHandler(c *gin.Context) error
		OnlyUnauthenticated(c *gin.Context) error
		OnlyExpiredSession(c *gin.Context) error
		CurrentAuthenticatedUser(ctx context.Context) entity.User
	}

	// Service encapsulates the authentication logic.
	Service interface {
		// authenticate authenticates a user using username and password.
		Login(ctx context.Context, req schema.LoginRequest) (schema.Session, error)
		RefreshSession(ctx context.Context, req refreshtokenschema.RefreshRequest) (schema.Session, error)
		Logout(ctx context.Context, req schema.LogoutRequest) error
	}
)

// RegisterHandlers registers handlers for different HTTP requests.
// A public key is required to verify the JWT token.
func RegisterHandlers(
	rg *routing.RouteGroup,
	service Service,
	authMid Middleware,
	publicKeyPEM string,
	logger log.Logger,
) {
	res := resource{
		service,
		authMid,
		publicKeyPEM,
		logger.With(context.Background(), "api", "auth"),
	}

	rgAuth := rg.Group("")

	rgAuth.Post(path.RefreshSessionPath, authMid.OnlyExpiredSession, res.RefreshSession)

	rgAuth.Post(path.LoginPath, authMid.OnlyUnauthenticated, res.login)
	rgAuth.Get(path.PublicKeyPath, res.publicKey)

	rgAuth.Delete(path.LogoutPath, authMid.AuthHandler, res.logout)
}

type resource struct {
	service      Service
	authMid      Middleware
	publicKeyPEM string
	logger       log.Logger
}

// login handles user basic authentication request.
// Summary	Login and retrieve an authentication JWT and a refresh JWT
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		creds	body		schema.LoginRequest	true	"Credentials"
//	@Success	200		{object}	schema.Session
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/login [post]
//	@Security	BasicAuth
//
//nolint:gofmt
func (r resource) login(c *gin.Context) {
	logger := r.logger.With(c.Request.Context(), "method", "login")

	var req schema.LoginRequest
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		return
	}

	ctx := c.Request.Context()
	session, err := r.service.Login(
		ctx,
		req,
	)

	if err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	c.SecureJSON(http.StatusOK, session)
}

// RefreshJWT handles refresh token request.
// Summary	login the user and refresh its JWT
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	schema.Session
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Param		creds	body		token.ConfirmRequestBase	true	"Credentials"
//	@Router		/refresh-jwt [post]
//	@Security	JWT Token
//
//nolint:gofmt
func (r resource) RefreshSession(c *gin.Context) {
	logger := r.logger.With(c.Request.Context(), "method", "refresh-session")

	var req refreshtokenschema.RefreshRequest
	err := c.Bind(&req)
	if err != nil {
		logger.Error(err)
		return
	}

	ctx := c.Request.Context()

	user := r.authMid.CurrentAuthenticatedUser(ctx)
	req.UserID = user.ID

	session, err := r.service.RefreshSession(ctx, req)
	if err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.SecureJSON(http.StatusOK, session)
}

// logout handles user logout request.
// Summary	Logout and delete the refresh JWTs associated to a user
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	nil
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/logout [delete]
//	@Security	JWT Token
//
//nolint:gofmt
func (r resource) logout(c *routing.Context) {
	logger := r.logger.With(c.Request.Context(), "method", "logout")

	ctx := c.Request.Context()
	user := r.authMid.CurrentAuthenticatedUser(ctx)

	err := r.service.Logout(ctx, schema.LogoutRequest{
		UserID: user.ID,
	})

	if err != nil {
		logger.Error(err)
		c.AbortWithError(http.StatusForbidden, err)
		return
	}

	c.SecureJSON(http.StatusOK)
}

// publicKey handles public key request.
// Summary	retrieve the public key used to verify the JWTs.
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	schema.PublicKeyResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/public-key [get]
//
//nolint:gofmt
func (r resource) publicKey(c *gin.Context) {
	c.SecureJSON(http.StatusOK, schema.PublicKeyResponse{r.publicKeyPEM})
}
