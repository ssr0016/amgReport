package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

// IsValidPassword checks if the password meets the minimum security criteria.
func IsValidPassword(password string) bool {
	// Password should be at least 8 characters long and contain at least one number and one letter.
	if len(password) < 8 {
		return false
	}

	var hasLetter, hasDigit bool
	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z':
			hasLetter = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}
