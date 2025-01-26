package optruck

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/yammerjp/optruck/pkg/config"
)

var version = "0.1.0"

func helpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	fmt.Print(`optruck - A CLI tool for managing secrets and creating templates with 1Password.

Usage:
  optruck <item> [options]

Description:
  optruck is a CLI tool for uploading secrets to 1Password Vaults and generating restoration templates. 
  By default, it **mirrors secrets**, meaning it uploads secrets from a data source to 1Password and generates a restoration template (--mirror).

Arguments:
  <item>                Name of the 1Password item to process. Required unless --interactive is used.

Actions (default: --mirror):
  --upload              Upload secrets to 1Password Vault.
                        Requires a data source (e.g., --env-file or --k8s-secret).
  --template            Generate a restoration template from the 1Password Vault.
                        Does not require a data source.
  --mirror              Upload secrets and generate a restoration template (default).
                        Combines the functionality of --upload and --template.

Data Source Options (default: --env-file):
  Specify where to fetch secrets from. Choose one of the following:
  --env-file <path>     Path to the .env file containing secrets (default: ".env").
  --k8s-secret <name>   Name of the Kubernetes Secret to fetch secrets from.
                        When this option is used, you can also specify:
                        - --k8s-namespace <name>: Kubernetes namespace (default: "default").

Output Options:
  --output <path>       Path to save the restoration template file (default: ".env.1password").
  --output-format <env|k8s>
                        Format of the output file:
                          - "env" for environment variable files.
                          - "k8s" for Kubernetes Secret manifests.
                        Defaults to "env" unless --k8s-* is used, in which case "k8s" is applied.

General Options:
  --vault <value>       1Password Vault (e.g., "Development" or "abcd1234efgh5678").
  --account <value>     1Password account (e.g., "my.1password.com" or "my.1password.example.com").
  --overwrite           Overwrite the existing 1Password item and the output file if they exist.
  --interactive         Enable interactive mode for selecting the item, account, and vault.
                        In this mode, <item> is optional.
  --log-level <level>   Set the log level (debug|info|warn|error). Defaults to "info".
  --log-output <path>   Set the log output (<file path>). If not specified, no logging is done.
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
  $ optruck MySecrets --mirror --k8s-secret my-secret --output kube-secret.yaml --output-format k8s
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
  - Default action is --mirror, which uploads secrets to 1Password and generates a restoration template.
  - Use --overwrite to update both the existing 1Password item and output file.
  - When using Kubernetes options, ensure kubectl is configured properly for accessing the desired cluster and namespace.
`)

	return nil
}

type CLI struct {
	Item string `arg:"" optional:"" name:"item" help:"Name of the 1Password item to process. Required unless --interactive is used."`

	// Actions
	Upload   bool `name:"upload" help:"Upload secrets to 1Password Vault."`
	Template bool `name:"template" help:"Generate a restoration template from the 1Password Vault."`
	Mirror   bool `name:"mirror" help:"Upload secrets and generate a restoration template (default)."` // Don't set kong's default value to distinguish mirror option in interactive mode

	// Data Source Options
	EnvFile      string `name:"env-file" help:"Path to the .env file containing secrets. (default: '.env')"` // Don't set kong's default value to distinguish env-file option in interactive mode
	K8sSecret    string `name:"k8s-secret" help:"Name of the Kubernetes Secret to fetch secrets from."`
	K8sNamespace string `name:"k8s-namespace" help:"Kubernetes namespace."` // Don't set kong's default value to distinguish k8s-namespace option in interactive mode

	// Output Options
	Output       string `name:"output" help:"Path to save the restoration template file. (default: '.env.1password' if format is env, otherwise '<name>-secret.yaml.1password' if format is k8s)"` // Don't set kong's default value to distinguish output option in interactive mode
	OutputFormat string `name:"output-format" help:"Format of the output file (env|k8s)."`                                                                                                         // Don't set kong's default value and don't use enum to distinguish output-format option in interactive mode

	// General Options
	Vault             string `name:"vault" help:"1Password Vault (e.g., 'Development' or 'abcd1234efgh5678')."`
	Account           string `name:"account" help:"1Password account (e.g., 'my.1password.com' or 'my.1password.example.com')."`
	Overwrite         bool   `name:"overwrite" help:"Overwrite the existing 1Password item and the output file if they exist."`
	OverwriteTarget   bool   `name:"overwrite-target" help:"Overwrite the existing 1Password item if it exists."`
	OverwriteTemplate bool   `name:"overwrite-template" help:"Overwrite the existing output file if it exists."`
	Interactive       bool   `name:"interactive" help:"Enable interactive mode for selecting the item, account, and vault." short:"i"`

	// Misc
	Version   bool   `name:"version" help:"Show the version of optruck."`
	LogLevel  string `name:"log-level" help:"Set the log level (debug|info|warn|error)." enum:"debug,info,warn,error" default:"info"`
	LogOutput string `name:"log-output" help:"Set the log output (<file path>). If not specified, no logging is done."`
}

func Run() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("optruck"),
		kong.Description("A CLI tool for managing secrets and creating templates with 1Password."),
		kong.UsageOnError(),
		kong.Help(helpPrinter),
	)

	// Handle version flag
	if cli.Version {
		fmt.Printf("optruck version %s\n", version)
		os.Exit(0)
	}

	builder := config.NewConfigBuilder().
		WithItem(cli.Item).
		WithVault(cli.Vault).
		WithAccount(cli.Account).
		WithEnvFile(cli.EnvFile).
		WithK8sSecret(cli.K8sSecret).
		WithK8sNamespace(cli.K8sNamespace).
		WithOutput(cli.Output).
		WithOutputFormat(cli.OutputFormat).
		WithOverwrite(cli.Overwrite).
		WithOverwriteTarget(cli.OverwriteTarget).
		WithOverwriteTemplate(cli.OverwriteTemplate).
		WithLogLevel(cli.LogLevel).
		WithLogFile(cli.LogOutput).
		WithUpload(cli.Upload).
		WithTemplate(cli.Template).
		WithMirror(cli.Mirror)

	if cli.Interactive {
		err := builder.SetConfigInteractively()
		if err != nil {
			ctx.Fatalf("%v", err)
		}
	}
	err := builder.SetDefaultIfEmpty()
	if err != nil {
		ctx.Fatalf("%v", err)
	}

	action, cleanup, err := builder.Build()
	if err != nil {
		ctx.Fatalf("%v", err)
	}

	defer cleanup()
	if err := action.Run(); err != nil {
		ctx.Fatalf("%v", err)
	}
}
