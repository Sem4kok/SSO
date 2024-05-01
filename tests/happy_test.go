package tests

import (
	ssov1 "SSO/contract/gen/go/sso"
	"SSO/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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

	jwToken := responseLogin.Token
	require.NotEmpty(t, jwToken)

	tokenParsed, err := jwt.Parse(jwToken, parseSecret)
}
