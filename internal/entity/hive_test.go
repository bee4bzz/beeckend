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
)

func TestHive(t *testing.T) {
	assert.ErrorContains(t, EmptyHive.Validate(), "CheptelID: cannot be blank; Name: cannot be blank.", "Hive should not be empty")
	assert.ErrorContains(t, InvalidHive.Validate(), "Notes:", "Hive should have a valid email")
	assert.NoError(t, ValidHive.Validate(), "Hive should be valid")
}
