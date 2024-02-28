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
		Name:  "ValidUser",
		Email: utils.ValidEmail(),
	}

	ValidCheptel = entity.Cheptel{
		Model: gorm.Model{
			ID: 2,
		},
		Name: "ValidCheptel",
		Hives: []entity.Hive{
			ValidHive,
		},
	}

	ValidHive = entity.Hive{
		Model: gorm.Model{
			ID: 3,
		},
		Name:      "ValidHive",
		CheptelID: 2,
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
