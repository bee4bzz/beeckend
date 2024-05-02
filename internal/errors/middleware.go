package errors

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gaetanDubuc/beeckend/internal/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler creates a middleware that handles panics and errors encountered during HTTP request processing.
func Handler(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) (err error) {
		defer func() {
			if e := recover(); e != nil {
				var ok bool
				if err, ok = e.(error); !ok {
					err = fmt.Errorf("%v", e)
				}

				logger.Errorf("recovered from panic (%v): %s", err, debug.Stack())
			}

			if err != nil {
				res := buildErrorResponse(err)
				if res.StatusCode() == http.StatusInternalServerError {
					logger.Errorf("encountered internal server error: %v", err)
				}
				c.Response.WriteHeader(res.StatusCode())
				if err = c.Write(res); err != nil {
					logger.Errorf("failed writing error response: %v", err)
				}
				c.Abort() // skip any pending handlers since an error has occurred
				err = nil // return nil because the error is already handled
			}
		}()
		return c.Next()
	}
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err error) ErrorResponse {
	switch err := err.(type) {
	case ErrorResponse:
		return err
	}

	switch err {
	case sql.ErrNoRows, gorm.ErrRecordNotFound, gorm.ErrForeignKeyViolated:
		return ErrForbiddenResp
	default:
		return ErrInternalServerResp
	}
}
