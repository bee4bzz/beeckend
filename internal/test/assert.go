package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func AssertNotEmpty(t *testing.T, actual gorm.Model) {
	assert.NotEmpty(t, actual.ID, "ID should not be empty")
	assert.NotEmpty(t, actual.CreatedAt.Unix(), "CreatedAt should not be empty")
	assert.NotEmpty(t, actual.UpdatedAt.Unix(), "UpdatedAt should not be empty")
}

func AssertCreated(t *testing.T, actual gorm.Model, now time.Time) {
	assert.NotEmpty(t, actual.ID, "ID should not be empty")
	assert.GreaterOrEqual(t, actual.CreatedAt.Unix(), now.Unix(), "CreatedAt should be greater than now")
	assert.GreaterOrEqual(t, actual.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}

func AssertUpdated(t *testing.T, actual gorm.Model, now time.Time) {
	assert.NotEmpty(t, actual.ID, "ID should not be empty")
	assert.LessOrEqual(t, actual.CreatedAt.Unix(), now.Unix(), "CreatedAt should be less than now")
	assert.GreaterOrEqual(t, actual.UpdatedAt.Unix(), now.Unix(), "UpdatedAt should be greater than now")
}
