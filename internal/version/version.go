package version

// Package version holds the application version information.

// These variables can be overridden at build time using -ldflags.
// Example:
//
//	go build -ldflags "-X github.com/jvanrhyn/passgen/internal/version.Version=v1.2.3 -X github.com/jvanrhyn/passgen/internal/version.Commit=abcdef -X github.com/jvanrhyn/passgen/internal/version.Date=2025-09-04T00:00:00Z"
var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)
