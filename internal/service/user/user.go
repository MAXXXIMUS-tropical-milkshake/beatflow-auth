package user

import (
	"context"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/logger"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userStorage core.UserStore
}

func New(userStorage core.UserStore) core.UserService {
	return &service{userStorage: userStorage}
}

func (s *service) DeleteUser(ctx context.Context, userID int) error {
	err := s.userStorage.DeleteUser(ctx, userID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *service) UpdateUser(ctx context.Context, user core.UpdateUser) (*core.User, error) {
	userFromDB, err := s.userStorage.GetUserByID(ctx, user.ID)
	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	if userFromDB.IsDeleted {
		logger.Log().Debug(ctx, core.ErrAlreadyDeleted.Error())
		return nil, core.ErrAlreadyDeleted
	}

	if user.Password != nil {
		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, err
		}

		user.Password.NewPassword = string(newPasswordHash)

		if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.PasswordHash), []byte(user.Password.OldPassword)); err != nil {
			logger.Log().Error(ctx, err.Error())
			return nil, core.ErrInvalidCredentials
		}
	}

	userID, err := s.userStorage.UpdateUser(ctx, user)
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

func (s *service) GetUser(ctx context.Context, user core.GetUser) (*core.User, error) {
	var retUser *core.User
	var err error

	if user.ID != nil {
		retUser, err = s.userStorage.GetUserByID(ctx, *user.ID)
	} else if user.Email != nil {
		retUser, err = s.userStorage.GetUserByEmail(ctx, *user.Email)
	} else if user.Username != nil {
		retUser, err = s.userStorage.GetUserByUsername(ctx, *user.Username)
	}

	if err != nil {
		logger.Log().Error(ctx, err.Error())
		return nil, err
	}

	return retUser, nil
}
