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
