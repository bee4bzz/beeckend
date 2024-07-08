package testutils

import (
	"net/http"
	"net/url"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/path"
	"github.com/gaetanDubuc/beeckend/internal/authenticator/schema"
	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/golang-jwt/jwt"
)

var (
	LoginRootTest = test.APITestCase[schema.Session]{
		Method: utils.String("POST"),
		URL: &url.URL{
			Path: path.LoginPath,
		},
		WantStatus:   utils.Int(http.StatusOK),
		WantResponse: utils.String(".*token.*refresh_token.*"),
	}

	RefreshSessionRootTest = test.APITestCase[schema.Session]{
		Method:       utils.String("POST"),
		URL:          &url.URL{Path: path.RefreshSessionPath},
		WantStatus:   utils.Int(http.StatusOK),
		WantResponse: utils.String(".*token.*refresh_token.*"),
	}

	PublicKeyRootTest = test.APITestCase[schema.PublicKeyResponse]{
		Method:     utils.String("GET"),
		URL:        &url.URL{Path: path.PublicKeyPath},
		WantStatus: utils.Int(http.StatusOK),
	}

	LogoutRootTest = test.APITestCase[any]{
		Method:     utils.String("DELETE"),
		URL:        &url.URL{Path: path.LogoutPath},
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
	scheme          string
	host            string
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
		WithScheme(c.scheme).
		WithHost(c.host).
		WithLogger(c.Logger).
		WithClient(&c.Client).Call()

	if err != nil {
		return schema.Session{}, err
	}

	utils.BearerAuthHeader(c.Header, response.JWT)

	return response, nil
}
