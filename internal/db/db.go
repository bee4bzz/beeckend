package db

import (
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"gorm.io/gorm"
)

var (
	GormConfig = gorm.Config{
		NowFunc: func() time.Time {
			// Spécifier la localisation temporelle que vous souhaitez utiliser
			return time.Now().UTC() // Par exemple, UTC
		},
	}
)

func InitGorm(dial gorm.Dialector) *gorm.DB {
	db, err := gorm.Open(dial, &GormConfig)

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&entity.User{}, &entity.Cheptel{}, &entity.Hive{}, &entity.HiveNote{}, &entity.CheptelNote{}, &entity.Album{})
	return db
}
