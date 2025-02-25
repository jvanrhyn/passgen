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

   - `--length`: Specify the desired length of the password (default is set by the `PASSWORD_LENGTH` environment variable).
   - `--numbers`: Include numbers in the password.
   - `--symbols`: Include symbols in the password.
   - `--clip`: Copy the generated password to the clipboard if allowed by the `CLIP_ALLOWED` environment variable.

## Design
The application is structured using best practices for Go projects:
- **cmd/**: Contains the command-line interface logic.
- **internal/**: Contains the password generation logic.
- **cmd/**: Contains the command-line interface logic.
- **internal/**: Contains the password generation logic.

The application is a command-line interface without a visual user interface.
