package schema

import (
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt/v5"
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
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// PublicKeyResponse.
// It contains the PublicKey which is a PEM encoded public key.
type PublicKeyResponse struct {
	PublicKey string `json:"key"`
}

func MakeUserClaim(user entity.User, expiration int) *Claims {
	return &Claims{
		ID:  user.ID,
		Exp: time.Now().UTC().Add(time.Duration(expiration) * time.Second).Unix(),
	}
}

type Claims struct {
	jwt.RegisteredClaims
	ID  uint  `json:"ID"`
	Exp int64 `json:"exp"`
}

func (c *Claims) Valid() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ID, validation.Required),
		validation.Field(&c.Exp, validation.Required, validation.Min(time.Now().Unix()).Exclusive().Error("Token is expired")),
	)
}

type ExpiredClaims struct {
	jwt.RegisteredClaims
	Claims
}

func (c *ExpiredClaims) Valid() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Claims.ID, validation.Required),
		validation.Field(&c.Exp, validation.Required, validation.Max(
			time.Now().Unix()).Error("Token is not expired")),
	)
}
