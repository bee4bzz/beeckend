package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func AssertCheptel(t *testing.T, expected, actual entity.Cheptel) {
	assert.Equal(t, expected.ID, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	for idx, v := range expected.Notes {
		assert.Equal(t, v.ID, actual.Notes[idx].ID)
	}
	for idx, v := range expected.Albums {
		assert.Equal(t, v.ID, actual.Albums[idx].ID)
	}
	for idx, v := range expected.Hives {
		assert.Equal(t, v.ID, actual.Hives[idx].ID)
	}
}

func AssertCheptels(t *testing.T, expected, actual []entity.Cheptel) {
	assert.Len(t, actual, len(expected))
	for i := range expected {
		AssertCheptel(t, expected[i], actual[i])
	}
}

func AssertCheptelCreated(t *testing.T, expected, actual entity.Cheptel, now time.Time) {
	AssertCheptel(t, expected, actual)
	utils.AssertCreated(t, actual.Model, now)
}

func AssertCheptelUpdated(t *testing.T, expected, actual entity.Cheptel, now time.Time) {
	AssertCheptel(t, expected, actual)
	utils.AssertUpdated(t, actual.Model, now)
}
