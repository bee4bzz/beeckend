package db

import (
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	zaplog "github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func NewLogger() log.Logger {
	config, err := utils.LoadConfig(".")

	if err != nil {
		panic("failed to load config")
	}

	var logger log.Logger
	if config.AppEnv == "development" {
		logger = zaplog.NewProduction()
	} else {
		logger = zaplog.NewDevelopment()
	}
	return logger
}

func NewGorm(dial gorm.Dialector, logger log.Logger) *gorm.DB {
	GormConfig := gorm.Config{
		NowFunc: func() time.Time {
			// Spécifier la localisation temporelle que vous souhaitez utiliser
			return time.Now().UTC() // Par exemple, UTC
		},
	}
	db, err := gorm.Open(dial, &GormConfig)

	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func NewGormAutoMigrate(dial gorm.Dialector, logger log.Logger) *gorm.DB {
	db := NewGorm(dial, logger)

	db.AutoMigrate(&entity.User{}, &entity.Cheptel{}, &entity.Hive{}, &entity.CheptelNote{}, &entity.HiveNote{}, &entity.Album{})
	return db
}

func NewGormWithMigrate(dial gorm.Dialector, sourceURL, databaseURL string, logger log.Logger) *gorm.DB {
	db := NewGorm(dial, logger)

	// make migration programmaticaly
	m, err := migrate.New(
		sourceURL,
		databaseURL)

	if err != nil {
		logger.Error("failed to create a new migrate instance: ", err)
	}
	if err := m.Up(); err != nil {
		logger.Error("failed to migrate up: ", err)
	}
	return db
}