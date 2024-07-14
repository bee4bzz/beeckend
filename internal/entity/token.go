package entity

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type TokenType string

const (
	RefreshTokenType TokenType = "refreshtokens"
)

var (
	TokenTypes = []any{
		RefreshTokenType,
	}
)

// Token is the data stored to challenge its owner and to verify its informations.
type Token struct {
	gorm.Model
	OwnerID     uint      `json:"owner_ID,omitempty" mapstructure:"owner_ID" example:"5a4fdee6-14d7-4929-8240-6cfb0b5a7456"`
	OwnerType   TokenType `json:"-" mapstructure:"-" example:"devicesecurity" gorm:"-"`
	HashedToken string    `json:"hashed_token" mapstructure:"-" example:"00812eb398682a7d492bd556132644ad279267f049344466f6f9cac2284fca4d"`
	Token       string    `json:"-" mapstructure:"-" gorm:"-" example:"00812eb398682a7d492bd556132644ad279267f049344466f6f9cac2284fca4d"`
}

func (t Token) TableName() string {
	return string(t.OwnerType)
}

// Validate Token structure.
func (t Token) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.OwnerID, validation.Required),
		validation.Field(&t.OwnerType, validation.Required, validation.In(TokenTypes...)),
		validation.Field(&t.HashedToken, validation.Required),
		validation.Field(&t.Token, validation.Required),
	)
}
