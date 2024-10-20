package core

import "errors"

var (
	// auth and user
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInvalidAuthConfig     = errors.New("invalid secret")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrUserNotFound          = errors.New("user not found")
	ErrInternal              = errors.New("internal error")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrRefreshTokenNotValid  = errors.New("refresh token not valid")
	ErrAlreadyDeleted        = errors.New("already deleted")

	// validation
	ErrValidationFailed = errors.New("validation failed")
)
