package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotNilUUID(t *testing.T) {
	err := NotNilUUID.Validate(uuid.Nil)
	assert.Error(t, err)
	err = NotNilUUID.Validate(uuid.New())
	assert.NoError(t, err)
}
