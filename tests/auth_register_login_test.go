package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/igortoigildin/goph-keeper/pkg/auth_v1"
	"github.com/igortoigildin/goph-keeper/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	login := gofakeit.Bird()
	pass := randomFakePassword()

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login: login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login: login,
		Password: pass,
	})
	require.NoError(t, err)
	
	token := respLogin.GetRefreshToken()
	require.NotEmpty(t, token)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}