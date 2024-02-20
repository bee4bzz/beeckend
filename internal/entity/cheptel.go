package entity

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

const (
	CheptelTableName = "Cheptels"
)

type Cheptel struct {
	gorm.Model
	Name   string        `gorm:"not null"`
	Hives  []Hive        `gorm:"constraint:OnDelete:CASCADE;"`
	Notes  []CheptelNote `gorm:"constraint:OnDelete:CASCADE;"`
	Albums []Album       `gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE; "`
}

// Validate Cheptel structure.
func (c Cheptel) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Hives, validation.Each(validation.Required)),
		validation.Field(&c.Notes, validation.Each(validation.Required)),
		validation.Field(&c.Albums, validation.Each(validation.Required)),
	)
}

var (
	ErrInvalidData = errors.New("invalid data")
)

type weather string

func (self *weather) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return ErrInvalidData
	}
	*self = weather(v)
	return nil
}

func (self weather) Value() (interface{}, error) {
	return string(self), nil
}

const (
	CLOUDY  weather = "CLOUDY"
	SUNNY   weather = "SUNNY"
	SNOWY   weather = "SNOWY"
	RAINY   weather = "RAINY"
	WINDY   weather = "WINDY"
	STORMY  weather = "STORMY"
	FOGGY   weather = "FOGGY"
	HAZY    weather = "HAZY"
	OTHER   weather = "OTHER"
	UNKNOWN weather = "UNKNOWN"
)

type CheptelNote struct {
	gorm.Model
	CheptelID        uint
	Name             string `gorm:"not null"`
	TemperatureDay   *float64
	TemperatureNight *float64
	Weather          weather `gorm:"default:'UNKNOWN'"`
	Flora            string  `gorm:"not null"`
	State            *string
	Observation      *string
}

// Validate CheptelNote structure.
func (c CheptelNote) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.CheptelID, validation.Required),
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Weather, validation.Required, validation.In(CLOUDY, SUNNY, SNOWY, RAINY, WINDY, STORMY, FOGGY, HAZY, OTHER, UNKNOWN)),
		validation.Field(&c.Flora, validation.Required),
		validation.Field(&c.State, validation.NilOrNotEmpty),
		validation.Field(&c.Observation, validation.NilOrNotEmpty),
	)
}
