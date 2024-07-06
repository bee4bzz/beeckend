package schema

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateRequest struct {
	UserID uint `json:"-" swaggertype:"string"`
}

func (m CreateRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required),
	)
}

type RefreshRequest struct {
	UserID  uint   `json:"-" swaggertype:"string"`
	TokenID uint   `json:"challenge,omitempty" swaggertype:"string" form:"challenge"`
	Token   string `json:"token,omitempty" form:"token"`
}

func (m RefreshRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required),
		validation.Field(&m.TokenID, validation.Required),
		validation.Field(&m.Token, validation.Required),
	)
}

type DeleteFromUserRequest struct {
	UserID uint `json:"-" swaggertype:"string"`
}

func (m DeleteFromUserRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UserID, validation.Required),
	)
}
