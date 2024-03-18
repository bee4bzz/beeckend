package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func AssertHiveNote(t *testing.T, expected, actual entity.HiveNote) {
	assert.Equal(t, expected.ID, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	assert.Equal(t, expected.NBRisers, actual.NBRisers, "NBRisers should be equal")
	assert.Equal(t, expected.Observation, actual.Observation, "Observation should be equal")
	assert.Equal(t, expected.HiveID, actual.HiveID, "HiveID should be equal")
	assert.Equal(t, expected.Operation, actual.Operation, "Operation should be equal")
	for i := range expected.Albums {
		assert.Equal(t, expected.Albums[i].ID, actual.Albums[i].ID, "Albums.ID should be equal")
	}
}

func AssertHiveNotes(t *testing.T, expected, actual []entity.HiveNote) {
	assert.Len(t, actual, len(expected))
	for i := range expected {
		AssertHiveNote(t, expected[i], actual[i])
	}
}

func AssertHiveNoteCreated(t *testing.T, expected, actual entity.HiveNote, now time.Time) {
	AssertHiveNote(t, expected, actual)
	utils.AssertCreated(t, actual.Model, now)
}

func AssertHiveNoteUpdated(t *testing.T, expected, actual entity.HiveNote, now time.Time) {
	AssertHiveNote(t, expected, actual)
	utils.AssertUpdated(t, actual.Model, now)
}
