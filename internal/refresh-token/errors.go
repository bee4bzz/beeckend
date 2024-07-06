package refreshtoken

import (
	e "errors"

	"gitlab.com/fogo-dev/infrastructure/web-api/internal/errors"
)

var (
	ErrRefreshJWT     = e.New("The jwt token can't be refreshed")
	ErrRefreshJWTResp = errors.ErrInternalServerResp.Join(ErrRefreshJWT)
)
