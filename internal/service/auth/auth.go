package auth

import (
	"context"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/jwt"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userStorage core.UserStore
	authConfig  core.AuthConfig
}

func NewConfig(secret string, tokenTTL int) core.AuthConfig {
	return core.AuthConfig{
		Secret:   secret,
		TokenTTL: tokenTTL,
	}
}

func New(userStorage core.UserStore, authConfig core.AuthConfig) core.AuthService {
	return &service{
		userStorage: userStorage,
		authConfig:  authConfig,
	}
}

func (s *service) Login(ctx context.Context, user core.User) (*string, error) {
	userFromDB, err := s.userStorage.GetUserByUsername(ctx, user.Username)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, core.ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(userFromDB.ID, s.authConfig)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	return token, nil
}

func (s *service) Signup(ctx context.Context, user core.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	user.PasswordHash = string(passwordHash)

	_, err = s.userStorage.AddUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, user core.User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	user.PasswordHash = string(passwordHash)

	_, err = s.userStorage.UpdateUser(ctx, user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}
