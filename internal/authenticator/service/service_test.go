package service

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/crypto"
	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	print((&crypto.Hasher{32}).Hash("test12345."))
	assert.Equal(t, 1, 0)
}
