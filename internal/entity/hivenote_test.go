package entity

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

var (
	EmptyHiveNote = HiveNote{
		Observation: utils.String(""),
	}
	InvalidHiveNote = HiveNote{
		NBRisers: MaxNBRisers + 1,
		Albums: []Album{
			InvalidAlbum,
		},
	}
	ValidHiveNote = HiveNote{
		HiveID:    1,
		Name:      utils.ValidName(),
		Operation: utils.ValidName(),
	}
)

func TestHiveNote(t *testing.T) {
	assert.ErrorContains(t, EmptyHiveNote.Validate(), "HiveID: cannot be blank; Name: cannot be blank; Observation: cannot be blank; Operation: cannot be blank.", "HiveNote should not be empty")
	err := InvalidHiveNote.Validate()
	assert.ErrorContains(t, err, "Albums:", "HiveNote should not be empty")
	assert.ErrorContains(t, err, "NBRisers: must be no greater than 5", "HiveNote should not be empty")
	assert.NoError(t, ValidHiveNote.Validate(), "HiveNote should be valid")
}
