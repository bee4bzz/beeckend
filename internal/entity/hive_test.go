package entity

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/stretchr/testify/assert"
)

var (
	EmptyHive   = Hive{}
	InvalidHive = Hive{
		Name:      utils.ValidName(),
		CheptelID: 1,
		Notes: []HiveNote{
			InvalidHiveNote,
		},
	}
	ValidHive = Hive{
		Name:      utils.ValidName(),
		CheptelID: 1,
	}

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

func TestHive(t *testing.T) {
	assert.ErrorContains(t, EmptyHive.Validate(), "CheptelID: cannot be blank; Name: cannot be blank.", "Hive should not be empty")
	assert.ErrorContains(t, InvalidHive.Validate(), "Notes:", "Hive should have a valid email")
	assert.NoError(t, ValidHive.Validate(), "Hive should be valid")
}

func TestHiveNote(t *testing.T) {
	assert.ErrorContains(t, EmptyHiveNote.Validate(), "HiveID: cannot be blank; Name: cannot be blank; Observation: cannot be blank; Operation: cannot be blank.", "HiveNote should not be empty")
	err := InvalidHiveNote.Validate()
	assert.ErrorContains(t, err, "Albums:", "HiveNote should not be empty")
	assert.ErrorContains(t, err, "NBRisers: must be no greater than 5", "HiveNote should not be empty")
	assert.NoError(t, ValidHiveNote.Validate(), "HiveNote should be valid")
}
