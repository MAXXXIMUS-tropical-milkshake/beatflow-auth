package core

import "context"

type (
	User struct {
		ID           int
		Username     string
		PasswordHash string
	}

	UserStore interface {
		AddUser(ctx context.Context, user User) (userID int, err error)
		GetUserByUsername(ctx context.Context, username string) (user *User, err error)
		UpdateUser(ctx context.Context, user User) (userID int, err error)
	}
)
