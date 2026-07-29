package utils

import (
	"errors"
	"strings"
	"unicode"
)

// IsValidEmail checks if the string looks like an emai
func IsValidEmail(email string) bool {
	if strings.ContainsAny(email, " \t\n\r") {
		return false
	}

	atIndex := strings.IndexByte(email, '@')
	if atIndex == -1 || atIndex != strings.LastIndexByte(email, '@') {
		return false
	}

	localPart := email[:atIndex]
	domainPart := email[atIndex+1:]

	if len(localPart) == 0 || len(domainPart) == 0 {
		return false
	}

	dotIndex := strings.LastIndexByte(domainPart, '.')
	if dotIndex == -1 {
		return false
	}

	if dotIndex == 0 || dotIndex == len(domainPart)-1 {
		return false
	}

	return true
}

// ValidatePasswordStrength checks if the password meets security requirements.
// It requires 8+ chars, 1 uppercase, 1 lowercase, 1 number, and 1 special character.
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

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

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
