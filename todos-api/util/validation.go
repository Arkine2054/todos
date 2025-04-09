package util

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func Validate(email string, password string) error {
	err := validation.Validate(email, validation.Required, is.Email)
	if err != nil {
		return err
	}
	err = validation.Validate(password, validation.Required, validation.Length(8, 50))
	if err != nil {
		return err
	}

	return nil
}
