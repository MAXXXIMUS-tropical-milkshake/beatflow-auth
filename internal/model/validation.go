package model

import (
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/validator"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
	userv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/user"
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

func ValidateUpdateUserRequest(v *validator.Validator, req *userv1.UpdateUserRequest) {
	for _, path := range req.UpdateMask.Paths {
		validateUpdateUserRequestPath(v, path)
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

func ValidateGetUserRequest(v *validator.Validator, req *userv1.GetUserRequest) {
	for _, path := range req.GetMask.Paths {
		validateGetUserRequestPath(v, path)
		if path == "username" {
			validateUsername(v, req.GetUser().GetUsername())
		} else if path == "email" {
			validateEmail(v, req.GetUser().GetEmail())
		} else if path == "user_id" {
			validateUserId(v, int(req.GetUser().GetUserId()))
		}
	}
}

func validateUserId(v *validator.Validator, userID int) {
	v.Check(validator.AtLeast(userID, 1), "user_id", "must be positive")
}

func validateEmail(v *validator.Validator, email string) {
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid")
}

func validateUsername(v *validator.Validator, username string) {
	v.Check(validator.Between(len(username), 2, 32), "username", "length must be between 2 and 32")
}

func validatePassword(v *validator.Validator, password string) {
	v.Check(validator.AtLeast(len(password), 8), "password", "must contain at least 8 characters")
}

func validateUpdateUserRequestPath(v *validator.Validator, path string) {
	v.Check(validator.OneOf(path, "username", "email", "password"), "path", "path should be one of username, email or password")
}

func validateGetUserRequestPath(v *validator.Validator, path string) {
	v.Check(validator.OneOf(path, "username", "email", "user_id"), "user_id", "path should be one of username, email or user_id")
}
