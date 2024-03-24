package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type Album struct {
	gorm.Model
	Name        string `gorm:"index:idx_name_owner_id,unique;not null"`
	Observation *string
	OwnerID     uint    `gorm:"index:idx_name_owner_id,unique;not null"`
	OwnerType   string  `gorm:"not null"`
	Photos      []Photo `gorm:"constraint:OnDelete:CASCADE;"`
}

// Validate Album structure.
func (a Album) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Observation, validation.NilOrNotEmpty),
		validation.Field(&a.OwnerID, validation.Required),
		validation.Field(&a.OwnerType, validation.Required),
	)
}

type Photo struct {
	gorm.Model
	Path    string `gorm:"not null"`
	AlbumID uint   `gorm:"not null"`
}

// Validate Photo structure.
func (p Photo) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Path, validation.Required),
		validation.Field(&p.AlbumID, validation.Required),
	)
}
