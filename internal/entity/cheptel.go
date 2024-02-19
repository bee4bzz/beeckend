package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

const (
	CheptelTableName = "Cheptels"
)

type Cheptel struct {
	gorm.Model
	Hives  []Hive
	Notes  []CheptelNote
	Albums []Album `gorm:"polymorphic:Owner;"`
}

// Validate Cheptel structure.
func (c Cheptel) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Hives, validation.Each(validation.Required)),
		validation.Field(&c.Notes, validation.Each(validation.Required)),
		validation.Field(&c.Albums, validation.Each(validation.Required)),
	)
}

type weather string

func (self *weather) Scan(value interface{}) error {
	*self = weather(value.(string))
	return nil
}

func (self weather) Value() (interface{}, error) {
	return string(self), nil
}

const (
	CLOUDY weather = "CLOUDY"
	SUNNY  weather = "SUNNY"
	SNOWY  weather = "SNOWY"
	RAINY  weather = "RAINY"
	WINDY  weather = "WINDY"
	STORMY weather = "STORMY"
	FOGGY  weather = "FOGGY"
	HAZY   weather = "HAZY"
	OTHER  weather = "OTHER"
)

type CheptelNote struct {
	gorm.Model
	CheptelID        uint   `gorm:"constraint:OnDelete:CASCADE;"`
	Name             string `gorm:"not null"`
	TemperatureDay   *float64
	TemperatureNight *float64
	Weather          *weather `sql:"type:ENUM('CLOUDY', 'SUNNY', 'SNOWY', 'RAINY', 'WINDY', 'STORMY', 'FOGGY', 'HAZY', 'OTHER')"`
	Flora            string   `gorm:"not null"`
	State            *string
	Observation      *string
}

// Validate CheptelNote structure.
func (c CheptelNote) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Weather, validation.NilOrNotEmpty, validation.In(CLOUDY, SUNNY, SNOWY, RAINY, WINDY, STORMY, FOGGY, HAZY, OTHER)),
		validation.Field(&c.Flora, validation.Required),
	)
}
