package token

import (
	e "errors"

	"github.com/gaetanDubuc/beeckend/internal/errors"
)

var (
	ErrGetToken     = e.New("An error occurred while getting the token!")
	ErrGetTokenResp = errors.ErrInternalServerResp.Join(ErrGetToken)

	ErrGetOwner     = e.New("An error occurred while getting the owner!")
	ErrGetOwnerResp = errors.ErrNotFoundResp.Join(ErrGetOwner)

	// errors during SendToken service.
	ErrSendToken     = e.New("An error occurred while sending the token!")
	ErrSendTokenResp = errors.ErrInternalServerResp.Join(ErrSendToken)

	// errors during ConfirmToken service.
	ErrUpdateOwner      = e.New("An error occurred while updating the owner information!")
	ErrUpdateOwnerResp  = errors.ErrInternalServerResp.Join(ErrUpdateOwner)
	ErrDeleteToken      = e.New("An error occurred while deleting the owner!")
	ErrDeleteTokenResp  = errors.ErrInternalServerResp.Join(ErrDeleteToken)
	ErrTokenInvalid     = e.New("The token is invalid!")
	ErrTokenInvalidResp = errors.ErrBadRequestResp.Join(ErrTokenInvalid)
	ErrTokenExpired     = e.New("The token is expired!")
	ErrTokenExpiredResp = errors.ErrForbiddenResp.Join(ErrTokenExpired)

	ErrCoolDown     = e.New("wait until the cooldown is over")
	ErrCoolDownResp = errors.ErrForbiddenResp.Join(ErrCoolDown)

	// errors during Creation service.
	ErrTokenCreation     = e.New("An error occurred while creating the token's informations!")
	ErrTokenCreationResp = errors.ErrInternalServerResp.Join(ErrTokenCreation)

	// errors during count service.
	ErrCountToken     = e.New("An error occurred while counting the tokens!")
	ErrCountTokenResp = errors.ErrInternalServerResp.Join(ErrCountToken)

	// errors during query service.
	ErrQueryToken     = e.New("An error occurred while querying the tokens!")
	ErrQueryTokenResp = errors.ErrInternalServerResp.Join(ErrQueryToken)

	ErrAlreadyConfirmed     = e.New("owner already confirmed")
	ErrAlreadyConfirmedResp = errors.ErrForbiddenResp.Join(ErrAlreadyConfirmed)

	ErrNoTokenFound     = e.New("no token found")
	ErrNoTokenFoundResp = errors.ErrNotFoundResp.Join(ErrNoTokenFound)
)
