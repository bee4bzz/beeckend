package test

import (
	"time"

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
		Cheptels: []entity.Cheptel{ValidCheptel, ValidCheptel2},
	}

	ValidUser2 = entity.User{
		Model: gorm.Model{
			ID: 2,
		},
		Name:     "ValidUser2",
		Email:    utils.ValidEmail(),
		Cheptels: []entity.Cheptel{ValidCheptel},
	}

	ValidUsers = []entity.User{ValidUser, ValidUser2}

	ValidCheptel = entity.Cheptel{
		Model: gorm.Model{
			ID: 1,
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

	ValidCheptel2 = entity.Cheptel{
		Model: gorm.Model{
			ID: 2,
		},
		Name: "ValidCheptel2",
		Hives: []entity.Hive{
			ValidHive2,
		},
	}

	ValidCheptelNote = entity.CheptelNote{
		Model: gorm.Model{
			ID: 1,
		},
		CheptelID: 1,
		Name:      "ValidCheptelNote",
		Flora:     "new flora",
		Weather:   entity.CLOUDY,
	}

	ValidChetpelAlbum = entity.Album{
		Model: gorm.Model{
			ID: 1,
		},
		Name:      "ValidChetpelAlbum",
		OwnerID:   1,
		OwnerType: "cheptels",
	}

	ValidPhoto = entity.Photo{
		Model: gorm.Model{
			ID: 1,
		},
		AlbumID: 1,
	}

	ValidHive = entity.Hive{
		Model: gorm.Model{
			ID: 1,
		},
		Name:      "ValidHive",
		CheptelID: 1,
		Notes:     []entity.HiveNote{ValidHiveNote},
	}

	ValidHive2 = entity.Hive{
		Model: gorm.Model{
			ID: 2,
		},
		Name:      "ValidHive2",
		CheptelID: 2,
		Notes:     []entity.HiveNote{},
	}

	ValidHiveNote = entity.HiveNote{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
		HiveID:    1,
		Name:      "ValidHiveNote",
		Operation: utils.ValidName(),
	}

	ValidHiveNoteAlbum = entity.Album{
		Model: gorm.Model{
			ID: 1,
		},
		Name:      "ValidHiveNoteAlbum",
		OwnerID:   1,
		OwnerType: "hive_notes",
	}
)
