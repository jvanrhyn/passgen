package cmd

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// buildAndRun builds the binary in a temp dir and executes it with provided args and env.
// It returns stdout and stderr content along with any execution error.
func buildAndRun(t *testing.T, args []string, env map[string]string) (string, string, error) {
	t.Helper()

	// Create temp working directory
	tempDir := t.TempDir()

	// Path to repo root where go.mod and main.go reside (one level up from this file)
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	// Copy a minimal .env into temp dir if env values provided
	var envBuf bytes.Buffer
	if env != nil {
		for k, v := range env {
			envBuf.WriteString(k + "=" + v + "\n")
		}
		if err := os.WriteFile(filepath.Join(repoRoot, ".env"), envBuf.Bytes(), 0o600); err != nil {
			t.Fatalf("failed to write .env: %v", err)
		}
		t.Cleanup(func() { _ = os.Remove(filepath.Join(repoRoot, ".env")) })
	}

	// Build the binary into tempDir (add .exe on Windows)
	binName := "passgen-test-bin"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(tempDir, binName)
	buildCmd := exec.Command("go", "build", "-o", binPath, "./")
	buildCmd.Dir = repoRoot
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, string(buildOut))
	}

	// Prepare execution
	runCmd := exec.Command(binPath, args...)
	runCmd.Dir = repoRoot
	// inherit existing env but ensure CLIP_ALLOWED and PASSWORD_LENGTH from .env are used by app loading .env
	runCmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr
	err = runCmd.Run()
	return stdout.String(), stderr.String(), err
}

func readGeneratedPassword(stdout string) string {
	// The app prints: "Generated Password: <password>" on a line
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Generated Password:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Generated Password:"))
		}
	}
	return ""
}

func TestCLI_DefaultLengthFromEnv(t *testing.T) {
	// .env defines PASSWORD_LENGTH=12, CLIP_ALLOWED=false to avoid clipboard dependency
	stdout, stderr, err := buildAndRun(t, []string{}, map[string]string{
		"PASSWORD_LENGTH": "12",
		"CLIP_ALLOWED":    "false",
	})
	if err != nil {
		t.Fatalf("run error: %v, stderr: %s", err, stderr)
	}
	pwd := readGeneratedPassword(stdout)
	if len(pwd) != 12 {
		t.Fatalf("expected default length 12, got %d (out: %q)", len(pwd), stdout)
	}
}

func TestCLI_LengthFlagOverridesEnv(t *testing.T) {
	stdout, stderr, err := buildAndRun(t, []string{"-l", "20"}, map[string]string{
		"PASSWORD_LENGTH": "10",
		"CLIP_ALLOWED":    "false",
	})
	if err != nil {
		t.Fatalf("run error: %v, stderr: %s", err, stderr)
	}
	pwd := readGeneratedPassword(stdout)
	if len(pwd) != 20 {
		t.Fatalf("expected length 20, got %d", len(pwd))
	}
}

func TestCLI_NumbersAndSymbolsFlagsAffectCharacterSet(t *testing.T) {
	stdout, stderr, err := buildAndRun(t, []string{"-l", "50", "-n", "-s"}, map[string]string{
		"PASSWORD_LENGTH": "8",
		"CLIP_ALLOWED":    "false",
	})
	if err != nil {
		t.Fatalf("run error: %v, stderr: %s", err, stderr)
	}
	pwd := readGeneratedPassword(stdout)
	if len(pwd) != 50 {
		t.Fatalf("expected length 50, got %d", len(pwd))
	}
	// ensure only characters from union set appear
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#%*_=+"
	re := regexp.MustCompile("^[" + regexp.QuoteMeta(allowed) + "]+$")
	if !re.MatchString(pwd) {
		t.Fatalf("password contains invalid characters: %q", pwd)
	}
}

func TestCLI_ClipboardDisabledMessage(t *testing.T) {
	stdout, _, err := buildAndRun(t, []string{"-c"}, map[string]string{
		"PASSWORD_LENGTH": "8",
		"CLIP_ALLOWED":    "false",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "Clipboard copying is disabled by environment settings") {
		t.Fatalf("expected clipboard disabled message, got: %q", stdout)
	}
}

func TestCLI_ClipboardEnabledAttemptsCopy(t *testing.T) {
	// When clipboard allowed and -c passed, it should attempt copy and print success.
	// On CI without clipboard provider, atotto/clipboard may still succeed on macOS.
	// We at least verify the success message appears when allowed.
	stdout, _, err := buildAndRun(t, []string{"-c"}, map[string]string{
		"PASSWORD_LENGTH": "8",
		"CLIP_ALLOWED":    "true",
	})
	if err != nil {
		// If clipboard backend missing, the app would exit non-zero after printing error.
		// In that case, accept either success text or an error line. We don't fail the test.
		// But prefer to pass when message exists.
	}
	if !(strings.Contains(stdout, "Password copied to clipboard") || strings.Contains(stdout, "Error copying to clipboard:")) {
		t.Fatalf("expected clipboard action message, got: %q", stdout)
	}
}

func TestCLI_MissingDotEnvFailsEarly(t *testing.T) {
	// When .env is missing, the app exits with an error about loading .env
	stdout, _, err := buildAndRun(t, []string{}, nil)
	if err == nil {
		t.Fatalf("expected non-zero exit when .env is missing, stdout=%q", stdout)
	}
	if !strings.Contains(stdout, "Error loading .env file") {
		t.Fatalf("expected error about loading .env, got: %q", stdout)
	}
}

func TestCLI_LettersOnly_NoFlags(t *testing.T) {
	// With no -n or -s, ensure only letters are used and default length applies.
	stdout, stderr, err := buildAndRun(t, []string{}, map[string]string{
		"PASSWORD_LENGTH": "10",
		"CLIP_ALLOWED":    "false",
	})
	if err != nil {
		t.Fatalf("run error: %v, stderr: %s", err, stderr)
	}
	pwd := readGeneratedPassword(stdout)
	if len(pwd) != 10 {
		t.Fatalf("expected length 10, got %d", len(pwd))
	}
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	re := regexp.MustCompile("^[" + regexp.QuoteMeta(allowed) + "]+$")
	if !re.MatchString(pwd) {
		t.Fatalf("letters-only expected, got: %q", pwd)
	}
}

func TestCLI_ZeroLength(t *testing.T) {
	stdout, stderr, err := buildAndRun(t, []string{"-l", "0"}, map[string]string{
		"PASSWORD_LENGTH": "8",
		"CLIP_ALLOWED":    "false",
	})
	if err != nil {
		t.Fatalf("run error: %v, stderr: %s", err, stderr)
	}
	pwd := readGeneratedPassword(stdout)
	if len(pwd) != 0 {
		t.Fatalf("expected empty password for length 0, got length %d", len(pwd))
	}
}

func TestCLI_NegativeLengthPanics(t *testing.T) {
	// Provide .env so loading succeeds, then override with -l -1 to trigger panic in generator
	_, _, err := buildAndRun(t, []string{"-l", "-1"}, map[string]string{
		"PASSWORD_LENGTH": "8",
		"CLIP_ALLOWED":    "false",
	})
	if err == nil {
		t.Fatalf("expected non-zero exit for negative length")
	}
}
