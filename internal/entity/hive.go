package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

const (
	HiveNotesKey      = "Notes"
	MaxNBRisers  uint = 5
)

type Hive struct {
	gorm.Model
	Name      string     `gorm:"index:idx_name_cheptel_id,unique;not null"`
	CheptelID uint       `gorm:"index:idx_name_cheptel_id,unique;not null"`
	Notes     []HiveNote `gorm:"constraint:OnDelete:CASCADE;"`
}

// Validate User structure.
func (h Hive) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Name, validation.Required),
		validation.Field(&h.CheptelID, validation.Required),
		validation.Field(&h.Notes),
	)
}

type HiveNote struct {
	gorm.Model
	HiveID      uint   `gorm:"index:idx_name_hive_id,unique;not null"`
	Name        string `gorm:"index:idx_name_hive_id,unique;not null"`
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
		validation.Field(&h.NBRisers, validation.Max(MaxNBRisers)),
		validation.Field(&h.Operation, validation.Required),
		validation.Field(&h.Observation, validation.NilOrNotEmpty),
		validation.Field(&h.Albums),
	)
}
