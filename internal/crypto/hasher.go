package crypto

import (
	"github.com/labstack/gommon/random"
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	TokenLength uint8
}

func NewHasher(TokenLength uint8) *Hasher {
	return &Hasher{TokenLength}
}

func (r *Hasher) GenerateToken() string {
	return random.String(r.TokenLength, random.Alphanumeric)
}

func (r *Hasher) Hash(value string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	return string(hashedValue), err
}

func (r *Hasher) AreSameHash(value string, hashedValue string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	return err == nil
}
