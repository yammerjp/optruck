# optruck

optruck is a CLI tool for managing secrets and creating templates with 1Password. It can upload secrets from .env files (default) or Kubernetes Secrets to 1Password, and generate templates for restoring them later.

## Prerequisites

1. Install [1Password CLI](https://1password.com/downloads/command-line/)
   - If using 1Password GUI app with CLI integration enabled, you're already signed in
   - Otherwise, run `eval $(op signin)` to sign in

2. Install optruck:
```bash
go install github.com/yammerjp/optruck@latest
```

## Quick Start

The easiest way to get started is to use interactive mode:

```bash
optruck -i
```

This will guide you through selecting:
- The item name in 1Password
- The 1Password account and vault
- The data source (local file or Kubernetes secret)

## Usage

```bash
optruck <item> [options]
```

### Arguments

- `<item>`: Name to save the secrets as in 1Password. Required unless --interactive is used.

### Target Options

- `--vault <value>`: 1Password Vault (e.g., "Development" or "abcd1234efgh5678")
- `--account <value>`: 1Password account (e.g., "my.1password.com" or "my.1password.example.com")
- `--overwrite`: Overwrite the existing 1Password item if it exists

### Data Source Options

- `--env-file <path>`: Path to the .env file containing secrets (default: ".env")
- `--k8s-secret <name>`: Name of the Kubernetes Secret to fetch secrets from
- `--k8s-namespace <name>`: Kubernetes namespace for --k8s-secret (default: "default")

### Output Options

- `--output <path>`: Path to save the template file (default: ".env.1password" or "<secret-name>-secret.yaml.1password")

### General Options

- `-i, --interactive`: Enable interactive mode to select item, account, and vault
- `--log-level <level>`: Set the log level (debug|info|warn|error|none). Defaults to "none"
- `-h, --help`: Show help for optruck
- `-v, --version`: Show the version of optruck

## Examples

1. Start in interactive mode (recommended for first use):
```bash
optruck --interactive
```

2. Basic usage with .env file (default source):
```bash
optruck MySecrets --vault MyVault --account my.1password.com
```

3. Use a specific .env file:
```bash
optruck MySecrets --env-file /path/to/custom.env
```

4. Upload from Kubernetes Secret:
```bash
optruck MySecrets --k8s-secret my-secret --k8s-namespace my-namespace
# -> Generates "my-secret-secret.yaml.1password"
```

## Notes

- op (1Password CLI) must be installed and configured
- When using Kubernetes options, ensure kubectl is configured properly

## License

[MIT](LICENSE)
