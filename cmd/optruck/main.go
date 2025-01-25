package optruck

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/yammerjp/optruck/pkg/actions"
	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
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
  --log-output <path>   Set the log output (<file path>). If not specified, output to stdout.
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
	Mirror   bool `name:"mirror" help:"Upload secrets and generate a restoration template (default)." default:"true"`

	// Data Source Options
	EnvFile      string `name:"env-file" help:"Path to the .env file containing secrets." default:".env"`
	K8sSecret    string `name:"k8s-secret" help:"Name of the Kubernetes Secret to fetch secrets from."`
	K8sNamespace string `name:"k8s-namespace" help:"Kubernetes namespace." default:"default"`

	// Output Options
	Output       string `name:"output" help:"Path to save the restoration template file." default:".env.1password"`
	OutputFormat string `name:"output-format" help:"Format of the output file (env|k8s)." enum:"env,k8s" default:"env"`

	// General Options
	Vault       string `name:"vault" help:"1Password Vault (e.g., 'Development' or 'abcd1234efgh5678')."`
	Account     string `name:"account" help:"1Password account (e.g., 'my.1password.com' or 'my.1password.example.com')."`
	Overwrite   bool   `name:"overwrite" help:"Overwrite the existing 1Password item and the output file if they exist."`
	Interactive bool   `name:"interactive" help:"Enable interactive mode for selecting the item, account, and vault."`

	// Misc
	Version   bool   `name:"version" help:"Show the version of optruck."`
	LogLevel  string `name:"log-level" help:"Set the log level (debug|info|warn|error)." enum:"debug,info,warn,error" default:"info"`
	LogOutput string `name:"log-output" help:"Set the log output destination. Use 'stdout' or 'stderr' for standard streams, or provide a file path." default:"stderr"`
}

func (cli CLI) validateItem(ctx *kong.Context) {
	if cli.Item == "" && !cli.Interactive {
		ctx.Fatalf("<item> is required unless --interactive or --help or --version is used")
	}
	if len(cli.Item) > 100 {
		// Limiting to 100 characters for simplicity; no deeper meaning behind this choice.
		ctx.Fatalf("item must be less than 100 characters")
	}
}

type ActionEnum int

const (
	ActionUpload ActionEnum = iota
	ActionTemplate
	ActionMirror
	ActionUnknown
)

func (cli CLI) validateAction(ctx *kong.Context) ActionEnum {
	if cli.Upload && !cli.Template && !cli.Mirror {
		// upload
		return ActionUpload
	}
	if !cli.Upload && cli.Template && !cli.Mirror {
		// template
		return ActionTemplate
	}
	if !cli.Upload && !cli.Template {
		// default or specified --mirror
		return ActionMirror
	}
	if !cli.Interactive {
		ctx.Fatalf("action must be one of upload, template, or mirror")
	}
	return ActionUnknown
}

func (cli CLI) BuildLogger(ctx *kong.Context) (*slog.Logger, func()) {
	var logLevel slog.Level
	switch cli.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		ctx.Fatalf("invalid log level: %s", cli.LogLevel)
	}
	var logWriter io.Writer
	var cleanup func() = func() {}

	if cli.LogOutput == "" {
		logWriter = os.Stderr
	} else {
		// Open log file
		f, err := os.OpenFile(cli.LogOutput, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			ctx.Fatalf("failed to open log file: %v", err)
		}
		logWriter = f
		cleanup = func() {
			f.Close()
		}
	}
	return slog.New(slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: logLevel})), cleanup
}

