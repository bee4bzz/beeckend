package test

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"gorm.io/gorm"
)

var (
	ValidUser = entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "ValidUser",
		Email:    utils.ValidEmail(),
		Cheptels: []entity.Cheptel{ValidCheptel},
	}

	ValidCheptel = entity.Cheptel{
		Model: gorm.Model{
			ID: 2,
		},
		Name: "ValidCheptel",
		Hives: []entity.Hive{
			ValidHive,
		},
		Notes:  []entity.CheptelNote{},
		Albums: []entity.Album{},
	}

	ValidHive = entity.Hive{
		Model: gorm.Model{
			ID: 3,
		},
		Name:      "ValidHive",
		CheptelID: 2,
		Notes:     []entity.HiveNote{ValidHiveNote},
	}

	ValidHiveNote = entity.HiveNote{
		Model: gorm.Model{
			ID: 4,
		},
		HiveID:    3,
		Name:      "ValidHiveNote",
		Operation: utils.ValidName(),
	}
)
