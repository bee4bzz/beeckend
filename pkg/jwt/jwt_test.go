package jwt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	secret := "secret-key"
	{
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"ID": "100",
		})
		tokenString, err := tok.SignedString([]byte(secret))
		assert.Nil(t, err)

		h := JWT(&jwt.MapClaims{}, Params[*gin.Context]{
			keyFunc: func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
		})
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users/", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)
		c, _ := gin.CreateTestContext(res)
		c.Request = req
		err = h(c)
		assert.Nil(t, err)
		token, ok := c.Get("JWT")
		assert.True(t, ok)
		if assert.NotNil(t, token) {
			claim := token.(*jwt.Token).Claims.(jwt.MapClaims)
			assert.Equal(t, "100", claim["ID"])
		}
	}

}
