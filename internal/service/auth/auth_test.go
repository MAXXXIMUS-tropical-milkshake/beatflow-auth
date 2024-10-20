package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	mocks "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const (
	password = "Qwerty123456"
)

func isUUID(val string) bool {
	_, err := uuid.Parse(val)
	return err == nil
}

func parseToken(tokenString, secret string) (userID *int, expiresAt *time.Time, err error) {
	errParseToken := errors.New("failed to parse token")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errParseToken
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, nil, errParseToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(float64)
		if !ok {
			return nil, nil, errParseToken
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			return nil, nil, errParseToken
		}

		idInt := int(id)
		expTime := time.Unix(int64(exp), 0)

		return &idInt, &expTime, nil
	}

	return nil, nil, errParseToken
}

func TestSignup_Success(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)
	refreshTokenStore := mocks.NewMockRefreshTokenStore(t)

	// config
	authConfig := core.AuthConfig{
		Secret:          "secret",
		AccessTokenTTL:  10,
		RefreshTokenTTL: 20,
	}

	// service
	authService := New(userStore, refreshTokenStore, authConfig)

	// vars
	userID := 1
	user := core.User{
		Username:     "alex123",
		Email:        "alex@gmail.com",
		PasswordHash: password,
	}
	userFromDB := &user
	userFromDB.ID = userID

	// mock behaviour
	userStore.EXPECT().AddUser(mock.Anything, mock.MatchedBy(func(user core.User) bool {
		err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		return err == nil
	})).Return(userID, nil).Once()
	userStore.EXPECT().GetUserByID(mock.Anything, userID).Return(userFromDB, nil).Once()

	retUser, err := authService.Signup(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, userFromDB, retUser)
}

func TestSignup_Fail(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)
	refreshTokenStore := mocks.NewMockRefreshTokenStore(t)

	// config
	authConfig := core.AuthConfig{
		Secret:          "secret",
		AccessTokenTTL:  10,
		RefreshTokenTTL: 20,
	}

	// service
	authService := New(userStore, refreshTokenStore, authConfig)

	// vars
	ctx := context.Background()
	user := core.User{
		PasswordHash: password,
	}
	wantErr := errors.New("internal error")

	tests := []struct {
		name      string
		behaviour func()
	}{
		{
			name: "AddUser error",
			behaviour: func() {
				userStore.EXPECT().AddUser(mock.Anything, mock.Anything).Return(0, wantErr).Once()
			},
		},
		{
			name: "GetUserByID error",
			behaviour: func() {
				userStore.EXPECT().AddUser(mock.Anything, mock.Anything).Return(0, nil).Once()
				userStore.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(nil, wantErr).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behaviour()

			_, err := authService.Signup(ctx, user)
			assert.ErrorIs(t, err, wantErr)
		})
	}
}

func TestLogin_Success(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)
	refreshTokenStore := mocks.NewMockRefreshTokenStore(t)

	// config
	atTTL := 10
	rtTTL := 20

	authConfig := core.AuthConfig{
		Secret:          "secret",
		AccessTokenTTL:  atTTL,
		RefreshTokenTTL: rtTTL,
	}

	// service
	authService := New(userStore, refreshTokenStore, authConfig)

	// vars
	userID := 1
	email := "alex@gmail.com"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	user := core.User{
		Email:        email,
		PasswordHash: password,
	}
	userFromDB := &core.User{
		ID:           userID,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	var retRefreshToken string
	wantTime := time.Now().Add(time.Duration(atTTL) * time.Minute)
	delta := 5 * time.Second

	// mock behaviour
	userStore.EXPECT().GetUserByEmail(mock.Anything, email).Return(userFromDB, nil).Once()
	refreshTokenStore.EXPECT().SetRefreshToken(mock.Anything, userID, mock.MatchedBy(func(tokenID string) bool {
		retRefreshToken = tokenID
		return isUUID(tokenID)
	}), time.Minute*time.Duration(rtTTL)).Return(nil).Once()

	at, rt, err := authService.Login(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, retRefreshToken, *rt)

	id, exp, err := parseToken(*at, authConfig.Secret)
	require.NoError(t, err)
	assert.Equal(t, userID, *id)
	assert.True(t, wantTime.Sub(*exp) <= delta)
}

func TestLogin_AlreadyDeletedUser(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)
	refreshTokenStore := mocks.NewMockRefreshTokenStore(t)

	// config
	atTTL := 10
	rtTTL := 20

	authConfig := core.AuthConfig{
		Secret:          "secret",
		AccessTokenTTL:  atTTL,
		RefreshTokenTTL: rtTTL,
	}

	// service
	authService := New(userStore, refreshTokenStore, authConfig)

	// vars
	userFromDB := &core.User{
		IsDeleted: true,
	}

	// mock behaviour
	userStore.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(userFromDB, nil).Once()

	_, _, err := authService.Login(context.Background(), core.User{})
	assert.ErrorIs(t, err, core.ErrAlreadyDeleted)
}

func TestLogin_InvalidPassword(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)
	refreshTokenStore := mocks.NewMockRefreshTokenStore(t)

	// config
	atTTL := 10
	rtTTL := 20

	authConfig := core.AuthConfig{
		Secret:          "secret",
		AccessTokenTTL:  atTTL,
		RefreshTokenTTL: rtTTL,
	}

	// service
	authService := New(userStore, refreshTokenStore, authConfig)

	// vars
	user := core.User{
		PasswordHash: password,
	}
	userFromDB := &core.User{
		PasswordHash: "12345678",
	}

	// mock behaviour
	userStore.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(userFromDB, nil).Once()

	_, _, err := authService.Login(context.Background(), user)
	assert.ErrorIs(t, err, core.ErrInvalidCredentials)
}

func TestLogin_Fail(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)
	refreshTokenStore := mocks.NewMockRefreshTokenStore(t)

	// config
	atTTL := 10
	rtTTL := 20

	authConfig := core.AuthConfig{
		Secret:          "secret",
		AccessTokenTTL:  atTTL,
		RefreshTokenTTL: rtTTL,
	}

	// service
	authService := New(userStore, refreshTokenStore, authConfig)

	// vars
	ctx := context.Background()
	wantError := errors.New("internal error")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	user := core.User{
		PasswordHash: password,
	}
	userFromDB := &core.User{
		PasswordHash: string(passwordHash),
	}

	tests := []struct {
		name      string
		behaviour func()
	}{
		{
			name: "GetUserByEmail error",
			behaviour: func() {
				userStore.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(nil, wantError).Once()
			},
		},
		{
			name: "SetRefreshToken error",
			behaviour: func() {
				userStore.EXPECT().GetUserByEmail(mock.Anything, mock.Anything).Return(userFromDB, nil).Once()
				refreshTokenStore.EXPECT().SetRefreshToken(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wantError).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behaviour()

			_, _, err := authService.Login(ctx, user)
			assert.ErrorIs(t, err, wantError)
		})
	}
}
