package request

import (
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/validator"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
)

func ValidateLoginRequest(v *validator.Validator, req *authv1.LoginRequest) {
	validateEmail(v, req.GetEmail())
	validatePassword(v, req.GetPassword())
}

func ValidateSignupRequest(v *validator.Validator, req *authv1.SignupRequest) {
	validateUsername(v, req.GetUsername())
	validateEmail(v, req.GetEmail())
	validatePassword(v, req.GetPassword())
}

func ValidateUpdateUserRequest(v *validator.Validator, req *authv1.UpdateUserRequest) {
	for _, path := range req.UpdateMask.Paths {
		validatePath(v, path)
		if path == "username" {
			validateUsername(v, req.User.GetUsername())
		} else if path == "email" {
			validateEmail(v, req.User.GetEmail())
		} else if path == "password" {
			validatePassword(v, req.User.GetPassword().OldPassword)
			validatePassword(v, req.User.GetPassword().NewPassword)
		}
	}
}
