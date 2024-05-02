// Package errors provides error handling utilities.
// these errors are used to return a consistent error response to the client and can be
// a base join for your custom error. see ../token/errors.go file for an example.
package errors

import (
	"errors"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	ErrNilResp = ErrorResponse{}

	ErrNotImplemented = errors.New("this feature is not yet implemented")
	// ErrNotImplementedResp is the error that is returned in
	// case an api or service is not implemented yet.
	ErrNotImplementedResp = ErrorResponse{
		Status:  http.StatusNotImplemented,
		Message: ErrNotImplemented.Error(),
		Errors:  ErrNotImplemented,
	}

	ErrBadRequest = errors.New("your request is in a bad format")
	// ErrBadRequestResp is the error that is returned in case of bad format request.
	ErrBadRequestResp = ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: ErrBadRequest.Error(),
		Errors:  ErrBadRequest,
	}

	ErrRequestIllFormedUUID = errors.New("uuid is ill formed and/or couldn't be parsed")
	// ErrRequestIllFormedUUIDResp is the error that returns in case the uuid in
	// request couldn't be found.
	ErrRequestIllFormedUUIDResp = ErrBadRequestResp.Join(ErrRequestIllFormedUUID)

	ErrRequestNotFoundUUID = errors.New("uuid can't be found in request")
	// ErrRequestNotFoundUUIDResp is the error that returns in case the uuid in
	// request couldn't be found.
	ErrRequestNotFoundUUIDResp = ErrBadRequestResp.Join(ErrRequestNotFoundUUID)

	ErrRequestIllFormedBody = errors.New("data is ill formed and/or couldn't be parsed")
	// ErrRequestIllFormedBodyResp is the error that returns in case we couldn't
	// parsed the body.
	ErrRequestIllFormedBodyResp = ErrBadRequestResp.Join(ErrRequestIllFormedBody)

	ErrRequestValidation = errors.New("data is ill formed")
	// ErrRequestValidationResp is the error that returns in case we couldn't
	// validate the request data in service.
	ErrRequestValidationResp = ErrBadRequestResp.Join(ErrRequestValidation)

	ErrUnauthorized = errors.New("you are not authenticated to perform the requested action")
	// ErrUnauthorizedResp is the error that returns in case of
	// tentative to access a ressource that user does not own and is not an administrator.
	ErrUnauthorizedResp = ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: ErrUnauthorized.Error(),
		Errors:  ErrUnauthorized,
	}

	ErrForbidden = errors.New("you are not authorized to perform the requested action")
	// ErrForbiddenResp is the error that returns in case of tentative to access
	// a ressource that user does not own and is not an administrator.
	ErrForbiddenResp = ErrorResponse{
		Status:  http.StatusForbidden,
		Message: ErrForbidden.Error(),
		Errors:  ErrForbidden,
	}

	ErrConflictResp = ErrorResponse{
		Status:  http.StatusConflict,
		Message: "conflict",
		Errors:  errors.New("conflict"),
	}

	ErrNotAllowedToAccessThisResource = errors.New("you are not allowed to access this resource")
	// ErrNotAllowedToAccessThisResourceResp is the error that returns in case of
	// tentative to access a ressource that user does not own and is not an
	// administrator.
	ErrNotAllowedToAccessThisResourceResp = ErrForbiddenResp.Join(ErrNotAllowedToAccessThisResource)

	ErrNotFound = errors.New("resource not found when it should be")
	// ErrNotFoundResp is the error that returns in case we couldn't find the
	// data.
	ErrNotFoundResp = ErrorResponse{
		Status:  http.StatusNotFound,
		Message: ErrNotFound.Error(),
		Errors:  ErrNotFound,
	}

	ErrInternalServer = errors.New("we encountered an error while processing your request")
	// ErrInternalServerResp is the error that returns in case something wrong happened server side.
	ErrInternalServerResp = ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: ErrInternalServer.Error(),
		Errors:  ErrInternalServer,
	}
	// ErrUnexpectedFailureResp is the error that returns in case something wrong
	// happened server side.
	ErrUnexpectedFailureResp = ErrInternalServerResp.Join(errors.New("unexpected failure has occurred"))
)

const RegexpErrorResponse = "(.)*status(.)*message(.)*"

// ErrorResponse is the response that represents an error.
// TODO: use wrappedErrors instead of Message.
type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Error message"`
	Errors  error  `json:"-"`
}

// Error is required by the error interface.
func (e ErrorResponse) Error() string {
	return e.Message
}

func (e ErrorResponse) Unwrap() error {
	return e.Errors
}

func (e ErrorResponse) Join(errs ...error) ErrorResponse {
	var newErrs []error
	newErrs = append(newErrs, errs...)
	newErrs = append(newErrs, e.Errors)
	e.Errors = errors.Join(newErrs...)
	e.Message = e.Errors.Error()
	return e
}

// StatusCode is required by routing.HTTPError interface.
func (e *ErrorResponse) StatusCode() int {
	return e.Status
}

// InvalidInput creates a new error response representing a data validation error (HTTP 400).
func InvalidInput(errs validation.Errors) ErrorResponse {
	var wrappedErrors []error
	for idx := range errs {
		wrappedErrors = append(wrappedErrors, errs[idx])
	}
	err := ErrRequestValidationResp.Join(wrappedErrors...)
	err.Message = strings.Join([]string{errs.Error(), ErrRequestValidationResp.Message}, "\n")
	return err
}
