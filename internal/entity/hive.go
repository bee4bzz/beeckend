package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type Hive struct {
	gorm.Model
	Name      string `gorm:"not null"`
	CheptelID uint
	Notes     []HiveNote `gorm:"constraint:OnDelete:CASCADE;"`
}

// Validate User structure.
func (h Hive) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Name, validation.Required),
		validation.Field(&h.CheptelID, validation.Required),
		validation.Field(&h.Notes, validation.Each(validation.Required)),
	)
}

type HiveNote struct {
	gorm.Model
	HiveID      uint
	Name        string `gorm:"not null"`
	NBRisers    uint   `gorm:"not null"`
	Operation   string `gorm:"not null"`
	Observation *string
	Albums      []Album `gorm:"polymorphic:Owner;"`
}

// Validate User structure.
func (h HiveNote) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.HiveID, validation.Required),
		validation.Field(&h.Name, validation.Required),
		validation.Field(&h.NBRisers, validation.Max(5)),
		validation.Field(&h.Operation, validation.Required),
		validation.Field(&h.Observation, validation.NilOrNotEmpty),
		validation.Field(&h.Albums, validation.Each(validation.Required)),
	)
}
