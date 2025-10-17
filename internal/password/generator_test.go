package password

import (
	"strings"
	"testing"
	"unicode"
)

const (
	lettersSet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbersSet     = "0123456789"
	symbolsSet     = "!@#$%^&*()_+-=[]{}|;:,.<>?~"
	safeSymbolsSet = "!@#$%^&*()_+-=[]{}:,."
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
	pw, err := GeneratePassword(16, false, false, false)
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
	pw, err := GeneratePassword(24, true, false, false)
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
	pw, err := GeneratePassword(18, false, true, false)
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
	pw, err := GeneratePassword(32, true, true, false)
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
	pw, err := GeneratePassword(0, true, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw != "" {
		t.Fatalf("expected empty string for zero length, got %q", pw)
	}
}

func TestGenerate_NegativeLengthReturnsEmpty(t *testing.T) {
	// Negative length should return empty string without error
	pw, err := GeneratePassword(-1, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pw != "" {
		t.Fatalf("expected empty string for negative length, got %q", pw)
	}
}

func TestGenerate_LargeLength(t *testing.T) {
	const n = 1024
	pw, err := GeneratePassword(n, true, true, false)
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
		pw, err := GeneratePassword(20, true, false, false)
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
		pw, err := GeneratePassword(20, false, true, false)
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
		pw, err := GeneratePassword(30, true, true, false)
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
	pw, err := GeneratePassword(2, true, true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pw) != 2 {
		t.Fatalf("expected length 2, got %d", len(pw))
	}
}

func TestGenerate_SafeSymbols(t *testing.T) {
	// Generate with safe symbols only
	for i := 0; i < 10; i++ {
		pw, err := GeneratePassword(20, false, true, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		// Check that password only contains safe symbols
		allowed := lettersSet + safeSymbolsSet
		if !allRunesInSet(pw, allowed) {
			t.Fatalf("password contains characters outside safe symbols set: %q", pw)
		}
		
		// Ensure at least one safe symbol is present
		hasSymbol := false
		for _, r := range pw {
			if strings.ContainsRune(safeSymbolsSet, r) {
				hasSymbol = true
				break
			}
		}
		if !hasSymbol {
			t.Fatalf("password should contain at least one safe symbol: %q", pw)
		}
		
		// Ensure no unsafe symbols are present
		unsafeSymbols := "|;<>?~"
		for _, r := range pw {
			if strings.ContainsRune(unsafeSymbols, r) {
				t.Fatalf("password should not contain unsafe symbol %c: %q", r, pw)
			}
		}
	}
}

func TestGenerate_SafeSymbolsWithNumbers(t *testing.T) {
	// Generate with safe symbols and numbers
	for i := 0; i < 10; i++ {
		pw, err := GeneratePassword(25, true, true, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		// Check that password only contains allowed characters
		allowed := lettersSet + numbersSet + safeSymbolsSet
		if !allRunesInSet(pw, allowed) {
			t.Fatalf("password contains characters outside allowed set: %q", pw)
		}
		
		// Ensure all types are present
		hasLetter := false
		hasNumber := false
		hasSafeSymbol := false
		
		for _, r := range pw {
			if unicode.IsLetter(r) {
				hasLetter = true
			}
			if unicode.IsDigit(r) {
				hasNumber = true
			}
			if strings.ContainsRune(safeSymbolsSet, r) {
				hasSafeSymbol = true
			}
		}
		
		if !hasLetter {
			t.Fatalf("password should contain at least one letter: %q", pw)
		}
		if !hasNumber {
			t.Fatalf("password should contain at least one number: %q", pw)
		}
		if !hasSafeSymbol {
			t.Fatalf("password should contain at least one safe symbol: %q", pw)
		}
		
		// Ensure no unsafe symbols
		unsafeSymbols := "|;<>?~"
		for _, r := range pw {
			if strings.ContainsRune(unsafeSymbols, r) {
				t.Fatalf("password should not contain unsafe symbol %c: %q", r, pw)
			}
		}
	}
}

func TestGenerate_SafeSymbolsFlagWithoutSymbols(t *testing.T) {
	// Safe symbols flag should have no effect when symbols are disabled
	pw, err := GeneratePassword(20, true, false, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Should only contain letters and numbers
	allowed := lettersSet + numbersSet
	if !allRunesInSet(pw, allowed) {
		t.Fatalf("password contains characters outside letters and numbers: %q", pw)
	}
}

