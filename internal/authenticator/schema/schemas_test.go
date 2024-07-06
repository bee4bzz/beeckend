package schema

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/entity"
	"gitlab.com/fogo-dev/infrastructure/web-api/internal/test"
	j "gitlab.com/fogo-dev/infrastructure/web-api/pkg/json"
)

func TestLoginRequest_Validate(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		req := LoginRequest{
			Username: test.ValidEmail,
			Password: test.ValidPassword,
		}
		assert.NoError(t, req.Validate())
	})
	t.Run("fail", func(t *testing.T) {
		testCases := []struct {
			LoginRequest
			errors []string
		}{
			{
				LoginRequest{
					Username: test.InvalidEmail,
				},
				[]string{
					"username: must be a valid email address",
					"password: cannot be blank",
				},
			},
			{
				LoginRequest{},
				[]string{
					"username: cannot be blank",
				},
			},
		}

		for _, tc := range testCases {
			for _, v := range tc.errors {
				assert.ErrorContains(t, tc.Validate(), v)
			}
		}
	})
}

var (
	UUID              = entity.GenerateUUID()
	FullLoginRequest  = fmt.Sprintf(`{"username":"%s","password":"%s"}`, test.ValidEmail, test.ValidPassword)
	FullLoginResponse = fmt.Sprintf(`{"%s":"%s","%s":"%s"}`, TokenJSONKey, test.ValidJWT, RefreshTokenJSONKey, test.ValidJWT)
)

func TestLoginRequest_JSON(t *testing.T) {
	list := []j.Testable{
		j.Test[LoginRequest]{Name: "LoginRequest", JSON: FullLoginRequest},
	}

	for _, test := range list {
		test.RunTest(t)
	}
}
func TestLoginResponse_JSON(t *testing.T) {
	list := []j.Testable{
		j.Test[Session]{Name: "Session", JSON: FullLoginResponse},
	}

	for _, test := range list {
		test.RunTest(t)
	}
}
