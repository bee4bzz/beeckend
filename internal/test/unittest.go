package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ServiceTestCases[T, E any] []ServiceTestCase[T, E]

type ServiceTestCase[T, E any] struct {
	Name          string
	Req           T
	WantedError   error
	RegisterMocks func()
	Ctx           context.Context
}

func (tc ServiceTestCase[T, E]) ShouldFailAndBeEmpty(t *testing.T, Method func(context.Context, T) (E, error)) {
	if tc.RegisterMocks != nil {
		tc.RegisterMocks()
	}

	shouldBeEmpty, err := Method(tc.Ctx, tc.Req)

	if tc.WantedError != nil {
		assert.ErrorIs(t, err, tc.WantedError)
	} else {
		assert.Error(t, err)
	}
	assert.Empty(t, shouldBeEmpty)
}

func (tc ServiceTestCase[T, E]) ShouldFail(t *testing.T, Method func(context.Context, T) error) {
	if tc.RegisterMocks != nil {
		tc.RegisterMocks()
	}

	err := Method(tc.Ctx, tc.Req)

	if tc.WantedError != nil {
		assert.ErrorIs(t, err, tc.WantedError)
	} else {
		assert.Error(t, err)
	}
}

func (tc ServiceTestCase[T, E]) ShouldSucceedAndNotEmpty(t *testing.T, Method func(context.Context, T) (E, error)) E {
	if tc.RegisterMocks != nil {
		tc.RegisterMocks()
	}

	result, err := Method(tc.Ctx, tc.Req)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	return result
}

func NewContext() (*gin.Context, *httptest.ResponseRecorder) {
	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	return ctx, res
}
