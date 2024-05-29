package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	cheptelapi "github.com/gaetanDubuc/beeckend/internal/cheptel/api"
	"github.com/gaetanDubuc/beeckend/internal/db"
	"github.com/gaetanDubuc/beeckend/internal/log"
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
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		logger.Error("Unable to connect to database:", err)
		panic(err)
	}

	go func() {
		conn, err := pool.Acquire(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
			panic(err)
		}
		defer conn.Release()

		_, err = conn.Exec(ctx, "LISTEN chat")
		if err != nil {
			logger.Fatal("cannot listen to chat:", err)
		}

		for {
			notification, err := conn.Conn().WaitForNotification(context.Background())
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error waiting for notification:", err)
				os.Exit(1)
			}

			fmt.Println("PID:", notification.PID, "Channel:", notification.Channel, "Payload:", notification.Payload)
		}
	}()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
		os.Exit(1)
	}
	_, err = conn.Conn().Exec(ctx, "NOTIFY chat, 'hello';")
	if err != nil {
		logger.Fatal("cannot listen to chat:", err)
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
	var upgrader = &websocket.Upgrader{} // use default options

	cheptelapi.RegisterHandlers(router, upgrader, logger)
}
