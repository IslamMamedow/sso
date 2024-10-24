package tests

import (
	"context"
	ssov1 "github.com/IslamMamedow/protos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sso/tests/suite"
	"testing"
	"time"
)

const (
	emptyAppId = 0
	appID      = 10
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := registerUser(ctx, st, email, pass)
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := loginUser(ctx, st, email, pass)
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

}

func TestRegisterLogin_DublicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := registerUser(ctx, st, email, pass)
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respReg, err = registerUser(ctx, st, email, pass)
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")

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
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "password is required",
		},
		{
			name:        "Register with empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with both empty",
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
		appId       int32
		expectedErr string
	}{
		{
			name:        "Login with empty password",
			email:       gofakeit.Email(),
			password:    "",
			appId:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with empty email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "email is required",
		},
		{
			name:        "Login with empty appId",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appId:       emptyAppId,
			expectedErr: "app_id is required",
		},

		{
			name:        "Login with wrong password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appId:       appID,
			expectedErr: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appId,
			})

			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})

	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}

func registerUser(
	ctx context.Context,
	st *suite.Suite,
	email string,
	password string,
) (*ssov1.RegisterResponse, error) {
	return st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
}

func loginUser(
	ctx context.Context,
	st *suite.Suite,
	email string,
	password string,
) (*ssov1.LoginResponse, error) {
	return st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
}
