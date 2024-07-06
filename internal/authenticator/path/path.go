package path

import "gitlab.com/fogo-dev/infrastructure/web-api/internal/config"

const (
	RefreshSessionPath = config.RefreshTokenGroup
	LoginPath          = "/login"
	PublicKeyPath      = "/public-key"
	LogoutPath         = "/logout"
)
