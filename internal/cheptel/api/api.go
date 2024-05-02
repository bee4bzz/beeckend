package api

import (
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/path"
	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.RouterGroup) {
	router.GET(path.Query, query)
	// router.POST(path.CheptelsGroup, create)
	// router.PUT(path.CheptelsGroup, update)
	// router.DELETE(path.MakeDeletePath(path.CheptelParam), delete)
}

// @Summary	Query the user's cheptels
// @Schemes
// @Tags		Cheptel
// @Accept		json
// @Produce	json
// @Success	200	{string}	query
// @Router		/cheptels [get]
func query(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "query"})
}
