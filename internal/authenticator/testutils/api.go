package testutils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/path"
	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/golang-jwt/jwt"
)

func BearerAuthHeader(header *http.Header, token string) {
	header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

var (
	LoginRootTest = test.APITestCase[schema.Session]{
		Method: utils.String("POST"),
		URL: &url.URL{
			Scheme: "http",
		},
		WantStatus:   utils.Int(http.StatusOK),
		WantResponse: utils.String(".*token.*refresh_token.*"),
	}

	RefreshSessionRootTest = test.APITestCase[schema.Session]{
		Method:       utils.String("POST"),
		URL:          utils.String(path.RefreshSessionPath),
		WantStatus:   utils.Int(http.StatusOK),
		WantResponse: utils.String(".*token.*refresh_token.*"),
	}

	PublicKeyRootTest = test.APITestCase[schema.PublicKeyResponse]{
		Method:     utils.String("GET"),
		URL:        utils.String(path.PublicKeyPath),
		WantStatus: utils.Int(http.StatusOK),
	}

	LogoutRootTest = test.APITestCase[any]{
		Method:     utils.String("DELETE"),
		URL:        utils.String(path.LogoutPath),
		WantStatus: utils.Int(http.StatusOK),
	}
)

func ExtractAuthenticationJWT[V any](tokenString string, verificationKey V, signinMethod string) string {
	parser := &jwt.Parser{
		ValidMethods: []string{signinMethod},
	}
	token, err := parser.Parse(tokenString, func(t *jwt.Token) (interface{}, error) { return verificationKey, nil })
	if err == nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		return claims["UUID"].(string)
	}
	return ""
}

type Client struct {
	http.Client
	BaseURL         string
	Email, Password string
	Header          *http.Header
	Logger          log.Logger
}

func (c *Client) Login() (schema.Session, error) {
	response, err := LoginRootTest.
		WithRequest(schema.LoginRequest{
			Username: c.Email,
			Password: c.Password,
		}).
		WithURL(c.URL).
		WithLogger(c.Logger).
		WithClient(&c.Client).Call()

	if err != nil {
		return schema.Session{}, err
	}

	BearerAuthHeader(c.Header, response.Token)

	return response, nil
}
