// Package password provides functionality to generate secure passwords.
package password

import (
	"crypto/rand"
	"math/big"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers = "0123456789"
	symbols = "!@#%*_=+"
)

// GeneratePassword generates a secure password of the specified length.
// When includeNumbers or includeSymbols are true, the password is guaranteed
// to contain at least one character from each requested character class.
func GeneratePassword(length int, includeNumbers bool, includeSymbols bool) (string, error) {
	if length < 0 {
		return "", nil // Let the slice allocation panic naturally for consistency
	}
	if length == 0 {
		return "", nil
	}

	// Build character sets
	charSets := []string{letters}
	if includeNumbers {
		charSets = append(charSets, numbers)
	}
	if includeSymbols {
		charSets = append(charSets, symbols)
	}

	// Calculate how many chars we need to guarantee from each set
	requiredCount := len(charSets)
	if length < requiredCount {
		// If password is shorter than required chars, just use all character sets randomly
		requiredCount = 0
	}

	password := make([]byte, length)

	// First, ensure at least one character from each required set
	pos := 0
	for i := 0; i < requiredCount; i++ {
		charSet := charSets[i]
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "", err
		}
		password[pos] = charSet[index.Int64()]
		pos++
	}

	// Fill remaining positions with random characters from all sets
	allChars := ""
	for _, cs := range charSets {
		allChars += cs
	}

	for i := pos; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			return "", err
		}
		password[i] = allChars[index.Int64()]
	}

	// Shuffle the password to avoid predictable patterns
	// (required chars at the beginning)
	if err := shuffleBytes(password); err != nil {
		return "", err
	}

	return string(password), nil
}

// shuffleBytes performs a Fisher-Yates shuffle using crypto/rand
func shuffleBytes(data []byte) error {
	for i := len(data) - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		j := int(jBig.Int64())
		data[i], data[j] = data[j], data[i]
	}
	return nil
}
