package db

import (
	"time"

	l "github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(dial gorm.Dialector, logger logger.Interface) *DB {
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

	return New(db.Session(&gorm.Session{}))
}

func NewGormWithMigrate(dial gorm.Dialector, sourceURL, databaseURL string, log l.Logger) *DB {
	db := NewGorm(dial, logger.Default)

	// make migration programmaticaly
	m, err := migrate.New(
		sourceURL,
		databaseURL)

	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		log.Info("failed to migrate up: ", err)
	}
	return db
}
