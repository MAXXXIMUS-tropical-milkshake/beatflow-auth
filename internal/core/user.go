package core

import (
	"context"
	"time"
)

type (
	User struct {
		ID           int
		Username     string
		Email        string
		PasswordHash string
		CreatedAt    time.Time
		UpdatedAt    time.Time
		IsDeleted    bool
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

	GetUser struct {
		ID       *int
		Username *string
		Email    *string
	}

	UserService interface {
		UpdateUser(ctx context.Context, user UpdateUser) (*User, error)
		DeleteUser(ctx context.Context, userID int) error
		GetUser(ctx context.Context, user GetUser) (*User, error)
	}

	UserStore interface {
		AddUser(ctx context.Context, user User) (userID int, err error)
		GetUserByUsername(ctx context.Context, username string) (user *User, err error)
		GetUserByID(ctx context.Context, userID int) (user *User, err error)
		GetUserByEmail(ctx context.Context, email string) (user *User, err error)
		UpdateUser(cxt context.Context, user UpdateUser) (userID int, err error)
		DeleteUser(ctx context.Context, userID int) error
	}
)
