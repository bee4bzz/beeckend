package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustMarshal(t *testing.T) {
	t.Run("should panic when json.Marshal fails", func(t *testing.T) {
		assert.Panics(t, func() {
			MustMarshal(make(chan int))
		})
	})

	t.Run("should return the marshaled value", func(t *testing.T) {
		expected := []byte(`{"key":"value"}`)
		actual := MustMarshal(map[string]string{"key": "value"})

		if string(actual) != string(expected) {
			t.Errorf("MustMarshal returned %s, expected %s", actual, expected)
		}
	})
}
