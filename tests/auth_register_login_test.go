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

	login := gofakeit.Email()
	pass := randomFakePassword()

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
	require.NotEmpty(t, token)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	login := gofakeit.Email()
	pass := randomFakePassword()

	resReg, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, resReg.GetUserId())

	resReg, err = st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.Error(t, err)
	assert.Empty(t, resReg.GetUserId())
	assert.ErrorContains(t, err, "already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		login       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			login:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty Login",
			login:       "",
			password:    randomFakePassword(),
			expectedErr: "login is required",
		},
		{
			name:        "Register with Both Empty",
			login:       "",
			password:    "",
			expectedErr: "login is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
				Login:    tt.login,
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
		login       string
		password    string
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			login:       gofakeit.BeerName(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Empty login",
			login:       "",
			password:    randomFakePassword(),
			expectedErr: "login is required",
		},
		{
			name:        "Login with both empty login && password",
			login:       "",
			password:    "",
			expectedErr: "login is required",
		},
		{
			name:        "Login with Non-Matching Password",
			login:       gofakeit.BookAuthor(),
			password:    randomFakePassword(),
			expectedErr: "failed to login",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &auth_v1.RegisterRequest{
				Login:    gofakeit.Email(),
				Password: randomFakePassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &auth_v1.LoginRequest{
				Login:    tt.login,
				Password: tt.password,
			})
			require.Error(t, err)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
