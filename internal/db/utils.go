package db

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"gorm.io/gorm"
)

func Seed(t *testing.T, db *gorm.DB) {
	t.Helper()
	err := db.Exec("PRAGMA foreign_keys = ON", nil).Error
	if err != nil {
		t.Fatal(err)
	}
	for _, ptr := range []any{&test.ValidUsers} {
		err := db.Create(ptr).Error
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Clean(t *testing.T, db *gorm.DB) {
	t.Helper()
	for _, ptr := range []any{&entity.User{}, &entity.Cheptel{}, &entity.Hive{}, &entity.HiveNote{}, &entity.CheptelNote{}, &entity.Album{}} {
		err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(ptr).Error
		if err != nil {
			t.Fatal(err)
		}
	}
}
