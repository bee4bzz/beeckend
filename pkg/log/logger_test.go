package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduction(t *testing.T) {
	assert.NotNil(t, NewProduction())
}

func TestNewDevelopment(t *testing.T) {
	assert.NotNil(t, NewDevelopment())
}

func TestNewForTest(t *testing.T) {
	logger, entries := NewForTest()
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 1")
	assert.Equal(t, 1, entries.Len())
	logger.Info("msg 2")
	logger.Info("msg 3")
	assert.Equal(t, 3, entries.Len())
	logentry := entries.TakeAll()
	assert.Equal(t, logentry[0].Message, "msg 1")
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 4")
	assert.Equal(t, 1, entries.Len())
}
