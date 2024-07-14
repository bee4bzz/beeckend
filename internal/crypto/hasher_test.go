package crypto

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	repo := NewHasher(10)
	t.Run("GenerateToken", func(t *testing.T) {
		token := repo.GenerateToken()
		if len(token) != 10 {
			t.Errorf("Expected token length to be 10, got %d", len(token))
		}
	})

	t.Run("hash", func(t *testing.T) {
		repo := NewHasher(10)
		hashedValue, err := repo.Hash(test.Password)

		assert.NoError(t, err)
		assert.Len(t, hashedValue, 60)
		assert.NotEqual(t, hashedValue, test.Password)
	})

	t.Run("hash", func(t *testing.T) {
		repo := NewHasher(10)
		hashedValue, _ := repo.Hash(test.Password)
		result := repo.AreSameHash(test.Password, hashedValue)

		assert.True(t, result)

		result = repo.AreSameHash(test.Password, test.Password)
		assert.False(t, result)
	})
}
