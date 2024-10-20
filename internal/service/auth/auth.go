package auth

import (
	"context"
	"time"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/jwt"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userStorage         core.UserStore
	refreshTokenStorage core.RefreshTokenStore
	authConfig          core.AuthConfig
}

func NewConfig(secret string, accessTokenTTL, refreshTokenTTL int) core.AuthConfig {
	return core.AuthConfig{
		Secret:          secret,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	}
}

func New(userStorage core.UserStore, refreshTokenStorage core.RefreshTokenStore, authConfig core.AuthConfig) core.AuthService {
	return &service{
		userStorage:         userStorage,
		refreshTokenStorage: refreshTokenStorage,
		authConfig:          authConfig,
	}
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (accesstoken, refreshtoken *string, err error) {
	userID, err := s.refreshTokenStorage.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, nil, err
	}

	err = s.refreshTokenStorage.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, nil, err
	}

	userFromDB, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, nil, err
	}

	accessToken, err := jwt.GenerateToken(userFromDB.ID, s.authConfig)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, nil, err
	}

	newRefreshToken := uuid.New().String()
	err = s.refreshTokenStorage.SetRefreshToken(ctx, userFromDB.ID, newRefreshToken, time.Minute*time.Duration(s.authConfig.RefreshTokenTTL))
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, nil, err
	}

	return accessToken, &newRefreshToken, nil
}

func (s *service) Login(ctx context.Context, user core.User) (accesstoken, refreshtoken *string, err error) {
	userFromDB, err := s.userStorage.GetUserByEmail(ctx, user.Email)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, nil, err
	}

	if userFromDB.IsDeleted {
		logger.Log().Error(ctx, core.ErrAlreadyDeleted.Error())
		return nil, nil, core.ErrAlreadyDeleted
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.PasswordHash), []byte(user.PasswordHash))
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, nil, core.ErrInvalidCredentials
	}

	accessToken, err := jwt.GenerateToken(userFromDB.ID, s.authConfig)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, nil, err
	}

	refreshToken := uuid.New().String()
	err = s.refreshTokenStorage.SetRefreshToken(ctx, userFromDB.ID, refreshToken, time.Minute*time.Duration(s.authConfig.RefreshTokenTTL))
	if err != nil {
		logger.Log().Info(ctx, err.Error())
		return nil, nil, err
	}

	return accessToken, &refreshToken, nil
}

func (s *service) Signup(ctx context.Context, user core.User) (*core.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	user.PasswordHash = string(passwordHash)

	userID, err := s.userStorage.AddUser(ctx, user)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	retUser, err := s.userStorage.GetUserByID(ctx, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	return retUser, nil
}
