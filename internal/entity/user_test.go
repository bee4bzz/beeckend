package entity

import (
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/utils"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
)

var (
	EmptyUser   = User{}
	InvalidUser = User{
		Name:  utils.ValidName(),
		Email: random.String(10),
		Cheptels: []Cheptel{
			InvalidCheptel,
		},
	}
	ValidUser = User{
		Name:  utils.ValidName(),
		Email: utils.ValidEmail(),
		Cheptels: []Cheptel{
			ValidCheptel,
		},
	}
)

func TestUser(t *testing.T) {
	assert.ErrorContains(t, EmptyUser.Validate(), "Email: cannot be blank; Name: cannot be blank.", "User should not be empty")
	err := InvalidUser.Validate()
	assert.ErrorContains(t, err, "Cheptels", "User should not be empty")
	assert.ErrorContains(t, err, "Email: must be a valid email address.", "User should not be empty")
	assert.NoError(t, ValidUser.Validate(), "User should be valid")
}
