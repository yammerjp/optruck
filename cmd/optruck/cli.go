package optruck

type CLI struct {
	Item string `arg:"" optional:"" name:"item" help:"Name of the 1Password item to process. Required unless --interactive is used."`

	// Data Source Options
	EnvFile      string `name:"env-file" help:"Path to the .env file containing secrets. (default: '.env')"` // Don't set kong's default value to distinguish env-file option in interactive mode
	K8sSecret    string `name:"k8s-secret" help:"Name of the Kubernetes Secret to fetch secrets from."`
	K8sNamespace string `name:"k8s-namespace" help:"Kubernetes namespace."` // Don't set kong's default value to distinguish k8s-namespace option in interactive mode

	// Output Options
	Output string `name:"output" help:"Path to save the restoration template file. (default: '.env.1password' if format is env, otherwise '<name>-secret.yaml.1password' if format is k8s)"` // Don't set kong's default value to distinguish output option in interactive mode

	// General Options
	Vault       string `name:"vault" help:"1Password Vault (e.g., 'Development' or 'abcd1234efgh5678')."`
	Account     string `name:"account" help:"1Password account (e.g., 'my.1password.com' or 'my.1password.example.com')."`
	Overwrite   bool   `name:"overwrite" help:"Overwrite the existing 1Password item if it exists."`
	Interactive bool   `name:"interactive" help:"Enable interactive mode for selecting the item, account, and vault." short:"i"`

	// Misc
	Version  VersionFlag `short:"v" help:"Show the version of optruck."`
	LogLevel string      `name:"log-level" help:"Set the log level (debug|info|warn|error)."`
}
