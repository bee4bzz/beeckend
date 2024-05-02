package schema

import (
	"fmt"
	"testing"

	"4d63.com/optional"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/test"
	j "gitlab.com/fogo-dev/infrastructure/web-api/pkg/json"
	val "gitlab.com/fogo-dev/infrastructure/web-api/pkg/validation"
)

var (
	validConfirmRequest = ConfirmRequest{
		OwnerUUID:  test.TokenAdminUnconfirmed.OwnerUUID,
		OwnerType:  entity.DeviceSecurityType,
		Expiration: 1,
		TokenUUID:  test.TokenAdminUnconfirmed.UUID,
		Token:      test.ValidToken,
	}

	ValidCreateRequest = CreateRequest{
		UUID:      optional.Optional[uuid.UUID]{test.TokenAdminUnconfirmed.UUID},
		OwnerUUID: test.TokenAdminUnconfirmed.OwnerUUID,
		OwnerType: entity.DeviceSecurityType,
	}
)

func TestConfirmRequest_Validate(t *testing.T) {
	t.Run("Succeed to validate the ConfirmRequestBase", func(t *testing.T) {
		err := validConfirmRequest.Validate()
		assert.NoError(t, err)
	})

	t.Run("Fail to validate the ConfirmRequestBase", func(t *testing.T) {
		nilTokenConfirmRequest := validConfirmRequest
		nilTokenConfirmRequest.Token = ""
		err := nilTokenConfirmRequest.Validate()
		assert.Error(t, err)

		nilOwnerRequestConfirmRequest := validConfirmRequest
		nilOwnerRequestConfirmRequest.TokenUUID = uuid.UUID{}
		err = nilOwnerRequestConfirmRequest.Validate()
		assert.ErrorContains(t, err, val.ErrNilUUID.Error())
	})

	t.Run("Succeed to validate the CreateRequest", func(t *testing.T) {
		err := ValidCreateRequest.Validate()
		assert.NoError(t, err)
	})

	t.Run("Fail to validate the CreateRequest", func(t *testing.T) {
		nilUUIDCreateSchema := ValidCreateRequest
		nilUUIDCreateSchema.UUID = optional.Optional[uuid.UUID]{uuid.Nil}
		err := nilUUIDCreateSchema.Validate()
		assert.Error(t, err)

		NilOwnerRequestCreateSchema := ValidCreateRequest
		NilOwnerRequestCreateSchema.OwnerUUID = uuid.UUID{}
		err = NilOwnerRequestCreateSchema.Validate()
		assert.Error(t, err)
	})
}

var (
	UUID                    = entity.GenerateUUID()
	ValidConfirmRequestJSON = fmt.Sprintf(`{"challenge":"%s","token":"%s"}`, UUID, test.ValidToken)
	ValidCreateRequestJSON  = fmt.Sprintf(`{"UUID":"%s", "owner_UUID":"%s"}`, UUID, UUID)
)

func TestJsonify(t *testing.T) {
	t.Run("succeed to jsonify", func(t *testing.T) {
		list := []j.Testable{
			j.Test[ConfirmRequest]{Name: "ConfirmRequest", JSON: ValidConfirmRequestJSON},
			j.Test[CreateRequest]{Name: "CreateRequest", JSON: ValidCreateRequestJSON},
		}

		for _, test := range list {
			test.RunTest(t)
		}
	})
}
