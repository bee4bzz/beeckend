package db

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"gorm.io/gorm"
)

// TODO: implment these functions in test package
func Seed(t *testing.T, db *gorm.DB, values ...any) {
	t.Helper()
	err := db.Exec("PRAGMA foreign_keys = ON", nil).Error
	if err != nil {
		t.Fatal(err)
	}
	for _, ptr := range values {
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
