package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

const (
	UsersKey  = "Users"
	HivesKey  = "Hives"
	AlbumsKey = "Albums"
)

type Cheptel struct {
	gorm.Model
	Name   string        `gorm:"not null"`
	Hives  []Hive        `gorm:"constraint:OnDelete:CASCADE;"`
	Notes  []CheptelNote `gorm:"constraint:OnDelete:CASCADE;"`
	Albums []Album       `gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE; "`
	Users  []User        `gorm:"many2many:user_cheptels;constraint:OnDelete:CASCADE;"`
}

// Validate Cheptel structure.
func (c Cheptel) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Hives),
		validation.Field(&c.Notes),
		validation.Field(&c.Albums),
	)
}
