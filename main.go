package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"github.com/jvanrhyn/passgen/internal/password"
	"github.com/jvanrhyn/passgen/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	// Load environment variables from .env file if it exists (optional)
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	var length int
	var useNumbers bool
	var useSymbols bool
	var clip bool
	var showVersion bool

	var rootCmd = &cobra.Command{
		Use:   "passgen",
		Short: "Generate a secure password",
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				fmt.Printf("passgen %s", version.Version)
				if version.BuildNumber != "" && version.BuildNumber != "0" {
					fmt.Printf(" (build %s)", version.BuildNumber)
				}
				if version.Commit != "" {
					fmt.Printf(" [%s]", version.Commit)
				}
				if version.Date != "" {
					fmt.Printf(" built %s", version.Date)
				}
				fmt.Println()
				return nil
			}

			// Validate password length
			if length < 0 {
				return fmt.Errorf("password length must be non-negative, got: %d", length)
			}
			if length > 10000 {
				return fmt.Errorf("password length too large (max 10000), got: %d", length)
			}

			// Generate password
			pwd, err := password.GeneratePassword(length, useNumbers, useSymbols)
			if err != nil {
				return fmt.Errorf("error generating password: %w", err)
			}

			// Print the generated password
			fmt.Println("Generated Password:", pwd)

			// Check if clipboard copying is allowed
			clipAllowed := os.Getenv("CLIP_ALLOWED")
			if clipAllowed == "true" && clip {
				if err := clipboard.WriteAll(pwd); err != nil {
					return fmt.Errorf("error copying to clipboard: %w", err)
				}
				fmt.Println("Password copied to clipboard")
			} else if clip {
				fmt.Println("Clipboard copying is disabled by environment settings")
			}

			return nil
		},
	}

	// Get default password length from environment or use 16 as default
	defaultLength := 16
	if maxLengthStr := os.Getenv("PASSWORD_LENGTH"); maxLengthStr != "" {
		if parsedLength, err := strconv.Atoi(maxLengthStr); err == nil && parsedLength > 0 {
			defaultLength = parsedLength
		}
	}

	rootCmd.Flags().IntVarP(&length, "length", "l", defaultLength, "Length of the password")
	rootCmd.Flags().BoolVarP(&useNumbers, "numbers", "n", false, "Include numbers in the password")
	rootCmd.Flags().BoolVarP(&useSymbols, "symbols", "s", false, "Include symbols in the password")
	rootCmd.Flags().BoolVarP(&clip, "clip", "c", false, "Copy password to clipboard")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "Show application version")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
