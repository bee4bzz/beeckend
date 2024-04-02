package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

const (
	HiveNotesKey = "Notes"
)

type Hive struct {
	gorm.Model
	Name      string     `gorm:"index:idx_name_cheptel_id,unique;not null"`
	CheptelID uint       `gorm:"index:idx_name_cheptel_id,unique;not null"`
	Cheptel   Cheptel    `gorm:"foreignKey:CheptelID;constraint:OnDelete:CASCADE;"`
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
