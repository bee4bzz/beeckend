package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func AssertCheptel(t *testing.T, expected, actual entity.Cheptel) {
	assert.NotEmpty(t, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	assert.Equal(t, expected.Notes, actual.Notes, "Notes should be equal")
	assert.Equal(t, expected.Albums, actual.Albums, "Albums should be equal")
	assert.Equal(t, expected.Hives, actual.Hives, "Hives should be equal")
}

func AssertCheptelCreated(t *testing.T, expected, actual entity.Cheptel, now time.Time) {
	AssertCheptel(t, expected, actual)
	utils.AssertCreated(t, expected.Model, actual.Model, now)
}

func AssertCheptelUpdated(t *testing.T, expected, actual entity.Cheptel, now time.Time) {
	AssertCheptel(t, expected, actual)
	utils.AssertUpdated(t, expected.Model, actual.Model, now)
}
