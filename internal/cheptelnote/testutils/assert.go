package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func AssertCheptelNote(t *testing.T, expected, actual entity.CheptelNote) {
	assert.Equal(t, expected.ID, actual.ID, "ID should be equal")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	assert.Equal(t, expected.TemperatureDay, actual.TemperatureDay, "TemperatureDay should be equal")
	assert.Equal(t, expected.TemperatureNight, actual.TemperatureNight, "TemperatureNight should be equal")
	assert.Equal(t, expected.Weather, actual.Weather, "Weather should be equal")
	assert.Equal(t, expected.Flora, actual.Flora, "Flora should be equal")
	assert.Equal(t, expected.State, actual.State, "State should be equal")
	assert.Equal(t, expected.Observation, actual.Observation, "Observation should be equal")
}

func AssertCheptelNotes(t *testing.T, expected, actual []entity.CheptelNote) {
	assert.Len(t, actual, len(expected))
	for i := range expected {
		AssertCheptelNote(t, expected[i], actual[i])
	}
}

func AssertCheptelNoteCreated(t *testing.T, expected, actual entity.CheptelNote, now time.Time) {
	AssertCheptelNote(t, expected, actual)
	utils.AssertCreated(t, actual.Model, now)
}

func AssertCheptelNoteUpdated(t *testing.T, expected, actual entity.CheptelNote, now time.Time) {
	AssertCheptelNote(t, expected, actual)
	utils.AssertUpdated(t, actual.Model, now)
}
