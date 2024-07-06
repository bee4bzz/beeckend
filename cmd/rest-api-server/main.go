package main

import (
	"context"
	"net/http"

	cheptelapi "github.com/gaetanDubuc/beeckend/internal/cheptel/api"
	cheptelrepository "github.com/gaetanDubuc/beeckend/internal/cheptel/repository"
	cheptelservice "github.com/gaetanDubuc/beeckend/internal/cheptel/service"
	cheptelmngrepository "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/repository"
	cheptelmngservice "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/service"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/internal/router"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
)

func main() {
	logger := utils.NewLogger()

	config, err := utils.LoadConfig(".")
	if err != nil {
		logger.Fatal("cannot load config:", err)
		panic(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		logger.Error("Unable to connect to database:", err)
		panic(err)
	}

	db := db.NewGormWithMigrate(
		postgres.Open(config.DBSource),
		"file://./migrations",
		config.DatabaseURL,
		logger)

	router, v1 := router.New(db)
	RegisterHandlers(v1, db, pool, logger)
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

func RegisterHandlers(router *gin.RouterGroup, db *db.DB, pool *pgxpool.Pool, logger log.Logger) {
	var upgrader = &websocket.Upgrader{} // use default options

	// Repositories
	cheptelRepository := cheptelrepository.NewGormRepository(db, pool, logger)
	cheptelManagerRepository := cheptelmngrepository.NewGormRepository(db)

	// Services
	cheptelManager := cheptelmngservice.NewService(
		cheptelManagerRepository,
		logger,
	)
	cheptelService := cheptelservice.NewService(
		cheptelRepository,
		cheptelManager,
		cheptelManagerRepository,
		logger,
	)

	// APIs
	cheptelapi.RegisterHandlers(router, cheptelService, upgrader, logger)
}
