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
