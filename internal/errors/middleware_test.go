package errors

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gaetanDubuc/beeckend/pkg/log"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

type MockDataWriterFailure struct {
	routing.DataWriter
}

func (m MockDataWriterFailure) Write(http.ResponseWriter, interface{}) error {
	return fmt.Errorf("test")
}

func TestHandler(t *testing.T) {
	t.Run("normal processing", func(t *testing.T) {
		logger, entries, _ := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerOK)
		assert.Nil(t, ctx.Next())
		assert.Zero(t, entries.Len())
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("error writing response", func(t *testing.T) {
		logger, entries, _ := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerOK)
		ctx.SetDataWriter(MockDataWriterFailure{routing.DefaultDataWriter})
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 2, entries.Len())
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("error processing", func(t *testing.T) {
		logger, entries, _ := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerError)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 1, entries.Len())
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("HTTP error processing", func(t *testing.T) {
		logger, entries, _ := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerHTTPError)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 0, entries.Len())
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("panic processing", func(t *testing.T) {
		logger, entries, _ := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerPanic)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 2, entries.Len())
		assert.Equal(t, http.StatusInternalServerError, res.Code)
	})

	t.Run("sql now row processing", func(t *testing.T) {
		logger, entries, _ := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerSQLNoRows)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 0, entries.Len())
		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("Duplicate key", func(t *testing.T) {
		logger, entries, _ := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, handlerErrDuplicate)
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 0, entries.Len())
		assert.Equal(t, http.StatusConflict, res.Code)
	})
	t.Run("not member error", func(t *testing.T) {
		logger, entries, _ := log.NewForTest()
		handler := Handler(logger)
		ctx, res := buildContext(handler, func(ctx *routing.Context) error {
			return u2herrors.ErrNotMember
		})
		assert.Nil(t, ctx.Next())
		assert.Equal(t, 0, entries.Len())
		assert.Equal(t, http.StatusForbidden, res.Code)
	})
}

func TestBuildErrorResponse(t *testing.T) {
	t.Run("HTTPErrorNotFound", func(t *testing.T) {
		res := buildErrorResponse(routing.NewHTTPError(http.StatusNotFound))
		assert.Equal(t, http.StatusNotFound, res.Status)
	})

	t.Run("HTTPErrorUnauthorized", func(t *testing.T) {
		res := buildErrorResponse(routing.NewHTTPError(http.StatusUnauthorized))
		assert.Equal(t, http.StatusUnauthorized, res.Status)
	})

	t.Run("HTTPErrorBadRequest", func(t *testing.T) {
		res := buildErrorResponse(routing.NewHTTPError(http.StatusBadRequest))
		assert.Equal(t, http.StatusBadRequest, res.Status)
	})

	t.Run("HTTPErrorBadRequest", func(t *testing.T) {
		res := buildErrorResponse(routing.NewHTTPError(http.StatusConflict))
		assert.Equal(t, http.StatusConflict, res.Status)
	})

	t.Run("validation.Errors", func(t *testing.T) {
		res := buildErrorResponse(validation.Errors{})
		assert.Equal(t, http.StatusBadRequest, res.Status)
	})

	t.Run("HTTPErrorForbidden", func(t *testing.T) {
		res := buildErrorResponse(routing.NewHTTPError(http.StatusForbidden))
		assert.Equal(t, http.StatusForbidden, res.Status)
	})

	t.Run("sql.ErrNoRows", func(t *testing.T) {
		res := buildErrorResponse(sql.ErrNoRows)
		assert.Equal(t, http.StatusForbidden, res.Status)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		res := buildErrorResponse(fmt.Errorf("test"))
		assert.Equal(t, http.StatusInternalServerError, res.Status)
	})

	t.Run("middleware.CombinatorError", func(t *testing.T) {
		res := buildErrorResponse(middleware.CombinatorError{
			CombinatorTypeErrorMsg: "test",
			Errors: []error{
				ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "A",
				},
				ErrorResponse{
					Status:  http.StatusForbidden,
					Message: "B",
				},
			},
		})
		assert.Equal(t, http.StatusBadRequest, res.Status)
	})
}

func buildContext(handlers ...routing.Handler) (*routing.Context, *httptest.ResponseRecorder) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://127.0.0.1/users", nil)
	return routing.NewContext(res, req, handlers...), res
}

func handlerOK(c *routing.Context) error {
	return c.Write("test")
}

func handlerError(c *routing.Context) error {
	return fmt.Errorf("abc")
}

func handlerHTTPError(c *routing.Context) error {
	return ErrNotFoundResp
}

func handlerPanic(c *routing.Context) error {
	panic("xyz")
}

func handlerSQLNoRows(c *routing.Context) error {
	return sql.ErrNoRows
}

func handlerErrDuplicate(c *routing.Context) error {
	return ErrDuplicate
}
