package validator

import (
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Between(value, lower, upper int) bool {
	return lower <= value && lower <= upper
}

func AtLeast(value, lower int) bool {
	return lower <= value
}

func OneOf(value string, other ...string) bool {
	for i := range other {
		if value == other[i] {
			return true
		}
	}
	return false
}
