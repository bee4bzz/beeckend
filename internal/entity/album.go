package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Name        string   `gorm:"not null"`
	paths       []string `gorm:"not null"`
	Observation *string
	OwnerID     uint
	OwnerType   string
}

// Validate Album structure.
func (a Album) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.paths, validation.Each(validation.Required)),
	)
}
