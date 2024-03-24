package entity

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

var (
	emptyAlbum   = Album{}
	InvalidAlbum = Album{
		Name:        "Album",
		OwnerID:     1,
		OwnerType:   "Hive",
		Observation: utils.String(""),
	}
	ValidAlbum = Album{
		Name:      "Album",
		OwnerID:   1,
		OwnerType: "Hive",
	}
)

func TestAlbum(t *testing.T) {
	assert.ErrorContains(t, emptyAlbum.Validate(), "Name: cannot be blank; OwnerID: cannot be blank; OwnerType: cannot be blank.", "Album should not be empty")
	assert.ErrorContains(t, InvalidAlbum.Validate(), "Observation: cannot be blank.", "Album should not be empty")
	assert.NoError(t, ValidAlbum.Validate(), "Album should be valid")
}

func TestPhoto(t *testing.T) {
	emptyPhoto := Photo{}
	InvalidPhoto := Photo{
		Path: "path",
	}
	ValidPhoto := Photo{
		Path:    "path",
		AlbumID: 1,
	}

	assert.ErrorContains(t, emptyPhoto.Validate(), "AlbumID: cannot be blank; Path: cannot be blank.", "Photo should not be empty")
	assert.ErrorContains(t, InvalidPhoto.Validate(), "AlbumID: cannot be blank.", "Photo should not be empty")
	assert.NoError(t, ValidPhoto.Validate(), "Photo should be valid")
}
