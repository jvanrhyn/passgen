# Password Generator

## Description
Password Generator is a simple command-line application written in Go that allows users to generate secure passwords. The application provides an intuitive interface for users to specify the desired length of the password.

## Usage
1. Clone the repository:
   ```bash
   git clone https://github.com/jvanrhyn/passgen.git
   cd passgen
   ```

2. Run the application with options:
   ```bash
   go run main.go --length <length> --numbers --symbols
   ```

   - `--length`: Specify the desired length of the password (default is set by the `PASSWORD_LENGTH` environment variable).
   - `--numbers`: Include numbers in the password.
   - `--symbols`: Include symbols in the password.

## Design
The application is structured using best practices for Go projects:
- **cmd/**: Contains the command-line interface logic.
- **internal/**: Contains the password generation logic.
- **cmd/**: Contains the command-line interface logic.
- **internal/**: Contains the password generation logic.

The application is a command-line interface without a visual user interface.
