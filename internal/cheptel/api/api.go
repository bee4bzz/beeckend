package api

import (
	"context"
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/path"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/pkg/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Service interface {
	Subscribe(ctx context.Context, user *entity.User, cheptels chan<- *[]entity.Cheptel) error
}

type Upgrader[T Conn] interface {
	Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (T, error)
}

type Conn interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
}

func RegisterHandlers[T Conn](router *gin.RouterGroup, service Service, upgrader Upgrader[T], logger log.Logger) {
	resource := &Resource[T]{upgrader, service, logger}
	router.GET(path.Query, resource.query)
	// router.POST(path.CheptelsGroup, create)
	// router.PUT(path.CheptelsGroup, update)
	// router.DELETE(path.MakeDeletePath(path.CheptelParam), delete)
}

type Resource[T Conn] struct {
	upgrader Upgrader[T]
	service  Service
	logger   log.Logger
}

// @Summary	Query the user's cheptels
// @Schemes
// @Tags		Cheptel
// @Accept		json
// @Produce	json
// @Success	200	{string}	query
// @Router		/cheptels [get]
func (r *Resource[T]) query(c *gin.Context) {
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

	ctx, cancel := context.WithCancel(c.Request.Context())
	chCheptels := make(chan *[]entity.Cheptel)
	err = r.service.Subscribe(ctx, &entity.User{}, chCheptels)
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
