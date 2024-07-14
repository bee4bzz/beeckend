package refreshtoken

import (
	e "errors"

	"github.com/gaetanDubuc/beeckend/internal/errors"
)

var (
	ErrRefreshJWT     = e.New("The jwt token can't be refreshed")
	ErrRefreshJWTResp = errors.ErrInternalServerResp.Join(ErrRefreshJWT)
)
