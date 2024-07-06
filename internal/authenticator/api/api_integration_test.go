package auth

/*
func TestIntegrationAPI_Login(t *testing.T) {
	_, baseURL, _, _, _ := authtestutils.InitData(t)

	t.Run("login succeed", func(t *testing.T) {
		// list of tests to run
		tests := []struct {
			Name               string
			UUID               uuid.UUID
			Username, Password string
			BasicAuth          bool
		}{
			{"login succeed with body authentication", test.UserNoAdminConfirmed.UUID, test.UserNoAdminConfirmed.Email, test.ValidPassword, false},
			{"login succeed with basic auth", test.UserNoAdminConfirmed.UUID, test.UserNoAdminConfirmed.Email, test.ValidPassword, true},
		}

		for _, infos := range tests {
			t.Run(infos.Name, func(t *testing.T) {
				header := authtestutils.GetAuthenticatedHeader(t, baseURL, infos.Username, infos.Password, infos.BasicAuth)

				// check if the token is valid and we can access a protected path
				user := usertestutils.GetMe(t, baseURL, header)
				assert.Equal(t, user.Email, infos.Username)
			})
		}
	})

	t.Run("login fail due to bad request or service", func(t *testing.T) {
		// list of tests to run
		tests := []struct {
			Name, Username, Password string
			WantStatus               int
		}{
			{
				"login with empty password fail",
				test.UserNoAdminConfirmed.Email,
				"",
				http.StatusBadRequest,
			},
			{
				"login with invalid email fail",
				test.InvalidEmail,
				test.ValidPassword,
				http.StatusBadRequest,
			},
			{
				"login with empty email fail",
				test.InvalidEmail,
				test.ValidPassword,
				http.StatusBadRequest,
			},
			{
				"login with wrong password fail",
				test.UserNoAdminConfirmed.Email,
				test.InvalidPassword,
				http.StatusUnauthorized,
			},
			{
				"login with an unknown user",
				testutils.ValidEmailForTest(),
				test.ValidPassword,
				http.StatusUnauthorized,
			},
		}

		for _, infos := range tests {
			t.Run(infos.Name, func(t *testing.T) {
				tc := authtestutils.Login(t, baseURL, infos.Username, infos.Password, false)
				result := tc.CopyWith(test.OpAPITestCase{
					WantStatus: optional.Of(infos.WantStatus),
				}).CheckEndpointFromClient(t)

				// check if the response contains the token
				assert.NotContains(t, result, "token")
				assert.NotContains(t, result, "refresh_token")
			})
		}
	})
}

func TestIntegrationAPI_PublicKey(t *testing.T) {
	_, baseURL, _, _, _ := authtestutils.InitData(t)

	t.Run("get the public key and verify a signed jwt", func(t *testing.T) {
		// we login to get a jwt
		result := authtestutils.Login(t, baseURL, test.UserNoAdminConfirmed.Email, test.ValidPassword, true).CheckEndpointFromClient(t)

		// we get the public key
		pubKeyRes := authtestutils.PublicKeyRootTest(baseURL).CheckEndpointFromClient(t)
		pubKey, err := crypt.DecodePEMToPublicKey([]byte(pubKeyRes["key"].(string)))

		assert.NoError(t, err)

		// we extract the uuid from the jwt and verify it
		UUID := authtestutils.ExtractAuthenticationJWT(result["token"].(string), pubKey.(*rsa.PublicKey), fogojwt.RS512)

		assert.Equal(t, test.UserNoAdminConfirmed.UUID.String(), UUID)
	})
}

func TestIntegrationAPI_Logout(t *testing.T) {
	cfg, baseURL, _, _, _ := authtestutils.InitData(t)

	t.Run("logout should work properly", func(t *testing.T) {
		// we login to create a refresh token
		user, header := usertestutils.CreateAndConfirmUser(t, baseURL, cfg)

		resp := auth.Login(t, baseURL, user.Email, test.ValidPassword, true).CheckEndpointFromClient(t)
		refreshToken := resp[RefreshTokenJSONKey].(string)

		// get the public key to verify the jwt
		pubKeyRes := authtestutils.PublicKeyRootTest(baseURL).CheckEndpointFromClient(t)
		pubKey, _ := crypt.DecodePEMToPublicKey([]byte(pubKeyRes["key"].(string)))

		// extract the refresh token and the refresh token id
		_, refreshToken = refreshtestutils.ExtractRefreshToken(refreshToken, pubKey, "RS512")

		// we logout
		response := authtestutils.LogoutRootTest(baseURL, header).CheckEndpointFromClient(t)

		assert.Empty(t, response)

		time.Sleep(time.Duration(cfg.JWTExpiration+1) * time.Second)

		// we check that the refresh token is not in the database anymore
		refreshtestutils.RefreshJWTRootTest(baseURL, header, refreshToken, user.GetUUID().String()).CopyWith(test.OpAPITestCase{
			WantStatus:   optional.Of(http.StatusNotFound),
			WantResponse: optional.Of(""),
		}).CheckEndpointFromClient(t)
	})

	t.Run("logout being unauthenticated should fail", func(t *testing.T) {
		authtestutils.LogoutRootTest(baseURL, nil).CopyWith(test.OpAPITestCase{
			WantStatus:   optional.Of(http.StatusUnauthorized),
			WantResponse: optional.Of(""),
		}).CheckEndpointFromClient(t)
	})
}

func TestIntegrationAPI_RefreshJWT(t *testing.T) {
	cfg, baseURL, _, _, _ := authtestutils.InitData(t)

	// get the public key to verify the jwt
	pubKeyRes := authtestutils.PublicKeyRootTest(baseURL).CheckEndpointFromClient(t)
	pubKey, err := crypt.DecodePEMToPublicKey([]byte(pubKeyRes["key"].(string)))

	assert.NoError(t, err)

	t.Run("succeed to refresh the jwt newly created", func(t *testing.T) {
		// login a first time to get the refresh token
		results := authtestutils.Login(t, baseURL, test.UserAdminConfirmed.Email, test.ValidPassword, false).CheckEndpointFromClient(t)
		header := authtestutils.BearerAuthHeader(results[TokenJSONKey].(string))
		AssertRefreshToken(t, results)

		// extract the refresh token and the refresh token id
		refreshTokenID, refreshToken := refreshtestutils.ExtractRefreshToken(results[RefreshTokenJSONKey].(string), pubKey, "RS512")

		t.Run("refresh the token without an expired session should fail", func(t *testing.T) {
			refreshtestutils.RefreshJWTRootTest(baseURL, header, refreshToken, refreshTokenID).CopyWith(
				test.OpAPITestCase{
					WantStatus:   optional.Of(http.StatusUnauthorized),
					WantResponse: optional.Of(errors.ErrUnauthorized.Error() + "(.)*"),
				},
			).CheckEndpointFromClient(t)
		})

		// refresh the jwt with the refresh token
		time.Sleep(time.Duration(cfg.JWTExpiration+1) * time.Second)
		results = refreshtestutils.RefreshJWTRootTest(baseURL, header, refreshToken, refreshTokenID).CheckEndpointFromClient(t)
		AssertRefreshToken(t, results)

		header = authtestutils.BearerAuthHeader(results[TokenJSONKey].(string))

		// login with the new jwt to check if it's valid
		user := usertestutils.GetMe(t, baseURL, header)

		assert.Equal(t, test.UserAdminConfirmed.Email, user.Email)
	})

	t.Run("refresh the token should fail", func(t *testing.T) {
		// login a first time to get the refresh token
		results := authtestutils.Login(t, baseURL, test.UserAdminConfirmed.Email, test.ValidPassword, false).CheckEndpointFromClient(t)
		header := authtestutils.BearerAuthHeader(results[TokenJSONKey].(string))
		AssertRefreshToken(t, results)

		// extract the refresh token and the refresh token id
		refreshTokenID, refreshToken := refreshtestutils.ExtractRefreshToken(results[RefreshTokenJSONKey].(string), pubKey, "RS512")
		refreshJWTRootTest := refreshtestutils.RefreshJWTRootTest(baseURL, header, refreshToken, refreshTokenID)
		// We use the token the first time here to make the test fail the second time
		time.Sleep(time.Duration(cfg.JWTExpiration+1) * time.Second)
		refreshJWTRootTest.CheckEndpointFromClient(t)

		tests := []test.OpAPITestCase{
			{
				Name:         "refresh a second time with the same token",
				WantStatus:   optional.Of(http.StatusNotFound),
				WantResponse: optional.Of(""),
			},
			{
				Name:         "refresh with body ill formed",
				Body:         optional.Of(refreshtestutils.MakeRefreshJWTBody("ill-formed", "ill-formed")),
				WantStatus:   optional.Of(http.StatusBadRequest),
				WantResponse: optional.Of(""),
			},
			{
				Name:         "refresh with blank information in body",
				Body:         optional.Of(refreshtestutils.MakeRefreshJWTBody("", "")),
				WantStatus:   optional.Of(http.StatusBadRequest),
				WantResponse: optional.Of(""),
			},
		}

		for _, tc := range tests {
			t.Run(tc.Name, func(t *testing.T) {
				tokentestutils.AssertFailedResponseClient(t, refreshJWTRootTest.CopyWith(tc))
			})
		}
	})
}

func AssertRefreshToken(t *testing.T, results map[string]any) {
	assert.Contains(t, results, RefreshTokenJSONKey)
	assert.Contains(t, results, TokenJSONKey)
	assert.NotEmpty(t, results[RefreshTokenJSONKey])
	assert.NotEmpty(t, results[TokenJSONKey])
}
*/
