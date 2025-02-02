package optruck

import (
	"fmt"

	"github.com/alecthomas/kong"
)

func helpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	fmt.Print(`optruck - A CLI tool for managing secrets and creating templates with 1Password.

Usage:
  optruck <item> [options]

Description:
  optruck helps you manage application secrets using 1Password. It can upload secrets from
  .env files (default) or Kubernetes Secrets to 1Password, and generate templates for
  restoring them later.

Arguments:
  <item>                Name to save the secrets as in 1Password. Required unless --interactive is used.

Target Options:
  --vault <value>       1Password Vault (e.g., "Development" or "abcd1234efgh5678").
  --account <value>     1Password account (e.g., "my.1password.com" or "my.1password.example.com").
  --overwrite           Overwrite the existing 1Password item if it exists.

Data Source Options:
  --env-file <path>     Path to the .env file containing secrets (default: ".env").
  --k8s-secret <name>   Name of the Kubernetes Secret to fetch secrets from.
  --k8s-namespace <name> Kubernetes namespace for --k8s-secret (default: "default").

Output Options:
  --output <path>       Path to save the template file (default: ".env.1password" or "<secret-name>-secret.yaml.1password").

General Options:
  -i, --interactive     Enable interactive mode to select item, account, and vault.
  --log-level <level>   Set the log level (debug|info|warn|error|none). Defaults to "none".
  -h, --help           Show help for optruck.
  -v, --version        Show the version of optruck.

Examples:
  # Start in interactive mode (recommended for first use)
  $ optruck --interactive

  # Basic usage with .env file (default source)
  $ optruck MySecrets --vault MyVault --account my.1password.com

  # Use a specific .env file
  $ optruck MySecrets --env-file /path/to/custom.env

  # Upload from Kubernetes Secret (generates YAML template)
  $ optruck MySecrets --k8s-secret my-secret --k8s-namespace my-namespace
  # -> Generates "my-secret-secret.yaml.1password"

Notes:
  - op (1Password CLI) must be installed and configured.
  - When using Kubernetes options, ensure kubectl is configured properly.
`)

	return nil
}
