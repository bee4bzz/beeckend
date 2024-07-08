package schema

import (
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenJSONKey        = "token"
	RefreshTokenJSONKey = "refresh_token"
)

// LoginRequest.
type LoginRequest struct {
	Username string `json:"username" example:"user@example.com"`
	Password string `json:"password" example:"test12345"`
}

// Validate LoginRequest structure.
func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required),
	)
}

type LogoutRequest struct {
	UserID uint `json:"-" example:"00000000-0000-0000-0000-000000000000"`
}

// Validate LogoutRequest structure.
func (r LogoutRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
	)
}

// Session is the response after a successful login
// It contains the authentication JWT and the refresh JWT.
type Session struct {
	JWT        string `json:"jwt"`
	RefreshJWT string `json:"refresh_jwt"`
}

// PublicKeyResponse.
// It contains the PublicKey which is a PEM encoded public key.
type PublicKeyResponse struct {
	PublicKey string `json:"key"`
}

func MakeUserClaim(user entity.User, expiration int) *Claims {
	now := time.Now().UTC()
	return &Claims{
		UserID: user.ID,
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
	UserID uint `json:"user_ID"`
}

func (c *Claims) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.UserID, validation.Required),
	)
}
