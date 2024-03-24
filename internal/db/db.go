package db

import (
	"log"
	"os"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	l "github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	zaplog "github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewLogger() l.Logger {
	config, err := utils.LoadConfig(".")

	if err != nil {
		panic("failed to load config")
	}

	var logger l.Logger
	if config.AppEnv == "development" {
		logger = zaplog.NewProduction()
	} else {
		logger = zaplog.NewDevelopment()
	}
	return logger
}

func NewGorm(dial gorm.Dialector, logger logger.Interface) *gorm.DB {
	GormConfig := gorm.Config{
		TranslateError: true,
		NowFunc: func() time.Time {
			// Spécifier la localisation temporelle que vous souhaitez utiliser
			return time.Now().UTC() // Par exemple, UTC
		},
		Logger: logger,
	}
	db, err := gorm.Open(dial, &GormConfig)

	if err != nil {
		panic("failed to connect database")
	}

	return db.Session(&gorm.Session{})
}

func NewGormForTest(dial gorm.Dialector) *gorm.DB {
	db := NewGorm(dial, logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{},
	))

	err := db.AutoMigrate(&entity.User{}, &entity.Cheptel{}, &entity.Hive{}, &entity.CheptelNote{}, &entity.HiveNote{}, &entity.Album{}, &entity.Photo{})
	if err != nil {
		panic("failed to migrate " + err.Error())
	}
	return db
}

func NewGormWithMigrate(dial gorm.Dialector, sourceURL, databaseURL string, log l.Logger) *gorm.DB {
	db := NewGorm(dial, logger.Default)

	// make migration programmaticaly
	m, err := migrate.New(
		sourceURL,
		databaseURL)

	if err != nil {
		log.Error("failed to create a new migrate instance: ", err)
	}
	if err := m.Up(); err != nil {
		log.Error("failed to migrate up: ", err)
	}
	return db
}
