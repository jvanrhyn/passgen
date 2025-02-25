package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/jvanrhyn/passgen/internal/password"
	"github.com/spf13/cobra"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	var length int
	var useNumbers bool
	var useSymbols bool

	var rootCmd = &cobra.Command{
		Use:   "passgen",
		Short: "Generate a secure password",
		Run: func(cmd *cobra.Command, args []string) {
			// Generate password
			password, err := password.GeneratePassword(length, useNumbers, useSymbols) // Use the flags for numbers and symbols
			if err != nil {
				fmt.Println("Error generating password:", err)
				os.Exit(1)
			}
			fmt.Println("Generated Password:", password)
		},
	}

	maxLengthStr := os.Getenv("PASSWORD_LENGTH")
	maxLength, err := strconv.Atoi(maxLengthStr)
	if err != nil {
		panic("PASSWORD_LENGTH environment variable is not set or invalid")
	}
	rootCmd.Flags().IntVarP(&length, "length", "l", maxLength, "Length of the password")
	rootCmd.Flags().BoolVarP(&useNumbers, "numbers", "n", false, "Include numbers in the password")
	rootCmd.Flags().BoolVarP(&useSymbols, "symbols", "s", false, "Include symbols in the password")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
