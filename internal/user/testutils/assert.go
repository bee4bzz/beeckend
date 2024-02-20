package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/stretchr/testify/assert"
)

func AssertUser(t *testing.T, expected, actual entity.User) {
	assert.NotEmpty(t, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	assert.Equal(t, expected.Email, actual.Email, "Email should be equal")
	assert.Equal(t, expected.Cheptels, actual.Cheptels, "Cheptels should be equal")
}

func AssertUserCreated(t *testing.T, expected, actual entity.User, now time.Time) {
	AssertUser(t, expected, actual)
	assert.GreaterOrEqual(t, actual.CreatedAt.Unix(), now.Unix(), "CreatedAt should be greater than now")
	assert.GreaterOrEqual(t, actual.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}

func AssertUserUpdated(t *testing.T, expected, actual entity.User, now time.Time) {
	AssertUser(t, expected, actual)
	assert.LessOrEqual(t, actual.CreatedAt.Unix(), now.Unix(), "CreatedAt should be less than now")
	assert.GreaterOrEqual(t, actual.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}
