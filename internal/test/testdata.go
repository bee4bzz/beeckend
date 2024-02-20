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
		Name:  utils.ValidName(),
		Email: utils.ValidEmail(),
		Cheptels: []entity.Cheptel{
			ValidCheptel,
		},
	}

	ValidCheptel = entity.Cheptel{
		Model: gorm.Model{
			ID: 1,
		},
		Name: utils.ValidName(),
		Hives: []entity.Hive{
			ValidHive,
		},
	}

	ValidHive = entity.Hive{
		Model: gorm.Model{
			ID: 1,
		},
		Name:      utils.ValidName(),
		CheptelID: 1,
		Notes:     []entity.HiveNote{ValidHiveNote},
	}

	ValidHiveNote = entity.HiveNote{
		Model: gorm.Model{
			ID: 1,
		},
		HiveID:    1,
		Name:      utils.ValidName(),
		Operation: utils.ValidName(),
	}
)
