package auth

import (
	"context"
	"testing"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	mocks "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core/mocks"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSignup_Success(t *testing.T) {
	// init server and client
	authService := mocks.NewMockAuthService(t)
	client := &server{authService: authService}

	// vars
	username := "alex123"
	email := "alex@gmail.com"
	password := "Qwerty123456"
	user := &authv1.SignupRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	coreUser := core.User{
		Username:     username,
		Email:        email,
		PasswordHash: password,
	}
	retUser := &core.User{
		ID:       1,
		Username: username,
		Email:    email,
	}

	// mock behaviour
	authService.EXPECT().Signup(mock.Anything, coreUser).Return(retUser, nil)

	res, err := client.Signup(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, int64(retUser.ID), res.UserId)
	assert.Equal(t, retUser.Email, res.Email)
	assert.Equal(t, retUser.Username, res.Username)
}

func TestSignup_ValidationErrors(t *testing.T) {
	// init server and client
	authService := mocks.NewMockAuthService(t)
	client := &server{authService: authService}

	// vars
	password := "Qwerty123456"
	username := "alex123"
	email := "alex@gmail.com"
	wantErr := status.Error(codes.InvalidArgument, core.ErrValidationFailed.Error())

	tests := []struct {
		name string
		user *authv1.SignupRequest
	}{
		{
			name: "empty username",
			user: &authv1.SignupRequest{
				Username: "",
				Email:    email,
				Password: password,
			},
		},
		{
			name: "less than 8 chars password",
			user: &authv1.SignupRequest{
				Username: username,
				Email:    email,
				Password: "123456",
			},
		},
		{
			name: "invalid email",
			user: &authv1.SignupRequest{
				Username: username,
				Email:    "alex@",
				Password: password,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Signup(context.Background(), tt.user)
			require.Equal(t, err.Error(), wantErr.Error())
		})
	}
}

func TestSignup_SignupError(t *testing.T) {
	// init server and client
	authService := mocks.NewMockAuthService(t)
	client := &server{authService: authService}

	// vars
	user := &authv1.SignupRequest{
		Username: "alex123",
		Email:    "alex@gmail.com",
		Password: "Qwerty123456",
	}

	tests := []struct {
		name      string
		behaviour func()
	}{
		{
			name: "email already exists",
			behaviour: func() {
				authService.EXPECT().Signup(mock.Anything, mock.Anything).Return(nil, core.ErrEmailAlreadyExists).Once()
			},
		},
		{
			name: "username already exists",
			behaviour: func() {
				authService.EXPECT().Signup(mock.Anything, mock.Anything).Return(nil, core.ErrUsernameAlreadyExists).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behaviour()

			_, err := client.Signup(context.Background(), user)
			st, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, st.Code(), codes.AlreadyExists)
		})
	}
}

func TestLogin_Success(t *testing.T) {
	// init server and client
	authService := mocks.NewMockAuthService(t)
	client := &server{authService: authService}

	// vars
	email := "alex@gmail.com"
	password := "Qwerty123456"
	user := &authv1.LoginRequest{
		Email:    email,
		Password: password,
	}
	coreUser := core.User{
		Email:        email,
		PasswordHash: password,
	}
	accessToken, refreshToken := "access_token", "refresh_token"

	// mock behaviour
	authService.EXPECT().Login(mock.Anything, coreUser).Return(&accessToken, &refreshToken, nil)

	res, err := client.Login(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, accessToken, res.AccessToken)
	assert.Equal(t, refreshToken, res.RefreshToken)
}

func TestLogin_LoginError(t *testing.T) {
	// init server and client
	authService := mocks.NewMockAuthService(t)
	client := &server{authService: authService}

	// vars
	user := &authv1.LoginRequest{
		Email:    "alex@gmail.com",
		Password: "Qwerty123456",
	}

	tests := []struct {
		name      string
		behaviour func()
	}{
		{
			name: "invalid credentials",
			behaviour: func() {
				authService.EXPECT().Login(mock.Anything, mock.Anything).Return(nil, nil, core.ErrInvalidCredentials).Once()
			},
		},
		{
			name: "user not found",
			behaviour: func() {
				authService.EXPECT().Login(mock.Anything, mock.Anything).Return(nil, nil, core.ErrUserNotFound).Once()
			},
		},
		{
			name: "user already deleted",
			behaviour: func() {
				authService.EXPECT().Login(mock.Anything, mock.Anything).Return(nil, nil, core.ErrAlreadyDeleted).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behaviour()

			_, err := client.Login(context.Background(), user)
			st, ok := status.FromError(err)
			require.True(t, ok)
			require.Equal(t, st.Code(), codes.Unauthenticated)
		})
	}
}
