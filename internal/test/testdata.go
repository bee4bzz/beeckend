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
		Name:     utils.ValidName(),
		Email:    utils.ValidEmail(),
		Cheptels: []entity.Cheptel{},
	}

	ValidCheptel = entity.Cheptel{
		Model: gorm.Model{
			ID: 1,
		},
		Name:  utils.ValidName(),
		Users: []entity.User{},
	}

	ValidHive = entity.Hive{
		Model: gorm.Model{
			ID: 1,
		},
		Name:      utils.ValidName(),
		CheptelID: 1,
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
