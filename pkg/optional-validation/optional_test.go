package optionalvalidation

import (
	"testing"

	"4d63.com/optional"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {
	optionalstringRule := By[string](is.UUID)
	err := optionalstringRule.Validate(optional.Of(""))
	assert.NoError(t, err)

	optionalUUIDRule := By[uuid.UUID](is.UUID)
	err = optionalUUIDRule.Validate(optional.Of(uuid.Nil))
	assert.NoError(t, err)

	err = BoolIsPresent.Validate(optional.Of(false))
	assert.NoError(t, err)

	err = BoolIsPresent.Validate(optional.Optional[bool](nil))
	assert.Error(t, err)
}
