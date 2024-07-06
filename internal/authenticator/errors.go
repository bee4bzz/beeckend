package auth

import (
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/errors"

	e "errors"
)

var (
	// internal error or user not found.
	ErrGetUser            = e.New("the user could not be found or an internal error occurred")
	ErrGetUserResp        = errors.ErrUnauthorizedResp.Join(ErrGetUser)
	ErrWrongPassword      = e.New("attempt to connect, wrong password")
	ErrWrongPasswordResp  = errors.ErrUnauthorizedResp.Join(ErrWrongPassword)
	ErrUnverifiedUser     = e.New("unverified user")
	ErrUnverifiedUserResp = errors.ErrForbiddenResp.Join(ErrUnverifiedUser)
	ErrUserNotFound       = e.New("User not found in context.")
	ErrUserNotFoundResp   = errors.ErrNotFoundResp.Join(ErrUserNotFound)

	ErrInvalidSession     = e.New("Invalid session.")
	ErrInvalidSessionResp = errors.ErrForbiddenResp.Join(ErrInvalidSession)

	ErrUserNotConfirmed = e.New("User is not confirmed yet, access is denied!")
	// ErrUserNotConfirmedResp is the error that returns in case user is not
	// confirmed.
	ErrUserNotConfirmedResp = errors.ErrForbiddenResp.Join(ErrUserNotConfirmed)

	ErrUserAlreadyConfirmed     = e.New("User is already confirmed, access is denied!")
	ErrUserAlreadyConfirmedResp = errors.ErrForbiddenResp.Join(ErrUserAlreadyConfirmed)

	ErrOnlySameUser     = e.New("You can't access the resource of another user.")
	ErrOnlySameUserResp = errors.ErrForbiddenResp.Join(ErrOnlySameUser)

	ErrOnlyExpiredSession     = e.New("You can only access the resource with an expired session.")
	ErrOnlyExpiredSessionResp = errors.ErrForbiddenResp.Join(ErrOnlyExpiredSession)
)
