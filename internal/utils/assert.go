package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func AssertCreated(t *testing.T, expected, actual gorm.Model, now time.Time) {
	assert.GreaterOrEqual(t, actual.CreatedAt.Unix(), now.Unix(), "CreatedAt should be greater than now")
	assert.GreaterOrEqual(t, actual.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}

func AssertUpdated(t *testing.T, expected, actual gorm.Model, now time.Time) {
	assert.LessOrEqual(t, actual.CreatedAt.Unix(), now.Unix(), "CreatedAt should be less than now")
	assert.GreaterOrEqual(t, actual.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}