func Run() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("optruck"),
		kong.Description("A CLI tool for managing secrets and creating templates with 1Password."),
		kong.UsageOnError(),
		kong.Help(helpPrinter),
	)
	logger, cleanup := cli.BuildLogger(ctx)
	defer cleanup()

	// Handle version flag
	if cli.Version {
		logger.Debug("optruck version", "version", version)
		fmt.Printf("optruck version %s\n", version)
		os.Exit(0)
	}

	if cli.Interactive {
		ctx.Fatalf("interactive is not implemented")
	}

	cli.validateItem(ctx)

	action := cli.validateAction(ctx)
	switch action {
	case ActionUpload:
		err := cli.BuildUploadConfig(ctx, logger).Run()
		if err != nil {
			ctx.Fatalf("failed to upload secrets: %v", err)
		}
	case ActionTemplate:
		err := cli.BuildTemplateConfig(ctx, logger).Run()
		if err != nil {
			ctx.Fatalf("failed to generate template: %v", err)
		}
	case ActionMirror:
		err := cli.BuildMirrorConfig(ctx, logger).Run()
		if err != nil {
			ctx.Fatalf("failed to mirror secrets: %v", err)
		}
	default:
		ctx.Fatalf("invalid action: %v", action)
	}
}

func (cli CLI) BuildUploadConfig(ctx *kong.Context, logger *slog.Logger) actions.UploadConfig {
	return actions.UploadConfig{
		Logger:     logger,
		Target:     cli.BuildOpTarget(ctx),
		DataSource: cli.BuildDataSource(ctx),
		Overwrite:  cli.Overwrite,
	}
}

func (cli CLI) BuildTemplateConfig(ctx *kong.Context, logger *slog.Logger) actions.TemplateConfig {
	return actions.TemplateConfig{
		Logger:    logger,
		Target:    cli.BuildOpTarget(ctx),
		Dest:      cli.BuildDest(ctx),
		Overwrite: cli.Overwrite,
	}
}

func (cli CLI) BuildMirrorConfig(ctx *kong.Context, logger *slog.Logger) actions.MirrorConfig {
	return actions.MirrorConfig{
		Logger:     logger,
		Target:     cli.BuildOpTarget(ctx),
		DataSource: cli.BuildDataSource(ctx),
		Dest:       cli.BuildDest(ctx),
		Overwrite:  cli.Overwrite,
	}
}

func (cli CLI) BuildOpTarget(ctx *kong.Context) op.Target {
	if cli.Vault == "" {
		ctx.Fatalf("vault is required")
	}
	if cli.Account == "" {
		ctx.Fatalf("account is required")
	}

	return op.Target{
		Account:  cli.Account,
		Vault:    cli.Vault,
		ItemName: cli.Item,
	}
}

func (cli CLI) BuildDataSource(ctx *kong.Context) datasources.Source {
	if cli.K8sSecret != "" {
		if cli.EnvFile != "" {
			ctx.Fatalf("cannot use both --k8s-secret and --env-file")
		}
		if cli.K8sNamespace == "" {
			ctx.Fatalf("k8s namespace is required")
		}
		return &datasources.K8sSecretSource{
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
			Client:     datasources.NewK8sClient(),
		}
	}

	if cli.EnvFile == "" {
		ctx.Fatalf("env file is required")
	}
	if cli.K8sSecret != "" {
		ctx.Fatalf("cannot use both --env-file and --k8s-secret")
	}

	return &datasources.EnvFileSource{Path: cli.EnvFile}
}

func (cli CLI) BuildDest(ctx *kong.Context) output.Dest {
	if cli.OutputFormat == "k8s" {
		if cli.K8sSecret == "" {
			ctx.Fatalf("k8s secret is required")
		}
		if cli.K8sNamespace == "" {
			ctx.Fatalf("k8s namespace is required")
		}
		return &output.K8sSecretTemplateDest{
			Path:       cli.Output,
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
		}
	}

	if cli.OutputFormat == "env" {
		return &output.EnvTemplateDest{
			Path: cli.Output,
		}
	}

	if cli.EnvFile != "" {
		return &output.EnvTemplateDest{
			Path: cli.Output,
		}
	}
	if cli.K8sSecret != "" {
		if cli.K8sNamespace == "" {
			ctx.Fatalf("k8s namespace is required")
		}
		return &output.K8sSecretTemplateDest{
			Path:       cli.Output,
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
		}
	}

	ctx.Fatalf("invalid output format: %s", cli.OutputFormat)
	return nil
}
