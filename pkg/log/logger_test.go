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
	logger, entries, b := NewForTest()
	assert.Equal(t, 0, entries.Len())
	assert.Equal(t, 0, b.Len())
	logger.Info("msg 1")
	assert.Equal(t, 1, entries.Len())
	assert.Contains(t, b.String(), "msg 1")
	logger.Info("msg 2")
	logger.Info("msg 3")
	assert.Equal(t, 3, entries.Len())
	assert.Contains(t, b.String(), "msg 2")
	assert.Contains(t, b.String(), "msg 3")
	logentry := entries.TakeAll()
	assert.Equal(t, logentry[0].Message, "msg 1")
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 4")
	assert.Equal(t, 1, entries.Len())
}

func TestLogger_With(t *testing.T) {
	logger, entries, _ := NewForTest()
	logger2 := logger.With("key", "value")
	logger2.Info("msg")

	logentry := entries.TakeAll()
	assert.Equal(t, logentry[0].Context[0].Key, "key")
}

func TestLogger_Named(t *testing.T) {
	logger, entries, _ := NewForTest()
	logger2 := logger.Named("name")
	logger2.Info("msg")

	logentry := entries.TakeAll()
	assert.Equal(t, "name", logentry[0].LoggerName)
}
