package schema

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	Type = "refresh"
)

type CreateRequest struct {
	UserID uint `json:"-" swaggertype:"string"`
}

func (m CreateRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required),
	)
}

type RefreshRequest struct {
	UserID  uint   `json:"-" swaggertype:"string"`
	TokenID uint   `json:"challenge,omitempty" form:"challenge"`
	Token   string `json:"token,omitempty" form:"token"`
}

func (m RefreshRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required),
		validation.Field(&m.TokenID, validation.Required),
		validation.Field(&m.Token, validation.Required),
	)
}

type DeleteFromUserRequest struct {
	UserID uint `json:"-" swaggertype:"string"`
}

func (m DeleteFromUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required),
	)
}

func MakeUserClaim(userID uint, expiration int, token string) *Claims {
	now := time.Now().UTC()
	return &Claims{
		UserID: userID,
		Type:   Type,
		Token:  token,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    "beeckend",
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiration) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"user_ID"`
	Type   string `json:"type"`
	Token  string `json:"token"`
}

func (c *Claims) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.UserID, validation.Required),
		validation.Field(&c.Type, validation.In(Type)),
		validation.Field(&c.Token, validation.Required),
	)
}
