package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func AssertUser(t *testing.T, expected, actual entity.User) {
	assert.NotEmpty(t, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	assert.Equal(t, expected.Email, actual.Email, "Email should be equal")
	for i, expectedCheptel := range expected.Cheptels {
		assert.Equal(t, expectedCheptel.ID, actual.Cheptels[i].ID, "Cheptel ID should be equal")
	}
}

func AssertUserCreated(t *testing.T, expected, actual entity.User, now time.Time) {
	AssertUser(t, expected, actual)
	utils.AssertCreated(t, actual.Model, now)
}

func AssertUserUpdated(t *testing.T, expected, actual entity.User, now time.Time) {
	AssertUser(t, expected, actual)
	utils.AssertUpdated(t, actual.Model, now)
}
