package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Cheptels []Cheptel `gorm:"many2many:user_cheptels;"`
}

type Cheptel struct {
	gorm.Model
	Hives  []Hive
	Notes  []CheptelNote
	Albums []Album `gorm:"polymorphic:Owner;"`
}

type weather string

type CheptelNote struct {
	gorm.Model
	CheptelID        uint
	Name             string
	TemperatureDay   float64
	TemperatureNight float64
	Weather          weather `sql:"type:ENUM('CLOUDY', 'SUNNY', 'SNOWY', 'RAINY', 'WINDY', 'STORMY', 'FOGGY', 'HAZY', 'OTHER')"`
	Flora            string
	State            string
	Observation      string
}

type Hive struct {
	gorm.Model
	Name      string
	CheptelID uint
	Notes     []HiveNote
}

type HiveNote struct {
	gorm.Model
	HiveID      uint
	Name        string
	NBRisers    uint
	Operation   string
	Observation string
	Albums      []Album `gorm:"polymorphic:Owner;"`
}

type Album struct {
	gorm.Model
	Name        string
	paths       []string
	Observation string
	OwnerID     uint
	OwnerType   string
}
