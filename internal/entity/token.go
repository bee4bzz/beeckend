package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gorm.io/gorm"
)

type TokenType string

const (
	RefreshTokenType TokenType = "refreshtokens"
)

type RefreshToken struct {
	gorm.Model
	Token
}

func (r RefreshToken) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Token, validation.Required),
	)
}

func (r *RefreshToken) GetCreatedAt() time.Time {
	return r.CreatedAt
}

func (r *RefreshToken) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}

// Token is the data stored to challenge its owner and to verify its informations.
type Token struct {
	OwnerID     uint   `json:"owner_ID,omitempty" mapstructure:"owner_ID" example:"5a4fdee6-14d7-4929-8240-6cfb0b5a7456"`
	HashedToken string `json:"hashed_token" mapstructure:"-" example:"00812eb398682a7d492bd556132644ad279267f049344466f6f9cac2284fca4d"`
	Token       string `json:"-" mapstructure:"-" gorm:"-" example:"00812eb398682a7d492bd556132644ad279267f049344466f6f9cac2284fca4d"`
}

// Validate Token structure.
func (t Token) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.OwnerID, validation.Required),
		validation.Field(&t.HashedToken, validation.Required),
		validation.Field(&t.Token, validation.Required),
	)
}

func (t *Token) GetHashedToken() string {
	return t.HashedToken
}

func (t *Token) SetHashedToken(hashedToken string) {
	t.HashedToken = hashedToken
}

func (t *Token) GetToken() string {
	return t.Token
}

func (t *Token) SetToken(token string) {
	t.Token = token
}
