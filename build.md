# Generic Go Builder (build.go)

This repository includes a small helper script `build.go` that builds all binaries in the module into `./bin` with automatic build number tracking and version injection.

It discovers all `main` packages in the module automatically (including the root package and any subfolders like `cmd/...`) and produces an executable for each.

## Quick start

- Build everything into `./bin` (dev build):

```sh
go run build.go
```

- Build a release version (without -dev tag):

```sh
go run build.go -release
# or
RELEASE=true go run build.go
```

- Outputs: one executable per `main` package.
  - The binary name comes from the package's folder name.
  - The module root builds to a binary named after the repository folder (e.g., `passgen`).
  - On Windows targets, `.exe` is appended.

## Build Number Tracking

The build system automatically tracks and increments build numbers for every build:

- Build numbers are stored in `.buildnum` file (auto-created if missing)
- Every build increments the build number by 1
- Dev builds append `-dev` tag to the build number (e.g., `5-dev`)
- Release builds use the plain build number (e.g., `5`)
- Build numbers are injected into the binary via `-ldflags`

The version information is accessible via the `internal/version` package and displayed with the `--version` flag.

Example output:
```
Dev build:     passgen dev (build 5-dev) [99cbc02] built 2025-10-17T19:19:24Z
Release build: passgen dev (build 5) [99cbc02] built 2025-10-17T19:19:24Z
```

## Optional .build file

You can include extra build targets via a `.build` file placed in the repository root. Each non-empty, non-comment line should contain one of:

- An import path (e.g., `./cmd/foo` or `github.com/you/repo/cmd/foo`)
- A relative directory path

Lines starting with `#` are ignored.

Example `.build`:

```
# Additional tools
./cmd/admin
./tools/migrate
```

These entries will be merged with the auto-discovered targets.

## Cross-compiling

`build.go` respects the `GOOS` environment variable for the output file extension:

```sh
# Cross-compile for Windows from macOS/Linux
GOOS=windows GOARCH=amd64 go run build.go
```

- The builder will append `.exe` if `GOOS=windows`.
- Standard Go cross-compilation rules apply; set `GOARCH` and other env vars as needed.

## Notes

- All binaries are written to `./bin` (created if missing).
- Build numbers are automatically incremented with each build.
- The `.buildnum` file is ignored by git (listed in `.gitignore`).
- Git commit hash and build date are automatically captured and injected.
- Stable output ordering for predictable diffs.
- If no `main` packages are found, the builder exits with an error.

## Troubleshooting

- Ensure you run from the module root (where `go.mod` resides).
- If `go list ./...` fails, confirm the module builds normally with `go build ./...`.
- For private modules, make sure your environment has access tokens configured for `go`.
- If build numbers seem incorrect, check the `.buildnum` file content.
