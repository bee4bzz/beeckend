package test

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/utils"
)

var (
	ValidUser = entity.User{
		Name:  utils.ValidName(),
		Email: utils.ValidEmail(),
		Cheptels: []entity.Cheptel{
			ValidCheptel,
		},
	}

	ValidCheptel = entity.Cheptel{
		Name: utils.ValidName(),
	}
)
