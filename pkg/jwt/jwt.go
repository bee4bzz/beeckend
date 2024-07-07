package jwt

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Context interface {
	Set(key string, value any)
	GetHeader(key string) string
	Header(key string, value string)
}

type JWTTokenHandler[T Context] func(c T, token *jwt.Token) error

func DefaultJWTTokenHandler[C Context](c C, token *jwt.Token) error {
	c.Set("JWT", token)
	return nil
}

// Params represents the parameters that are be used with the JWT handler.
type Params[C Context] struct {
	// auth realm. Defaults to "API".
	Realm string
	// the allowed signing method. This is required and should be the actual method that you use to create JWT token.
	// It defaults to "HS256".
	SigningMethod string
	// a function that handles the parsed JWT token. Defaults to DefaultJWTTokenHandler,
	// which stores the token in the context with the key "JWT".
	TokenHandler JWTTokenHandler[C]
	// a function to get a dynamic VerificationKey
	Keyfunc func(t *jwt.Token) (interface{}, error)
}

func JWT[C Context](claims jwt.Claims, p Params[C]) func(C) error {
	if p.Realm == "" {
		p.Realm = "API"
	}
	if p.SigningMethod == "" {
		p.SigningMethod = "HS256"
	}
	if p.TokenHandler == nil {
		p.TokenHandler = DefaultJWTTokenHandler
	}
	parser := jwt.NewParser(
		jwt.WithExpirationRequired(),
		jwt.WithStrictDecoding(),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{p.SigningMethod}),
	)
	return func(c C) error {
		header := c.GetHeader("Authorization")
		message := ""
		if strings.HasPrefix(header, "Bearer ") {
			token, err := parser.ParseWithClaims(
				header[7:],
				claims,
				p.Keyfunc,
			)
			if err == nil && token.Valid {
				err = p.TokenHandler(c, token)
			}
			if err == nil {
				return nil
			}
			message = err.Error()
		}

		c.Header("WWW-Authenticate", `Bearer realm="`+p.Realm+`"`)

		return errors.New(message)
	}
}
