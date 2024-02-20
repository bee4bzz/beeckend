package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string    `gorm:"not null"`
	Email    string    `gorm:"unique;not null"`
	Cheptels []Cheptel `gorm:"many2many:user_cheptels;constraint:OnDelete:CASCADE;"`
}

// Validate User structure.
func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Cheptels),
	)
}
