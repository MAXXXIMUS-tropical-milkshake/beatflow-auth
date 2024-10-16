package core

import (
	"context"
	"time"
)

type (
	AuthService interface {
		Login(ctx context.Context, user User) (accessToken *string, refreshToken *string, err error)
		Signup(ctx context.Context, user User) (*User, error)
		UpdateUser(ctx context.Context, user UpdateUser) (*User, error)
		RefreshToken(ctx context.Context, refreshToken string) (*string, *string, error)
	}

	AuthConfig struct {
		Secret          string
		AccessTokenTTL  int
		RefreshTokenTTL int
	}

	RefreshTokenStore interface {
		SetRefreshToken(ctx context.Context, userID int, tokenID string, expiresIn time.Duration) error
		GetRefreshToken(ctx context.Context, tokenID string) (int, error)
		DeleteRefreshToken(ctx context.Context, prevTokenID string) error
	}
)
