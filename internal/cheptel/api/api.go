package api

import (
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/path"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func RegisterHandlers(router *gin.RouterGroup, upgrader *websocket.Upgrader, logger log.Logger) {
	resource := &Resource{upgrader, logger}
	router.GET(path.Query, resource.query)
	// router.POST(path.CheptelsGroup, create)
	// router.PUT(path.CheptelsGroup, update)
	// router.DELETE(path.MakeDeletePath(path.CheptelParam), delete)
}

type Resource struct {
	upgrader *websocket.Upgrader
	logger   log.Logger
}

// @Summary	Query the user's cheptels
// @Schemes
// @Tags		Cheptel
// @Accept		json
// @Produce	json
// @Success	200	{string}	query
// @Router		/cheptels [get]
func (r *Resource) query(c *gin.Context) {
	conn, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		r.logger.Error("upgrade:", err)
		panic(err)
	}
	defer conn.Close()
	for {
		err = conn.WriteMessage(0, []byte("query"))
		if err != nil {
			r.logger.Error("write:", err)
			break
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "query"})
}
