package main

import (
	"net/http"

	cheptelapi "github.com/gaetanDubuc/beeckend/internal/cheptel/api"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
)

func main() {
	logger := utils.NewLogger()

	config, err := utils.LoadConfig(".")
	if err != nil {
		logger.Fatal("cannot load config:", err)
	}

	db := db.NewGormWithMigrate(
		postgres.Open(config.DBSource),
		"file://./migrations",
		config.DatabaseURL,
		logger)

	router, v1 := utils.NewRouter(db)
	RegisterHandlers(v1, db, logger)
	server := utils.NewServer(config.ServerAddress, router)

	go func() {
		// service connections
		logger.Infof("server is running at %v", config.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	utils.GracefulShutdown(server, logger)
}

func RegisterHandlers(router *gin.RouterGroup, db *db.DB, logger log.Logger) {
	cheptelapi.RegisterHandlers(router)
}
