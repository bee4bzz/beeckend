package schema

import (
	"4d63.com/optional"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	optionalvalidation "gitlab.com/fogo-dev/infrastructure/web-api/pkg/optional-validation"
	val "gitlab.com/fogo-dev/infrastructure/web-api/pkg/validation"
)

type ConfirmRequest struct {
	OwnerUUID  uuid.UUID        `json:"-" swaggertype:"string"`
	OwnerType  entity.TokenType `json:"-" swaggertype:"string"`
	Expiration int64            `json:"-"`
	TokenUUID  uuid.UUID        `json:"challenge,omitempty" swaggertype:"string"`
	Token      string           `json:"token,omitempty"`
}

// Validate ConfirmRequest.
func (m ConfirmRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OwnerUUID, val.NotNilUUID),
		validation.Field(&m.OwnerType, validation.Required, validation.In(entity.TokenTypes...)),
		validation.Field(&m.TokenUUID, val.NotNilUUID),
		validation.Field(&m.Token, validation.Required),
	)
}

// CreateRequest is the body expected during tokenData creation.
type CreateRequest struct {
	UUID      optional.Optional[uuid.UUID] `json:"UUID,omitempty" swaggertype:"string" format:"UUID"`
	OwnerUUID uuid.UUID                    `json:"owner_UUID" swaggertype:"string"`
	OwnerType entity.TokenType             `json:"-" swaggertype:"string"`
}

// Validate CreateRequest.
func (m CreateRequest) Validate() error {
	err := validation.ValidateStruct(&m,
		validation.Field(&m.UUID, optionalvalidation.NotNilUUID),
		validation.Field(&m.OwnerUUID, val.NotNilUUID),
		validation.Field(&m.OwnerType, validation.Required, validation.In(entity.TokenTypes...)),
	)
	return err
}
