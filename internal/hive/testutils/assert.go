package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/stretchr/testify/assert"
)

func AssertHive(t *testing.T, expected, actual entity.Hive) {
	assert.NotEmpty(t, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	assert.Equal(t, expected.Notes, actual.Notes, "Notes should be equal")
}

func AssertHiveCreated(t *testing.T, expected, actual entity.Hive, now time.Time) {
	AssertHive(t, expected, actual)
	assert.GreaterOrEqual(t, actual.CreatedAt.Unix(), now.Unix(), "CreatedAt should be greater than now")
	assert.GreaterOrEqual(t, actual.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}

func AssertHiveUpdated(t *testing.T, expected, actual entity.Hive, now time.Time) {
	AssertHive(t, expected, actual)
	assert.LessOrEqual(t, actual.CreatedAt.Unix(), now.Unix(), "CreatedAt should be less than now")
	assert.GreaterOrEqual(t, actual.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}
