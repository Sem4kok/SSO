package tests

import (
	ssov1 "SSO/contract/gen/go/sso"
	"SSO/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	appID          = 1
	emptyAppId     = 0
	appSecret      = "test-secret"
	passDefaultLen = 11
)

func TestRegisterLogin_Login(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	responseRep, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Password: password,
		Email:    email,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, responseRep.GetUserId())

	responseLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	// fixing time when we have logged in
	loginTime := time.Now()
	const DeltaSeconds = 1

	jwToken := responseLogin.Token
	require.NotEmpty(t, jwToken)

	tokenParsed, err := jwt.Parse(jwToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	// check is parsedToken valid
	jwtClaims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, jwtClaims["email"].(string))
	assert.Equal(t, appID, int(jwtClaims["app_id"].(float64)))
	assert.Equal(t, responseRep.GetUserId(), int64(jwtClaims["uid"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), jwtClaims["exp"].(float64), DeltaSeconds)
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			expectedErr: "email is required",
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			password:    "",
			expectedErr: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			appID:       appID,
			expectedErr: "invalid credentials",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    gofakeit.Password(true, true, true, true, false, passDefaultLen),
			appID:       emptyAppId,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, passDefaultLen),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestRegisterLogin_DoubleRegister(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	responseRep, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Password: password,
		Email:    email,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, responseRep.GetUserId())

	responseRepSecond, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Password: password,
		Email:    email,
	})
	require.Error(t, err)
	assert.Empty(t, responseRepSecond.GetUserId())
	assert.Contains(t, err.(error).Error(), "user already exists")
}

func TestRegisterLogin_DoubleRegister_Login(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passDefaultLen)

	responseRep, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Password: password,
		Email:    email,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, responseRep.GetUserId())

	responseRepSecond, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Password: password,
		Email:    email,
	})
	require.Error(t, err)
	assert.Empty(t, responseRepSecond.GetUserId())
	assert.Contains(t, err.(error).Error(), "user already exists")

	responseLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	// fixing time when we have logged in
	loginTime := time.Now()
	const DeltaSeconds = 1

	jwToken := responseLogin.Token
	require.NotEmpty(t, jwToken)

	tokenParsed, err := jwt.Parse(jwToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	// check is parsedToken valid
	jwtClaims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, jwtClaims["email"].(string))
	assert.Equal(t, appID, int(jwtClaims["app_id"].(float64)))
	assert.Equal(t, responseRep.GetUserId(), int64(jwtClaims["uid"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), jwtClaims["exp"].(float64), DeltaSeconds)
}
