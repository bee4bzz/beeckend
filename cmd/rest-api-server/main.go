package main

import (
	"context"
	"net/http"
	"os"

	authapi "github.com/gaetanDubuc/beeckend/internal/authenticator/api"
	authmiddleware "github.com/gaetanDubuc/beeckend/internal/authenticator/middleware"
	authservice "github.com/gaetanDubuc/beeckend/internal/authenticator/service"
	cheptelapi "github.com/gaetanDubuc/beeckend/internal/cheptel/api"
	cheptelrepository "github.com/gaetanDubuc/beeckend/internal/cheptel/repository"
	cheptelservice "github.com/gaetanDubuc/beeckend/internal/cheptel/service"
	cheptelmngrepository "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/repository"
	cheptelmngservice "github.com/gaetanDubuc/beeckend/internal/cheptelmanager/service"
	"github.com/gaetanDubuc/beeckend/internal/crypto"
	"github.com/gaetanDubuc/beeckend/internal/db"
	refreshtokrepository "github.com/gaetanDubuc/beeckend/internal/refresh-token/repository"
	refreshtokservice "github.com/gaetanDubuc/beeckend/internal/refresh-token/service"
	"github.com/gaetanDubuc/beeckend/internal/router"
	tokenservice "github.com/gaetanDubuc/beeckend/internal/token"
	userrepository "github.com/gaetanDubuc/beeckend/internal/user/repository"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	crypt "github.com/gaetanDubuc/beeckend/pkg/crypto"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
)

// @title						Beeckend API
// @version					NA
// @BasePath					/api/v1
// @securityDefinitions.apikey	JWT Token
// @in							header
// @name						Authorization
// @securityDefinitions.basic	BasicAuth
// @Description				JWT Token: You must precede the token with `Bearer `
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

	router, v1 := router.New(os.Stdout, db)
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

func RegisterHandlers(router *gin.RouterGroup, db *db.DB, pool *pgxpool.Pool, logger *log.Logger) {
	// Utilities
	upgrader := &websocket.Upgrader{} // use default options
	hasher := crypto.NewHasher(32)

	// Load private and public keys
	fb := "key"
	publicKey, privateKey, err := crypt.LoadOrGenerateKeys(fb)

	if err != nil {
		panic(err)
	}

	// Repositories
	refreshTokenRepository := refreshtokrepository.New(db)
	userRepository := userrepository.NewGormRepository(db)
	cheptelRepository := cheptelrepository.NewGormRepository(db, pool, logger)
	cheptelManagerRepository := cheptelmngrepository.NewGormRepository(db)

	// Services
	tokenService := tokenservice.New(refreshTokenRepository, hasher, logger)
	refreshTokenService := refreshtokservice.New(
		tokenService,
		hasher,
		86400,
		jwt.SigningMethodRS256,
		func(t *jwt.Token) (interface{}, error) {
			return privateKey, nil
		},
		logger)
	authService := authservice.New(
		userRepository,
		refreshTokenService,
		hasher,
		3600,
		jwt.SigningMethodRS256,
		func(t *jwt.Token) (interface{}, error) {
			return privateKey, nil
		},
		logger)

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

	// Middlewares
	authMiddleware := authmiddleware.New(jwt.SigningMethodRS256.Name,
		func(t *jwt.Token) (interface{}, error) {
			return publicKey, nil
		}, userRepository, logger)

	// APIs
	authapi.RegisterHandlers(router, authService, authMiddleware, "", logger)
	cheptelapi.RegisterHandlers(router, cheptelService, authMiddleware, upgrader, logger)
}
