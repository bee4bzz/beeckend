package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse_Error(t *testing.T) {
	e := ErrorResponse{
		Message: "abc",
	}
	assert.Equal(t, "abc", e.Error())
}

func TestErrorResponse_StatusCode(t *testing.T) {
	e := ErrorResponse{
		Status: 400,
	}
	assert.Equal(t, 400, e.StatusCode())
}

func TestInternalServerError(t *testing.T) {
	res := ErrInternalServerResp
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	assert.ErrorIs(t, res, ErrInternalServer)
	err := errors.New("test")
	res = ErrInternalServerResp.Join(err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode())
	assert.ErrorIs(t, res, err)
}

func TestNotFound(t *testing.T) {
	res := ErrNotFoundResp
	assert.Equal(t, http.StatusNotFound, res.StatusCode())
	assert.ErrorIs(t, res, ErrNotFound)
	err := errors.New("test")
	res = ErrNotFoundResp.Join(err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode())
	assert.ErrorIs(t, res, err)
}

func TestUnauthorized(t *testing.T) {
	res := ErrUnauthorizedResp
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.ErrorIs(t, res, ErrUnauthorized)
	err := errors.New("test")
	res = ErrUnauthorizedResp.Join(err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode())
	assert.ErrorIs(t, res, err)
}

func TestForbidden(t *testing.T) {
	res := ErrForbiddenResp
	assert.Equal(t, http.StatusForbidden, res.StatusCode())
	assert.ErrorIs(t, res, ErrForbidden)
	err := errors.New("test")
	res = ErrForbiddenResp.Join(err)
	assert.Equal(t, http.StatusForbidden, res.StatusCode())
	assert.ErrorIs(t, res, err)
}

func TestBadRequest(t *testing.T) {
	res := ErrBadRequestResp
	assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.ErrorIs(t, res, ErrBadRequest)
	err := errors.New("test")
	res = ErrBadRequestResp.Join(err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	assert.ErrorIs(t, res, err)
}

func TestInvalidInput(t *testing.T) {
	err1 := fmt.Errorf("1")
	err2 := fmt.Errorf("2")
	err := InvalidInput(validation.Errors{
		"xyz": err2,
		"abc": err1,
	})
	assert.Equal(t, http.StatusBadRequest, err.Status)
	assert.ErrorIs(t, err, err1)
	assert.ErrorIs(t, err, err2)
	assert.Contains(t, err.Message, "abc: 1")
	assert.Contains(t, err.Message, "xyz: 2")
}
