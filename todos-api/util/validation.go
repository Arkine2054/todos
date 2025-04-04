package util

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateEmail(email string) error {
	return validation.Validate(&email, validation.Required, is.Email)
}

func ValidatePassword(password string) error {
	return validation.Validate(&password, validation.Required, validation.Length(8, 50))
}
