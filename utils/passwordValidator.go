package utils

import (
	"errors"
	"fmt"
	"unicode"
)

const PASSWORD_MIN_LENGTH = 8
const PASSWORD_MAX_LENGTH = 24

// PasswordValidator checks for a password validity.
//
// Returns error if:
//   - the password is too short (less than 8 characters)
//   - the password is too long (more than 24 characters)
//   - the password contains no upper case character
//   - the password contains no lower case character
//   - the password contains no number
//   - the password contains no special character
func PasswordValidator(password string) error {
	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	if !passwordLengthIsValid(password) {
		return fmt.Errorf("password length should be between %d and %d", PASSWORD_MIN_LENGTH, PASSWORD_MAX_LENGTH)
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		return errors.New("password should contain at least one character of each category: uppercase, lowercas, number, specialcharacter")
	}

	return nil
}

func passwordLengthIsValid(password string) bool {
	if len(password) < PASSWORD_MIN_LENGTH {
		return false
	}
	if len(password) > PASSWORD_MAX_LENGTH {
		return false
	}
	return true
}
