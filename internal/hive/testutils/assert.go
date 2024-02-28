package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
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
	utils.AssertCreated(t, actual.Model, now)
}

func AssertHiveUpdated(t *testing.T, expected, actual entity.Hive, now time.Time) {
	AssertHive(t, expected, actual)
	utils.AssertUpdated(t, actual.Model, now)
}
