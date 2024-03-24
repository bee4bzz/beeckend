package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func AssertAlbum(t *testing.T, expected, actual entity.Album) {
	assert.Equal(t, expected.ID, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt, "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt, "UpdatedAt should not be empty")
	assert.Equal(t, expected.Name, actual.Name, "Name should be equal")
	assert.Equal(t, expected.Observation, actual.Observation, "Observation should be equal")
	assert.Equal(t, expected.OwnerID, actual.OwnerID, "OwnerID should be equal")
	assert.Equal(t, expected.OwnerType, actual.OwnerType, "OwnerType should be equal")
	for idx, v := range expected.Photos {
		assert.Equal(t, v, actual.Photos[idx], "Photos should be equal")
	}
}

func AssertAlbums(t *testing.T, expected, actual []entity.Album) {
	assert.Len(t, actual, len(expected))
	for i := range expected {
		AssertAlbum(t, expected[i], actual[i])
	}
}

func AssertAlbumCreated(t *testing.T, expected, actual entity.Album, now time.Time) {
	AssertAlbum(t, expected, actual)
	utils.AssertCreated(t, actual.Model, now)
}

func AssertAlbumUpdated(t *testing.T, expected, actual entity.Album, now time.Time) {
	AssertAlbum(t, expected, actual)
	utils.AssertUpdated(t, actual.Model, now)
}
