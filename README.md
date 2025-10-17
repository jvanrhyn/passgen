# Password Generator

## Description
Password Generator is a command-line application written in Go that allows users to generate secure passwords. The application provides options for specifying the desired length of the password, as well as including numbers and symbols. Additionally, users can copy the generated password to the clipboard.

## Usage
1. Clone the repository:
   ```bash
   git clone https://github.com/jvanrhyn/passgen.git
   cd passgen
   ```

2. Run the application with options:
   ```bash
   go run main.go --length <length> --numbers --symbols --clip
   ```

   or alternatively use the short form:
   ```bash
   go run main.go -l <length> -n -s -c
   ```

3. Run a compiled version of the application, build the application first:
   ```bash
   go build
   ```
   and then run the compiled application:
   ```bash
   ./passgen --length <length> --numbers --symbols --clip
   ```
   or alternatively use the short form:
   ```bash
   ./passgen -l <length> -n -s -c
   ```

### Options

- `--length, -l`: Specify the desired length of the password (default: 16, or value from `PASSWORD_LENGTH` environment variable)
- `--numbers, -n`: Include numbers in the password (guarantees at least one number)
- `--symbols, -s`: Include symbols in the password (guarantees at least one symbol)
- `--safe-symbols`: Use only shell-safe symbols, excluding potentially problematic characters (`|;<>?~`). Use with `-s` flag.
- `--clip, -c`: Copy the generated password to the clipboard (requires `CLIP_ALLOWED=true` environment variable)
- `--version, -v`: Show application version information

### Symbol Sets

When using the `--symbols` flag, you can choose between two symbol sets:

**Full symbol set (default):**
```
!@#$%^&*()_+-=[]{}|;:,.<>?~
```

**Safe symbol set (`--safe-symbols`):**
```
!@#$%^&*()_+-=[]{}:,.
```

The safe symbol set excludes characters that may cause issues in shell environments or certain systems:
- `|` (pipe) - command separator
- `;` (semicolon) - command separator  
- `<` `>` (angle brackets) - I/O redirection
- `?` (question mark) - glob pattern
- `~` (tilde) - home directory expansion

Use `--safe-symbols` when generating passwords that will be used in shell scripts, configuration files, or environments where these characters might be interpreted specially.

### Environment Variables (Optional)

- `PASSWORD_LENGTH`: Set default password length (default: 16)
- `CLIP_ALLOWED`: Enable clipboard functionality (`true` or `false`, default: disabled)

These can be set in a `.env` file in the project root, or as system environment variables.

### Examples

Generate a 20-character password with numbers and symbols:
```bash
./passgen -l 20 -n -s
```

Generate a 20-character password with numbers and safe symbols only:
```bash
./passgen -l 20 -n -s --safe-symbols
```

Generate a password using defaults (16 characters, letters only):
```bash
./passgen
```

Generate a password and copy to clipboard:
```bash
CLIP_ALLOWED=true ./passgen -l 24 -n -s -c
```

Generate a shell-safe password for use in scripts:
```bash
./passgen -l 32 -n -s --safe-symbols
```

## Design
The application is structured using best practices for Go projects:
- **cmd/**: Contains the command-line interface logic.
- **internal/**: Contains the password generation logic.

The application is a command-line interface without a visual user interface.
