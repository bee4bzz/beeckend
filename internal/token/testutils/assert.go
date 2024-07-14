package testutils

import (
	"testing"
	"time"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/stretchr/testify/assert"
)

func AssertTokenCreated(t *testing.T, expected, actual entity.Token, now time.Time) {
	test.AssertCreated(t, actual.Model, now)
	AssertToken(t, expected, actual)
}

func AssertTokenUpdated(t *testing.T, expected, actual entity.Token, now time.Time) {
	test.AssertUpdated(t, actual.Model, now)
	AssertToken(t, expected, actual)
}

func AssertToken(t *testing.T, expected, actual entity.Token) {
	test.AssertNotEmpty(t, actual.Model)
	assert.Equal(t, expected.OwnerID, actual.OwnerID)
	if expected.Token != "" {
		assert.Equal(t, expected.Token, actual.Token)
	} else {
		assert.Empty(t, actual.Token)
	}
	if expected.HashedToken != "" {
		assert.Equal(t, expected.HashedToken, actual.HashedToken)
	}
	assert.NotEmpty(t, actual.HashedToken)
}

func AssertTokens(t *testing.T, want, actual []entity.Token) {
	if len(want) != len(actual) {
		t.Errorf("want %d tokens, actual %d", len(want), len(actual))
		return
	}

	for i, w := range want {
		g := actual[i]
		AssertToken(t, w, g)
	}
}
