package password

import (
	"crypto/rand"
	"math/big"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers = "0123456789"
	symbols = "!@#$%^&*()-_=+[]{}|;:,.<>?/"
)

// GeneratePassword generates a secure password of the specified length.
func GeneratePassword(length int, includeNumbers bool, includeSymbols bool) (string, error) {
	allChars := letters
	if includeNumbers {
		allChars += numbers
	}
	if includeSymbols {
		allChars += symbols
	}

	password := make([]byte, length)

	for i := range password {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[index.Int64()]
	}

	return string(password), nil
}
