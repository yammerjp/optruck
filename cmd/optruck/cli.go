package optruck

import "github.com/yammerjp/optruck/internal/util/interactiverunner"

type InteractiveFlag bool

type CLI struct {
	// max length is 100
	Item string `arg:"" optional:"" name:"item" help:"Name or ID of the 1Password item to process."`

	// Target Options
	Account   string `name:"account" help:"1Password account (e.g., 'my.1password.com' or 'my.1password.example.com')."`
	Vault     string `name:"vault" help:"1Password Vault Name or ID (e.g., 'Development' or 'abcd1234efgh5678')."`
	Overwrite bool   `name:"overwrite" help:"Overwrite the existing 1Password item if it exists."`

	// Data Source Options
	EnvFile      string `name:"env-file" type:"existingfile" optional:"" help:"Path to the .env file containing secrets.(default: '.env')" xor:"source-type,source-k8s"`
	K8sSecret    string `name:"k8s-secret" optional:"" help:"Name of the Kubernetes Secret to fetch secrets from." xor:"source-type"`
	K8sNamespace string `name:"k8s-namespace" optional:"" help:"Kubernetes namespace.(default: 'default')" xor:"source-k8s"`

	// Output Options
	Output string `name:"output" type:"path" help:"Path to save the restoration template file. (default: '.env.1password' if format is env, otherwise '<name>-secret.yaml.1password' if format is k8s)"` // Don't set kong's default value

	// General Options
	Version     VersionFlag     `short:"v" help:"Show the version of optruck."`
	LogLevel    string          `name:"log-level" help:"Set the log level (debug|info|warn|error|none)." enum:"debug,info,warn,error,none" default:"none"`
	Interactive InteractiveFlag `name:"interactive" help:"Enable interactive mode for selecting the item, account, and vault." short:"i"`

	runner interactiverunner.InteractiveRunner
}
