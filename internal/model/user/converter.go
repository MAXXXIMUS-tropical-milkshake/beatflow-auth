package user

import (
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
		} else if path == "password" {
			user.Password = new(core.UpdatePassword)
			user.Password.OldPassword = req.GetUser().Password.OldPassword
			user.Password.NewPassword = req.GetUser().Password.NewPassword
		}
	}

	return user
}

func FromGetUserRequest(req *userv1.GetUserRequest) *core.GetUser {
	user := new(core.GetUser)

	for _, path := range req.GetMask.Paths {
		if path == "username" {
			user.Username = &req.GetUser().Username
		} else if path == "user_id" {
			userIDInt := int(req.GetUser().UserId)
			user.ID = &userIDInt
		} else if path == "email" {
			user.Email = &req.GetUser().Email
		}
	}

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
