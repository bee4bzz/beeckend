package db

import (
	"log"
	"os"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var (
	tables = []interface{}{
		&entity.User{},
		&entity.Cheptel{},
		&entity.CheptelAlbum{},
		&entity.Hive{},
		&entity.CheptelNote{},
		&entity.HiveNote{},
		&entity.HiveNoteAlbum{},
		&entity.Photo{},
	}
)

func NewGormForTest(dial gorm.Dialector) *gorm.DB {
	db := NewGorm(dial, logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{},
	))

	err := db.AutoMigrate(
		tables...,
	)
	if err != nil {
		panic("failed to migrate " + err.Error())
	}
	return db
}

// TODO: implment these functions in test package
func Seed(t *testing.T, db *gorm.DB, values ...any) {
	t.Helper()
	err := db.Exec("PRAGMA foreign_keys = ON", nil).Error
	if err != nil {
		t.Fatal(err)
	}
	for _, ptr := range values {
		err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(ptr).Error
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Clean(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, ptr := range tables {
		err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(ptr).Error
		if err != nil {
			t.Fatal(err)
		}
	}
}
