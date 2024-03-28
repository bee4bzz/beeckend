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
			ValidChetpelAlbum,
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

	ValidChetpelAlbum = entity.Album{
		Model: gorm.Model{
			ID: 2,
		},
		Name:      "ValidChetpelAlbum",
		OwnerID:   2,
		OwnerType: "cheptels",
	}

	ValidPhoto = entity.Photo{
		Model: gorm.Model{
			ID: 3,
		},
		AlbumID: 2,
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

	ValidHiveNoteAlbum = entity.Album{
		Model: gorm.Model{
			ID: 3,
		},
		Name:      "ValidHiveNoteAlbum",
		OwnerID:   4,
		OwnerType: "hive_notes",
	}
)
