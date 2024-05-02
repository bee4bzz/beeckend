package token

var (
	ErrAlreadyConfirmed     = e.New("owner already confirmed")
	ErrAlreadyConfirmedResp = errors.ErrForbiddenResp.Join(ErrAlreadyConfirmed)

	ErrNoTokenFound     = e.New("no token found")
	ErrNoTokenFoundResp = errors.ErrNotFoundResp.Join(ErrNoTokenFound)
)
