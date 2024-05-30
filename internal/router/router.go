package router

import (
	"database/sql"

	docs "github.com/gaetanDubuc/beeckend/api"
	dbx "github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/errors"
	"github.com/gaetanDubuc/beeckend/pkg/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @BasePath /api/v1

func New(db *dbx.DB) (*gin.Engine, *gin.RouterGroup) {
	r := gin.Default()
	engine := gin.New()
	engine.Use(gin.Logger(), gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		errorResponse := buildErrorResponse(err)
		c.JSON(errorResponse.StatusCode(), errorResponse)
	}))

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := r.Group(docs.SwaggerInfo.BasePath)

	v1.Use(
		middleware.JSONContentType,
		db.TransactionHandler(),
	)
	return r, v1
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err any) errors.ErrorResponse {
	switch err := err.(type) {
	case errors.ErrorResponse:
		return err
	}

	switch err {
	case sql.ErrNoRows, gorm.ErrRecordNotFound, gorm.ErrForeignKeyViolated:
		return errors.ErrForbiddenResp
	default:
		return errors.ErrInternalServerResp
	}
}
