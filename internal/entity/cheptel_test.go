package entity

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

var (
	EmptyCheptel   = Cheptel{}
	InvalidCheptel = Cheptel{
		Hives: []Hive{
			InvalidHive,
		},
		Notes: []CheptelNote{
			InvalidNote,
		},
		Albums: []Album{
			InvalidAlbum,
		},
	}
	ValidCheptel = Cheptel{
		Name: utils.ValidName(),
	}
)

func TestCheptel(t *testing.T) {
	assert.ErrorContains(t, EmptyCheptel.Validate(), "Name: cannot be blank.", "Cheptel should not be empty")
	assert.ErrorContains(t, InvalidCheptel.Validate(), "Name: cannot be blank; Notes:", "Cheptel should not be empty")
	assert.NoError(t, ValidCheptel.Validate(), "Cheptel should be valid")
}
