package user

import (
	"context"
	"errors"
	"testing"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	mocks "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestDeleteUser_Success(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)

	// service
	userService := New(userStore)

	// vars
	userID := 1

	// mock behaviour
	userStore.EXPECT().DeleteUser(mock.Anything, userID).Return(nil).Once()

	err := userService.DeleteUser(context.Background(), userID)
	assert.NoError(t, err)
}

func TestDeleteUser_Fail(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)

	// service
	userService := New(userStore)

	// vars
	wantErr := errors.New("internal error")

	// mock behaviour
	userStore.EXPECT().DeleteUser(mock.Anything, mock.Anything).Return(wantErr).Once()

	err := userService.DeleteUser(context.Background(), 1)
	assert.ErrorIs(t, err, wantErr)
}

func TestUpdateUser_Success(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)

	// service
	userService := New(userStore)

	// vars
	userID := 1
	updateUser := core.UpdateUser{
		ID:       userID,
		Username: &[]string{"alex123"}[0],
	}
	updatedUser := &core.User{
		ID:       userID,
		Username: "alex123",
	}

	// mock behaviour
	userStore.EXPECT().GetUserByID(mock.Anything, userID).Return(&core.User{}, nil).Once()
	userStore.EXPECT().UpdateUser(mock.Anything, updateUser).Return(userID, nil).Once()
	userStore.EXPECT().GetUserByID(mock.Anything, userID).Return(updatedUser, nil).Once()

	user, err := userService.UpdateUser(context.Background(), updateUser)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser, user)
}

func TestUpdateUser_UpdatePasswordSuccess(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)

	// service
	userService := New(userStore)

	// vars
	userID := 1
	oldPassword := "Qwerty123456"
	newPassword := "12345678"
	hashedOldPassword, err := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	require.NoError(t, err)
	updateUser := core.UpdateUser{
		ID: userID,
		Password: &core.UpdatePassword{
			OldPassword: oldPassword,
			NewPassword: newPassword,
		},
	}
	userFromDB := &core.User{
		PasswordHash: string(hashedOldPassword),
	}

	// mock behaviour
	userStore.EXPECT().GetUserByID(mock.Anything, userID).Return(userFromDB, nil).Once()
	userStore.EXPECT().UpdateUser(mock.Anything, mock.MatchedBy(func(user core.UpdateUser) bool {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password.NewPassword), []byte(newPassword))
		return err == nil
	})).Return(userID, nil).Once()
	userStore.EXPECT().GetUserByID(mock.Anything, userID).Return(&core.User{}, nil).Once()

	_, err = userService.UpdateUser(context.Background(), updateUser)
	assert.NoError(t, err)
}

func TestUpdateUser_AlreadyDeleted(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)

	// service
	userService := New(userStore)

	// vars
	userFromDB := &core.User{
		IsDeleted: true,
	}

	// mock behaviour
	userStore.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(userFromDB, nil).Once()

	_, err := userService.UpdateUser(context.Background(), core.UpdateUser{})
	assert.ErrorIs(t, err, core.ErrAlreadyDeleted)
}

func TestUpdateUser_Fail(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)

	// service
	userService := New(userStore)

	// vars
	ctx := context.Background()
	wantErr := errors.New("internal error")

	// tests
	tests := []struct {
		name      string
		behaviour func()
	}{
		{
			name: "first GetUserByID error",
			behaviour: func() {
				userStore.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(nil, wantErr).Once()
			},
		},
		{
			name: "UpdateUser error",
			behaviour: func() {
				userStore.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(&core.User{}, nil).Once()
				userStore.EXPECT().UpdateUser(mock.Anything, mock.Anything).Return(0, wantErr).Once()
			},
		},
		{
			name: "second GetUserByID error",
			behaviour: func() {
				userStore.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(&core.User{}, nil).Once()
				userStore.EXPECT().UpdateUser(mock.Anything, mock.Anything).Return(1, nil).Once()
				userStore.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(nil, wantErr).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.behaviour()

			_, err := userService.UpdateUser(ctx, core.UpdateUser{})
			assert.ErrorIs(t, err, wantErr)
		})
	}
}

func TestGetUser_Success(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)

	// service
	userService := New(userStore)

	// vars
	userID := 1
	getUser := core.GetUser{
		ID: &userID,
	}
	userFromDB := &core.User{
		ID:       userID,
		Username: "alex123",
		Email:    "alex@gmail.com",
	}

	// mock behaviour
	userStore.EXPECT().GetUserByID(mock.Anything, userID).Return(userFromDB, nil).Once()

	user, err := userService.GetUser(context.Background(), getUser)
	assert.NoError(t, err)
	assert.Equal(t, userFromDB, user)
}

func TestGetUser_Fail(t *testing.T) {
	t.Parallel()

	// store
	userStore := mocks.NewMockUserStore(t)

	// service
	userService := New(userStore)

	// vars
	getUser := core.GetUser{
		ID: &[]int{1}[0],
	}
	wantErr := errors.New("internal error")

	// mock behaviour
	userStore.EXPECT().GetUserByID(mock.Anything, mock.Anything).Return(nil, wantErr).Once()

	_, err := userService.GetUser(context.Background(), getUser)
	assert.ErrorIs(t, err, wantErr)
}
