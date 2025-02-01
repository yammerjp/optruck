# optruck

**Note:** This project is currently under active development and not yet stable.

optruck is a CLI tool that streamlines secret management with 1Password by providing a simple way to upload secrets and generate restoration templates.

## Features

- Upload secrets from environment files or Kubernetes Secrets to 1Password
- Generate restoration templates from 1Password items
- Support for both environment file and Kubernetes Secret formats
- Interactive mode for selecting items, accounts, and vaults

## Installation

```bash
go install github.com/yammerjp/optruck@latest
```

## Getting Started

### Prerequisites

1. Install [1Password CLI](https://1password.com/downloads/command-line/)
   - If using 1Password GUI app with CLI integration enabled, you're already signed in
   - Otherwise, run `eval $(op signin)` to sign in

2. Install optruck:
```bash
go install github.com/yammerjp/optruck@latest
```

### Quick Start

1. Create a `.env` file with your secrets:
```bash
# .env
API_KEY=your-secret-api-key
DATABASE_URL=your-database-url
```

2. Upload secrets to 1Password and generate a template:
```bash
optruck MySecrets --vault MyVault --account my.1password.com
```

This will:
- Create a new item named "MySecrets" in your 1Password vault
- Upload the secrets from `.env`
- Generate a template file `.env.1password` for restoration

3. To restore secrets later, use the generated template with 1Password CLI:
```bash
op inject -i .env.1password -o .env
```

## Common Use Cases

1. Mirror secrets from a .env file (default behavior):
```bash
optruck MySecrets --vault MyVault --account my.1password.com
```

2. Use a custom .env file:
```bash
optruck MySecrets --env-file /path/to/custom.env
```

3. Upload Kubernetes Secret and generate K8s template:
```bash
optruck MySecrets --k8s-secret my-secret --output kube-secret.yaml
```

4. Specify Kubernetes namespace:
```bash
optruck MySecrets --k8s-secret my-secret --k8s-namespace my-namespace
```

## Advanced Usage

### Data Source Options (default: --env-file)

- `--env-file <path>`: Path to the .env file containing secrets (default: ".env")
- `--k8s-secret <name>`: Name of the Kubernetes Secret to fetch secrets from
- `--k8s-namespace <name>`: Kubernetes namespace (default: "default")

### Output Options

- `--output <path>`: Path to save the restoration template file (default: ".env.1password")
- `--output-format <env|k8s>`: Format of the output file
  - "env" for environment variable files
  - "k8s" for Kubernetes Secret manifests

### General Options

- `--vault <value>`: 1Password Vault (e.g., "Development" or "abcd1234efgh5678")
- `--account <value>`: 1Password account (e.g., "my.1password.com" or "my.1password.example.com")
- `--overwrite`: Overwrite existing 1Password item and output file
- `--interactive`: Enable interactive mode for selecting item, account, and vault
- `--log-level <level>`: Set log level (debug|info|warn|error)
- `--log-output <path>`: Set the log output (<file path>). If not specified, no logging is done.
- `-h, --help`: Show help
- `--version`: Show version

## Notes

- Use `--overwrite` to update existing 1Password items 
- When using Kubernetes options, ensure kubectl is properly configured
- The tool requires appropriate 1Password CLI configuration and authentication

## License

[MIT](LICENSE)


## TODO

- [ ] Add tests for interactive mode
- [ ] More descriptive error messages to help user
- [ ] Fix logging
- [ ] Add tests for kube client
- [ ] E2E testing (with dummy commands(kubectl, op))
- [ ] Register kubernetes secret base64 decoded data as 1Password item and mark it as 'decoded'
- [ ] Support more data source of kubernetes secret (not only type Opaque)
- [ ] Allow to select a part of secret on .env file to upload
- [ ] Add to check op cli and kubectl are installed
