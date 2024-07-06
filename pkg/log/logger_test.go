package log

import (
	"bytes"
	"context"
	"net/http"
	"reflect"
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
	logger2 := logger.With(context.Background(), "key", "value")
	assert.False(t, reflect.DeepEqual(logger2, logger))

	req := buildRequest("abc", "123")
	ctx := WithRequest(context.Background(), req)
	logger3 := logger2.With(ctx)
	assert.False(t, reflect.DeepEqual(logger3, logger2))

	logger3.Info("msg")

	logentry := entries.TakeAll()
	assert.Equal(t, logentry[0].Context[0].Key, "key")
	assert.Equal(t, logentry[0].Context[1].Key, "request_id")
	assert.Equal(t, logentry[0].Context[2].Key, "correlation_id")
}

func TestLogger_Named(t *testing.T) {
	logger, entries, _ := NewForTest()
	logger2 := logger.Named("name")
	logger2.Info("msg")

	logentry := entries.TakeAll()
	assert.Equal(t, "name", logentry[0].LoggerName)
}

func TestWithRequest(t *testing.T) {
	req := buildRequest("abc", "123")
	ctx := WithRequest(context.Background(), req)
	assert.Equal(t, "abc", ctx.Value(requestIDKey).(string))
	assert.Equal(t, "123", ctx.Value(correlationIDKey).(string))

	req = buildRequest("", "123")
	ctx = WithRequest(context.Background(), req)
	assert.NotEmpty(t, ctx.Value(requestIDKey).(string))
	assert.Equal(t, "123", ctx.Value(correlationIDKey).(string))
}

func Test_getCorrelationID(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", bytes.NewBufferString(""))
	assert.Empty(t, getCorrelationID(req))
	req.Header.Set("X-Correlation-ID", "test")
	assert.Equal(t, "test", getCorrelationID(req))
}

func Test_getRequestID(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", bytes.NewBufferString(""))
	assert.Empty(t, getRequestID(req))
	req.Header.Set("X-Request-ID", "test")
	assert.Equal(t, "test", getRequestID(req))
}

func buildRequest(requestID, correlationID string) *http.Request {
	req, _ := http.NewRequest("GET", "http://example.com", bytes.NewBufferString(""))
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}
	if correlationID != "" {
		req.Header.Set("X-Correlation-ID", correlationID)
	}
	return req
}
