package core

import "context"

type (
	AuthService interface {
		Login(ctx context.Context, user User) (*string, error)
		Signup(ctx context.Context, user User) error
		UpdatePassword(ctx context.Context, user User) error
	}

	AuthConfig struct {
		Secret   string
		TokenTTL int
	}
)
