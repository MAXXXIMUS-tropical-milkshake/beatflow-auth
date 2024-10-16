package core

import "context"

type (
	User struct {
		ID           int
		Username     string
		Email        string
		PasswordHash string
	}

	UpdateUser struct {
		ID       int
		Username *string
		Email    *string
		Password *UpdatePassword
	}

	UpdatePassword struct {
		OldPassword string
		NewPassword string
	}

	UserStore interface {
		AddUser(ctx context.Context, user User) (userID int, err error)
		GetUserByUsername(ctx context.Context, username string) (user *User, err error)
		GetUserByID(ctx context.Context, userID int) (user *User, err error)
		GetUserByEmail(ctx context.Context, email string) (user *User, err error)
		UpdateUser(cxt context.Context, user UpdateUser) (userID int, err error)
	}
)
