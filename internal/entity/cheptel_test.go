package entity

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/labstack/gommon/random"
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

var (
	EmptyNote = CheptelNote{
		Weather:     weather(""),
		State:       utils.String(""),
		Observation: utils.String(""),
	}
	InvalidNote = CheptelNote{
		CheptelID: 1,
		Name:      utils.ValidName(),
		Flora:     utils.ValidName(),
		Weather:   weather(random.String(10)),
	}
	ValidFilledNote = CheptelNote{
		CheptelID:        1,
		Name:             utils.ValidName(),
		TemperatureDay:   utils.Float64(10),
		TemperatureNight: utils.Float64(10),
		Flora:            utils.ValidName(),
		Weather:          CLOUDY,
		State:            utils.String(random.String(10)),
		Observation:      utils.String(random.String(10)),
	}
	ValidNote = CheptelNote{
		CheptelID: 1,
		Name:      utils.ValidName(),
		Flora:     utils.ValidName(),
		Weather:   UNKNOWN,
	}
)

func TestChecptelNote(t *testing.T) {
	assert.ErrorContains(t, EmptyNote.Validate(), "CheptelID: cannot be blank; Flora: cannot be blank; Name: cannot be blank; Observation: cannot be blank; State: cannot be blank; Weather: cannot be blank.", "Cheptel should not be empty")
	assert.ErrorContains(t, InvalidNote.Validate(), "Weather: must be a valid value.", "Cheptel should not be empty")
	assert.NoError(t, ValidNote.Validate(), "Cheptel should be valid")
	assert.NoError(t, ValidFilledNote.Validate(), "Cheptel should be valid")
}

func TestWeather(t *testing.T) {
	v, err := CLOUDY.Value()
	assert.NoError(t, err, "Weather should be valid")
	assert.Equal(t, "CLOUDY", v, "Weather should be valid")

	var w weather
	err = w.Scan("CLOUDY")
	assert.NoError(t, err, "Weather should be valid")
	assert.Equal(t, CLOUDY, w, "Weather should be valid")

	err = w.Scan(5)
	assert.ErrorIs(t, err, ErrInvalidData, "Weather should be invalid")
}
