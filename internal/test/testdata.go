package test

import (
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	Password            = "password12345"
	HashedPassword, err = bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)

	// User has:
	// 2 cheptels
	//  - ValidCheptel
	//  - ValidCheptel2
	// 1 cheptel album
	//  - ValidCheptelAlbum
	// 1 cheptel note
	//  - ValidCheptelNote
	// 1 photo
	//  - ValidPhoto
	// 1 hive
	//  - ValidHive
	// 1 hive note
	//  - ValidHiveNote
	// 1 hive note album
	//  - ValidHiveNoteAlbum
	ValidUser = entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name:           "ValidUser",
		Email:          utils.ValidEmail(),
		HashedPassword: string(HashedPassword),
		Cheptels:       []entity.Cheptel{ValidCheptel, ValidCheptel2},
	}

	ValidUser2 = entity.User{
		Model: gorm.Model{
			ID: 2,
		},
		Name:           "ValidUser2",
		Email:          utils.ValidEmail(),
		HashedPassword: string(HashedPassword),
		Cheptels:       []entity.Cheptel{ValidCheptel},
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
		Albums: []entity.CheptelAlbum{
			ValidCheptelAlbum,
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

	ValidCheptelAlbum = entity.CheptelAlbum{
		Album: entity.Album{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			Name:    "ValidCheptelAlbum",
			OwnerID: 1,
		},
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

	ValidHiveNoteAlbum = entity.HiveNoteAlbum{
		Album: entity.Album{
			Model: gorm.Model{
				ID: 1,
			},
			Name:    "ValidHiveNoteAlbum",
			OwnerID: 1},
	}
)
