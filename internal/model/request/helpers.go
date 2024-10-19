package request

import (
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/model/validator"
)

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
	v.Check(validator.OneOf(path, "username", "email", "password"), "path", "path should be one of username, email or password")
}
