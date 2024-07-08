package router

import (
	"io"

	docs "github.com/gaetanDubuc/beeckend/api"
	dbx "github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/pkg/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

func New(out io.Writer, db *dbx.DB) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(out), gin.Recovery())

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := r.Group(docs.SwaggerInfo.BasePath)

	v1.Use(
		middleware.JSONContentType,
		middleware.SecurityHeaders,
		db.TransactionHandler(),
	)
	return r, v1
}
