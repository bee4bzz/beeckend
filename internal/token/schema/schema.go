package schema

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	validation "github.com/go-ozzo/ozzo-validation"
)

type ConfirmRequest struct {
	OwnerID    uint             `json:"-" swaggertype:"string"`
	OwnerType  entity.TokenType `json:"-" swaggertype:"string"`
	Expiration int              `json:"-"`
	TokenID    uint             `json:"challenge,omitempty" swaggertype:"string"`
	Token      string           `json:"token,omitempty"`
}

// Validate ConfirmRequest.
func (m ConfirmRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OwnerID, validation.Required),
		validation.Field(&m.OwnerType, validation.Required, validation.In(entity.TokenTypes...)),
		validation.Field(&m.TokenID, validation.Required),
		validation.Field(&m.Token, validation.Required),
	)
}

// CreateRequest is the body expected during tokenData creation.
type CreateRequest struct {
	OwnerID   uint             `json:"owner_ID" swaggertype:"string"`
	OwnerType entity.TokenType `json:"-" swaggertype:"string"`
}

// Validate CreateRequest.
func (m CreateRequest) Validate() error {
	err := validation.ValidateStruct(&m,
		validation.Field(&m.OwnerID, validation.Required),
		validation.Field(&m.OwnerType, validation.Required, validation.In(entity.TokenTypes...)),
	)
	return err
}
