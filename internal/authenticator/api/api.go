package api

import (
	"context"
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/path"
	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	refreshsessionschema "github.com/gaetanDubuc/beeckend/internal/refresh-session/schema"
	"github.com/gin-gonic/gin"
)

type (
	Middleware interface {
		AuthHandler(c *gin.Context)
		OnlyUnauthenticated(c *gin.Context)
		CurrentAuthenticatedUser(ctx context.Context) entity.User
	}

	// Service encapsulates the authentication logic.
	Service interface {
		// authenticate authenticates a user using username and password.
		Login(ctx context.Context, req schema.LoginRequest) (string, error)
	}

	RefreshSessionService interface {
		Create(ctx context.Context, req refreshsessionschema.CreateRequest) (string, error)
		DeleteFromUser(ctx context.Context, req refreshsessionschema.DeleteFromUserRequest) error
	}

	UserService interface {
		Get(ctx context.Context, user *entity.User) error
	}
)

// RegisterHandlers registers handlers for different HTTP requests.
// A public key is required to verify the JWT token.
func RegisterHandlers(
	rg *gin.RouterGroup,
	service Service,
	refreshSessionService RefreshSessionService,
	userService UserService,
	authMid Middleware,
	publicKeyPEM string,
	logger log.Logger,
) {
	res := resource{
		service,
		refreshSessionService,
		userService,
		authMid,
		publicKeyPEM,
		logger.With(context.Background(), "api", "auth"),
	}

	rgAuth := rg.Group("")

	rgAuth.POST(path.LoginPath, authMid.OnlyUnauthenticated, res.login)
	rgAuth.GET(path.PublicKeyPath, res.publicKey)

	rgAuth.DELETE(path.LogoutPath, authMid.AuthHandler, res.logout)
}

type resource struct {
	service               Service
	refreshSessionService RefreshSessionService
	userService           UserService
	authMid               Middleware
	publicKeyPEM          string
	logger                log.Logger
}

// login handles user basic authentication request.
// Summary	Login and retrieve an authentication JWT and a refresh JWT
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		creds	body		schema.LoginRequest	true	"Credentials"
//	@Success	200		{object}	schema.Session
//	@Router		/login [post]
//	@Security	BasicAuth
//
//nolint:gofmt
func (r resource) login(c *gin.Context) {
	logger := r.logger.With(c.Request.Context(), "method", "login")

	var req schema.LoginRequest
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()
	JWT, err := r.service.Login(
		ctx,
		req,
	)

	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := entity.User{
		Email: req.Username,
	}

	err = r.userService.Get(ctx, &user)

	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	refreshJWT, err := r.refreshSessionService.Create(ctx, refreshsessionschema.CreateRequest{
		UserID: user.ID,
	})

	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.SecureJSON(http.StatusOK, schema.Session{
		JWT:        JWT,
		RefreshJWT: refreshJWT,
	})
}

// logout handles user logout request.
// Summary	Logout and delete the refresh JWTs associated to a user
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	nil
//	@Router		/logout [delete]
//	@Security	JWT Token
//
//nolint:gofmt
func (r resource) logout(c *gin.Context) {
	logger := r.logger.With(c.Request.Context(), "method", "logout")

	ctx := c.Request.Context()
	user := r.authMid.CurrentAuthenticatedUser(ctx)

	err := r.refreshSessionService.DeleteFromUser(ctx, refreshsessionschema.DeleteFromUserRequest{
		UserID: user.ID,
	})

	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}

// publicKey handles public key request.
// Summary	retrieve the public key used to verify the JWTs.
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	schema.PublicKeyResponse
//	@Router		/public-key [get]
//
//nolint:gofmt
func (r resource) publicKey(c *gin.Context) {
	c.SecureJSON(http.StatusOK, schema.PublicKeyResponse{r.publicKeyPEM})
}
