package request

import (
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
)

func FromLoginRequest(req *authv1.LoginRequest) *core.User {
	return &core.User{
		Email:        req.GetEmail(),
		PasswordHash: req.GetPassword(),
	}
}

func FromSignupRequest(req *authv1.SignupRequest) *core.User {
	return &core.User{
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: req.GetPassword(),
	}
}

func FromUpdateUserRequest(req *authv1.UpdateUserRequest, userID int) *core.UpdateUser {
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
