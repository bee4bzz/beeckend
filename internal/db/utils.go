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
		&entity.CheptelPhoto{},
		&entity.Hive{},
		&entity.CheptelNote{},
		&entity.HiveNote{},
		&entity.HiveNoteAlbum{},
		&entity.HiveNotePhoto{},
	}
)

func NewGormForTest(dial gorm.Dialector) *DB {
	db := NewGorm(dial, logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{},
	))

	err := db.DB().AutoMigrate(
		tables...,
	)
	if err != nil {
		panic("failed to migrate " + err.Error())
	}
	return db
}

// TODO: implment these functions in test package
func Seed(t *testing.T, db *DB, values ...any) {
	t.Helper()

	for _, ptr := range values {
		err := db.DB().Clauses(clause.OnConflict{DoNothing: true}).Create(ptr).Error
		if err != nil {
			t.Fatal(err)
		}
	}

	// Delete sequences
	script := `
	DO $$ 
	DECLARE
		r RECORD;
		seq_name text;
		tbl_name text;
		max_id int;
	BEGIN
		FOR r IN (
			SELECT c.relname AS seq_name, 
				substring(c.relname from 1 for length(c.relname) - 7) AS tbl_name
			FROM pg_class c
			JOIN pg_namespace n ON n.oid = c.relnamespace
			WHERE c.relkind = 'S' AND n.nspname = 'public'
		) LOOP
			-- Determine the maximum ID currently in the corresponding table
			EXECUTE 'SELECT MAX(id) FROM ' || r.tbl_name INTO max_id;

			-- If there are no records, start with 1, otherwise start with max_id + 1
			IF max_id IS NULL THEN
				max_id := 1;
			ELSE
				max_id := max_id + 1;
			END IF;

			-- Restart the sequence with the determined max_id
			EXECUTE 'ALTER SEQUENCE ' || r.seq_name || ' RESTART WITH ' || max_id;
		END LOOP;
	END $$;
	`

	// Exécution du script
	if err := db.db.Exec(script).Error; err != nil {
		log.Fatalf("Error executing script: %v", err)
	}
}

func Clean(t *testing.T, db *DB) {
	t.Helper()
	for _, ptr := range tables {
		err := db.DB().Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(ptr).Error
		if err != nil {
			t.Fatal(err)
		}
	}
}
