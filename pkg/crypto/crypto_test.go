package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyToPem(t *testing.T) {
	{
		ValidRSASigningKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			panic(err)
		}

		p, err := PublicKeyToPem(ValidRSASigningKey.Public())
		assert.NoError(t, err)

		pb, err := jwt.ParseRSAPublicKeyFromPEM(p)
		assert.NoError(t, err)

		assert.Equal(t, ValidRSASigningKey.Public(), pb)
	}
	{
		_, err := PublicKeyToPem("invalid")
		assert.Error(t, err)
	}
}
