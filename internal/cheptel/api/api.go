package api

import (
	"context"
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/path"
	"github.com/gaetanDubuc/beeckend/internal/cheptel/schema"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/pkg/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Service interface {
	Subscribe(ctx context.Context, user *entity.User, cheptels chan<- *[]entity.Cheptel) error
	Create(ctx context.Context, req schema.CreateRequest) (entity.Cheptel, error)
}

type Upgrader[T Conn] interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (T, error)
}

type Conn interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type Middleware interface {
	AuthHandler(c *gin.Context)
	CurrentAuthenticatedUser(ctx context.Context) entity.User
}

func RegisterHandlers[T Conn](
	router *gin.RouterGroup,
	service Service,
	authMiddleware Middleware,
	upgrader Upgrader[T],
	logger log.Logger,
) {
	resource := &Resource[T]{upgrader, service, authMiddleware, logger}
	router.GET(path.Query, authMiddleware.AuthHandler, resource.query)
	router.POST(path.CheptelsGroup, authMiddleware.AuthHandler, resource.create)
	// router.PUT(path.CheptelsGroup, update)
	// router.DELETE(path.MakeDeletePath(path.CheptelParam), delete)
}

type Resource[T Conn] struct {
	upgrader       Upgrader[T]
	service        Service
	authMiddleware Middleware
	logger         log.Logger
}

// @Summary	Query the user's cheptels
// @Schemes
// @Tags		Cheptel
// @Accept		json
// @Produce	json
// @Success	200	{string}	query
// @Router		/cheptels [get]
// @Security	JWT Token
func (r *Resource[T]) query(c *gin.Context) {
	ctx := c.Request.Context()
	user := r.authMiddleware.CurrentAuthenticatedUser(ctx)

	conn, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		r.logger.Error("upgrade:", err)
		panic(err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			r.logger.Error("close:", err)
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	chCheptels := make(chan *[]entity.Cheptel)
	err = r.service.Subscribe(ctx, &user, chCheptels)
	if err != nil {
		r.logger.Error("subscribe:", err)
		panic(err)
	}
	defer cancel()

	for cheptel := range chCheptels {
		err = conn.WriteMessage(websocket.BinaryMessage, json.MustMarshal(cheptel))
		if err != nil {
			r.logger.Error("server write:", err)
			panic(err)
		}
	}
}

// @Summary	Create a cheptel
// @Schemes
// @Tags		Cheptel
// @Accept		json
// @Produce	json
// @Success	200	{string}	create
// @Param		cheptel	body	schema.CreateRequest	true	"Cheptel to create"
// @Router		/cheptels [post]
// @Security	JWT Token
func (r *Resource[T]) create(c *gin.Context) {
	logger := r.logger.With(c.Request.Context(), "method", "create")

	ctx := c.Request.Context()
	user := r.authMiddleware.CurrentAuthenticatedUser(ctx)

	var req schema.CreateRequest
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	req.UserID = user.ID

	cheptel, err := r.service.Create(ctx, req)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.SecureJSON(http.StatusOK, cheptel)
}
