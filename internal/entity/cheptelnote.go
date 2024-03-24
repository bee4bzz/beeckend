package entity

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

var (
	ErrInvalidData  = errors.New("invalid data")
	CheptelNotesKey = "Notes"
)

const (
	CLOUDY  Weather = "CLOUDY"
	SUNNY   Weather = "SUNNY"
	SNOWY   Weather = "SNOWY"
	RAINY   Weather = "RAINY"
	WINDY   Weather = "WINDY"
	STORMY  Weather = "STORMY"
	FOGGY   Weather = "FOGGY"
	HAZY    Weather = "HAZY"
	OTHER   Weather = "OTHER"
	UNKNOWN Weather = "UNKNOWN"
)

type CheptelNote struct {
	gorm.Model
	CheptelID        uint
	Name             string `gorm:"not null"`
	TemperatureDay   *float64
	TemperatureNight *float64
	Weather          Weather `gorm:"default:'UNKNOWN'"`
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

type Weather string

func (self *Weather) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return ErrInvalidData
	}
	*self = Weather(v)
	return nil
}

func (self Weather) Value() (interface{}, error) {
	return string(self), nil
}
