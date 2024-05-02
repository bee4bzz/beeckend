package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
)

func AssertCheptel(t *testing.T, expected, actual entity.Cheptel) {
	assert.Equal(t, expected.ID, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	assert.Len(t, actual.Notes, len(expected.Notes), "Notes should have the same length")
	for idx, v := range actual.Notes {
		assert.Equal(t, v.ID, actual.Notes[idx].ID)
	}
	assert.Len(t, actual.Albums, len(expected.Albums), "Albums should have the same length")
	for idx, v := range actual.Albums {
		assert.Equal(t, v.ID, actual.Albums[idx].ID)
	}
	assert.Len(t, actual.Hives, len(expected.Hives), "Hives should have the same length")
	for idx, v := range actual.Hives {
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
	test.AssertCreated(t, actual.Model, now)
}

func AssertCheptelUpdated(t *testing.T, expected, actual entity.Cheptel, now time.Time) {
	AssertCheptel(t, expected, actual)
	test.AssertUpdated(t, actual.Model, now)
}
