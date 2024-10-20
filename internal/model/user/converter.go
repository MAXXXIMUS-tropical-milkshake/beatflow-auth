package user

import (
	"strings"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	userv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/user"
)

func FromUpdateUserRequest(req *userv1.UpdateUserRequest, userID int) *core.UpdateUser {
	user := new(core.UpdateUser)
	user.ID = userID

	for _, path := range req.UpdateMask.Paths {
		if path == "username" {
			user.Username = &req.GetUser().Username
		} else if path == "email" {
			user.Email = &req.GetUser().Email
		} else if strings.HasPrefix(path, "password") {
			user.Password = new(core.UpdatePassword)
			user.Password.OldPassword = req.GetUser().GetPassword().GetOldPassword()
			user.Password.NewPassword = req.GetUser().GetPassword().GetNewPassword()
		}
	}

	return user
}

func FromGetUserRequest(req *userv1.GetUserRequest) *core.User {
	user := new(core.User)
	user.ID = int(req.GetUserId())

	return user
}

func ToGetUserResponse(user core.User) *userv1.GetUserResponse {
	return &userv1.GetUserResponse{
		UserId:    int64(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		IsDeleted: user.IsDeleted,
	}
}

func ToUpdateUserResponse(user core.User) *userv1.UpdateUserResponse {
	return &userv1.UpdateUserResponse{
		UserId:   int64(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}
}
