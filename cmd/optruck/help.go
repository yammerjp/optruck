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
  optruck is a CLI tool for uploading secrets to 1Password Vaults and generating restoration templates. 
  By default, it **mirrors secrets**, meaning it uploads secrets from a data source to 1Password and generates a restoration template.

Arguments:
  <item>                Name of the 1Password item to process. Required unless --interactive is used.

Data Source Options (default: --env-file):
  Specify where to fetch secrets from. Choose one of the following:
  --env-file <path>     Path to the .env file containing secrets (default: ".env").
  --k8s-secret <name>   Name of the Kubernetes Secret to fetch secrets from.
                        When this option is used, you can also specify:
                        - --k8s-namespace <name>: Kubernetes namespace (default: "default").

Output Options:
  --output <path>       Path to save the restoration template file (default: ".env.1password").

General Options:
  --vault <value>       1Password Vault (e.g., "Development" or "abcd1234efgh5678").
  --account <value>     1Password account (e.g., "my.1password.com" or "my.1password.example.com").
  --overwrite           Overwrite the existing 1Password item if it exists.
  --interactive         Enable interactive mode for selecting the item, account, and vault.
                        In this mode, <item> is optional.
  --log-level <level>   Set the log level (debug|info|warn|error). Defaults to "info".
  -h, --help            Show help for optruck.
  --version            Show the version of optruck.

Examples:
  # Default: Mirror secrets from a .env file to 1Password and generate a template
  $ optruck MySecrets --vault MyVault --account my.1password.com
  # -> Uploads secrets from the .env file to 1Password and generates a restoration template.

  # Use a custom .env file for upload and template generation
  $ optruck MySecrets --env-file /path/to/custom.env
  # -> Uploads secrets from the specified .env file and generates a restoration template.

  # Upload secrets from Kubernetes Secret and generate a template in Kubernetes format
  $ optruck MySecrets --k8s-secret my-secret --output kube-secret.yaml --output-format k8s
  # -> Fetches secrets from the specified Kubernetes Secret, uploads them to 1Password, and generates a Kubernetes Secret manifest.

  # Specify Kubernetes namespace with Kubernetes Secret
  $ optruck MySecrets --k8s-secret my-secret --k8s-namespace my-namespace
  # -> Fetches secrets from the Kubernetes Secret in the specified namespace, then uploads them to 1Password.

  # Generate a restoration template only, without uploading secrets
  $ optruck MySecrets --template --output /path/to/template.env --output-format env
  # -> Generates a restoration template from the 1Password Vault.

  # Upload secrets only (no template generation)
  $ optruck MySecrets --upload --env-file /path/to/custom.env
  # -> Uploads secrets from the specified .env file to 1Password without generating a template.

  # Interactive mode for selecting item, account, and vault
  $ optruck --interactive
  # -> Allows you to select the item, Vault, and account interactively.

Notes:
  - When using Kubernetes options, ensure kubectl is configured properly for accessing the desired cluster and namespace.
`)

	return nil
}
