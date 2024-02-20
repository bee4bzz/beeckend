package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	emptyAlbum   = Album{}
	InvalidAlbum = Album{
		Name:      "Album",
		OwnerID:   1,
		OwnerType: "Hive",
		Paths:     []string{""},
	}
	ValidAlbum = Album{
		Name:      "Album",
		Paths:     []string{"path"},
		OwnerID:   1,
		OwnerType: "Hive",
	}
)

func TestAlbum(t *testing.T) {
	assert.ErrorContains(t, emptyAlbum.Validate(), "Name: cannot be blank; OwnerID: cannot be blank; OwnerType: cannot be blank.", "Album should not be empty")
	assert.ErrorContains(t, InvalidAlbum.Validate(), "Paths: (0: cannot be blank.).", "Album should not be empty")
	assert.NoError(t, ValidAlbum.Validate(), "Album should be valid")
}
