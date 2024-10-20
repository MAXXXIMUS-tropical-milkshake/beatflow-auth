package model

import (
	"strings"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/validator"
	authv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/auth"
	userv1 "github.com/MAXXXIMUS-tropical-milkshake/beatflow-protos/gen/go/user"
	fieldmask "google.golang.org/protobuf/types/known/fieldmaskpb"
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
	validateMask(v, req.GetUpdateMask())
	for _, path := range req.GetUpdateMask().GetPaths() {
		validatePath(v, path)
		if path == "username" {
			validateUsername(v, req.User.GetUsername())
		} else if path == "email" {
			validateEmail(v, req.User.GetEmail())
		} else if strings.HasPrefix(path, "password") {
			validatePassword(v, req.User.GetPassword().GetOldPassword())
			validatePassword(v, req.User.GetPassword().GetNewPassword())
		}
	}
}

func ValidateGetUserRequest(v *validator.Validator, req *userv1.GetUserRequest) {
	validateID(v, int(req.GetUserId()))
}

func validateID(v *validator.Validator, id int) {
	v.Check(id > 0, "id", "must be positive")
}

func validateMask(v *validator.Validator, mask *fieldmask.FieldMask) {
	v.Check(mask != nil, "mask", "mask is required")
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

func validatePath(v *validator.Validator, path string) {
	v.Check(validator.OneOf(path, "username", "email") || validator.HasPrefix(path, "password"), "path", "path should be one of username, email or password")
}
