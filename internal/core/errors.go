package core

import "errors"

var (
	// auth
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidAuthConfig  = errors.New("invalid secret")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)
