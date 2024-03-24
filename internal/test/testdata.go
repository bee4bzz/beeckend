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
		Notes: []entity.CheptelNote{ValidCheptelNote},
		Albums: []entity.Album{
			ValidAlbum,
		},
	}

	ValidCheptelNote = entity.CheptelNote{
		Model: gorm.Model{
			ID: 1,
		},
		CheptelID: 2,
		Name:      "ValidCheptelNote",
		Flora:     "new flora",
		Weather:   entity.CLOUDY,
	}

	ValidAlbum = entity.Album{
		Model: gorm.Model{
			ID: 2,
		},
		Name:      "ValidAlbum",
		OwnerID:   2,
		OwnerType: "cheptels",
	}

	ValidHive = entity.Hive{
		Model: gorm.Model{
			ID: 3,
		},
		Name:      "ValidHive",
		CheptelID: 2,
		Notes:     []entity.HiveNote{ValidHiveNote},
	}

	ValidHive2 = entity.Hive{
		Model: gorm.Model{
			ID: 4,
		},
		Name:      "ValidHive2",
		CheptelID: 2,
		Notes:     []entity.HiveNote{},
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
