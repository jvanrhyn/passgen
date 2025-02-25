package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "passgen",
	Short: "A simple password generator",
	Long:  `A simple command line application to generate secure passwords.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
