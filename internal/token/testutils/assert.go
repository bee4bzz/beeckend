package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	entitytestutils "gitlab.com/fogo-dev/infrastructure/web-api/internal/entity/testutils"
)

func AssertTokenCreated(t *testing.T, expected, actual entity.Token, now time.Time) {
	entitytestutils.AssertsBaseCreated(t, actual.Base, now)
	AssertToken(t, expected, actual)
}

func AssertTokenUpdated(t *testing.T, expected, actual entity.Token, now time.Time) {
	entitytestutils.AssertsBaseUpdated(t, actual.Base, now)
	AssertToken(t, expected, actual)
}

func AssertToken(t *testing.T, want, got entity.Token) {
	entitytestutils.AssertsBase(t, got.Base)
	assert.Equal(t, want.OwnerUUID, got.OwnerUUID)
	if want.Token != "" {
		assert.Equal(t, want.Token, got.Token)
	} else {
		assert.Empty(t, got.Token)
	}
	if want.HashedToken != "" {
		assert.Equal(t, want.HashedToken, got.HashedToken)
	}
	assert.NotEmpty(t, got.HashedToken)
}

func AssertTokens(t *testing.T, want, got []entity.Token) {
	if len(want) != len(got) {
		t.Errorf("want %d tokens, got %d", len(want), len(got))
		return
	}

	for i, w := range want {
		g := got[i]
		AssertToken(t, w, g)
	}
}
