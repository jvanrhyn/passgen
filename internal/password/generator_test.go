package password

import (
	"strings"
	"testing"
	"unicode"
)

const (
	lettersSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbersSet = "0123456789"
	symbolsSet = "!@#$%^&*()_+-=[]{}|;:,.<>?~"
)

func allRunesInSet(s, allowed string) bool {
	allowedRunes := map[rune]struct{}{}
	for _, r := range allowed {
		allowedRunes[r] = struct{}{}
	}
	for _, r := range s {
		if _, ok := allowedRunes[r]; !ok {
			return false
		}
	}
	return true
}

func TestGenerate_LettersOnly(t *testing.T) {
	pw, err := GeneratePassword(16, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 16 {
		t.Fatalf("expected length 16, got %d", len(pw))
	}
	if !allRunesInSet(pw, lettersSet) {
		t.Fatalf("password contains characters outside letters set: %q", pw)
	}
}

func TestGenerate_LettersAndNumbers(t *testing.T) {
	pw, err := GeneratePassword(24, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 24 {
		t.Fatalf("expected length 24, got %d", len(pw))
	}
	allowed := lettersSet + numbersSet
	if !allRunesInSet(pw, allowed) {
		t.Fatalf("password contains characters outside allowed set: %q", pw)
	}
}

func TestGenerate_LettersAndSymbols(t *testing.T) {
	pw, err := GeneratePassword(18, false, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 18 {
		t.Fatalf("expected length 18, got %d", len(pw))
	}
	allowed := lettersSet + symbolsSet
	if !allRunesInSet(pw, allowed) {
		t.Fatalf("password contains characters outside allowed set: %q", pw)
	}
}

func TestGenerate_AllCharacterClasses(t *testing.T) {
	pw, err := GeneratePassword(32, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 32 {
		t.Fatalf("expected length 32, got %d", len(pw))
	}
	allowed := lettersSet + numbersSet + symbolsSet
	if !allRunesInSet(pw, allowed) {
		t.Fatalf("password contains characters outside allowed set: %q", pw)
	}
}

func TestGenerate_ZeroLength(t *testing.T) {
	pw, err := GeneratePassword(0, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw != "" {
		t.Fatalf("expected empty string for zero length, got %q", pw)
	}
}

func TestGenerate_NegativeLengthReturnsEmpty(t *testing.T) {
	// Negative length should return empty string without error
	pw, err := GeneratePassword(-1, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw != "" {
		t.Fatalf("expected empty string for negative length, got %q", pw)
	}
}

func TestGenerate_LargeLength(t *testing.T) {
	const n = 1024
	pw, err := GeneratePassword(n, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != n {
		t.Fatalf("expected length %d, got %d", n, len(pw))
	}
	// quick sanity: string should be printable characters from the set
	for _, r := range pw {
		if unicode.IsSpace(r) {
			t.Fatalf("unexpected whitespace in password")
		}
	}
}

func TestGenerate_GuaranteesNumbers(t *testing.T) {
	// Generate multiple times to ensure numbers are always present
	for i := 0; i < 10; i++ {
		pw, err := GeneratePassword(20, true, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		hasNumber := false
		for _, r := range pw {
			if unicode.IsDigit(r) {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			t.Fatalf("password should contain at least one number: %q", pw)
		}
	}
}

func TestGenerate_GuaranteesSymbols(t *testing.T) {
	// Generate multiple times to ensure symbols are always present
	for i := 0; i < 10; i++ {
		pw, err := GeneratePassword(20, false, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		hasSymbol := false
		for _, r := range pw {
			if strings.ContainsRune(symbolsSet, r) {
				hasSymbol = true
				break
			}
		}
		if !hasSymbol {
			t.Fatalf("password should contain at least one symbol: %q", pw)
		}
	}
}

func TestGenerate_GuaranteesAllCharacterTypes(t *testing.T) {
	// Generate multiple times to ensure all types are present
	for i := 0; i < 10; i++ {
		pw, err := GeneratePassword(30, true, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		hasLetter := false
		hasNumber := false
		hasSymbol := false
		
		for _, r := range pw {
			if unicode.IsLetter(r) {
				hasLetter = true
			}
			if unicode.IsDigit(r) {
				hasNumber = true
			}
			if strings.ContainsRune(symbolsSet, r) {
				hasSymbol = true
			}
		}
		
		if !hasLetter {
			t.Fatalf("password should contain at least one letter: %q", pw)
		}
		if !hasNumber {
			t.Fatalf("password should contain at least one number: %q", pw)
		}
		if !hasSymbol {
			t.Fatalf("password should contain at least one symbol: %q", pw)
		}
	}
}

func TestGenerate_ShortPasswordWithAllTypes(t *testing.T) {
	// When password length is less than number of character types,
	// it should still work without error
	pw, err := GeneratePassword(2, true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 2 {
		t.Fatalf("expected length 2, got %d", len(pw))
	}
}
