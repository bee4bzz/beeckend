package api

import (
	"context"
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	authschema "github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/internal/refresh-session/path"
	refreshtokenschema "github.com/gaetanDubuc/beeckend/internal/refresh-session/schema"
	"github.com/gin-gonic/gin"
)

type (
	Middleware interface {
		AuthHandler(c *gin.Context)
		CurrentAuthenticatedUser(ctx context.Context) entity.User
	}

	// Service encapsulates the authentication logic.
	Service interface {
		Refresh(ctx context.Context, req refreshtokenschema.RefreshRequest) (string, error)
	}

	AuthService interface {
		GenerateJWT(ctx context.Context, userID uint) (string, error)
	}
)

// RegisterHandlers registers handlers for different HTTP requests.
// A public key is required to verify the JWT token.
func RegisterHandlers(
	rg *gin.RouterGroup,
	service Service,
	authService AuthService,
	authMid Middleware,
	publicKeyPEM string,
	logger log.Logger,
) {
	res := resource{
		service,
		authService,
		authMid,
		publicKeyPEM,
		logger.With(context.Background(), "api", "auth"),
	}

	rgRefreshSession := rg.Group("")

	rgRefreshSession.POST(
		path.RefreshSessionPath,
		authMid.AuthHandler,
		authMid.AuthHandler,
		res.RefreshSession)
}

type resource struct {
	service      Service
	authService  AuthService
	authMid      Middleware
	publicKeyPEM string
	logger       log.Logger
}

// RefreshJWT handles refresh token request.
// Summary	login the user and refresh its JWT
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	schema.Session
//	@Param		creds	body		refreshtokenschema.RefreshRequest	true	"Credentials"
//	@Router		/refresh-jwt [post]
//	@Security	JWT Token
//
//nolint:gofmt
func (r *resource) RefreshSession(c *gin.Context) {
	logger := r.logger.With(c.Request.Context(), "method", "refresh-session")

	var req refreshtokenschema.RefreshRequest
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := c.Request.Context()

	user := r.authMid.CurrentAuthenticatedUser(ctx)
	req.UserID = user.ID

	RefreshJWT, err := r.service.Refresh(ctx, req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	JWT, err := r.authService.GenerateJWT(ctx, user.ID)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.SecureJSON(http.StatusOK, authschema.Session{
		JWT:        JWT,
		RefreshJWT: RefreshJWT,
	})
}

// publicKey handles public key request.
// Summary	retrieve the public key used to verify the JWTs.
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	schema.PublicKeyResponse
//	@Router		/refresh-jwt/public-key [get]
//
//nolint:gofmt
func (r resource) publicKey(c *gin.Context) {
	c.SecureJSON(http.StatusOK, schema.PublicKeyResponse{r.publicKeyPEM})
}
