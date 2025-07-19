package auth

import (
	"errors"
	"regexp"
)

// Валидируем пароль: должен быть минимум из символов и содержать буквы с цифрами
func isValidPassword(password string) error {
	if len(password) < 5 {
		return errors.New("password must be at least 5 characters long")
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasLetter || !hasDigit {
		return errors.New("password must contain both letters and numbers")
	}

	return nil
}
