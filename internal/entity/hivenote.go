package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

const (
	MaxNBRisers uint = 5
)

type HiveNote struct {
	gorm.Model
	HiveID      uint   `gorm:"index:idx_name_hive_id,unique;not null"`
	Hive        Hive   `gorm:"foreignKey:HiveID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name        string `gorm:"index:idx_name_hive_id,unique;not null"`
	NBRisers    uint   `gorm:"not null"`
	Operation   string `gorm:"not null"`
	Observation *string
	Albums      []HiveNoteAlbum `gorm:"constraint:OnDelete:CASCADE;foreignKey:OwnerID"`
}

// Validate HiveNote structure.
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
