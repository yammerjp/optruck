package cmd

import (
	"github.com/alecthomas/kong"
)

// CLI defines the structure of the CLI commands and flags.
type CLI struct {
	Upload   UploadCmd   `cmd:"" help:"Upload secrets to a 1Password Vault."`
	Template TemplateCmd `cmd:"" help:"Generate a restoration template file."`
	Mirror   MirrorCmd   `cmd:"" help:"Upload secrets and generate a template in one step."`
	Version  VersionCmd  `cmd:"" help:"Show version information."`
}

// SharedFlags defines flags that are common across multiple commands.
type SharedFlags struct {
	Vault   string `help:"Name of the 1Password Vault."`
	Account string `help:"1Password account email address."`
	Item    string `help:"Name of the 1Password item where secrets will be stored or referenced." required:""`

	Overwrite bool `help:"Overwrite existing entries in the Vault and a template file." default:"false"`

	Interactive bool `help:"Interactive mode." default:"false"`

	Namespace string `help:"Namespace of the Kubernetes Secret." default:"default"`
}

func Run() {
	var cli CLI
	ctx := kong.Parse(&cli,
		kong.Name("optruck"),
		kong.Description("A command-line tool for managing secrets between 1Password Vault and your applications, with support for template-based restoration."),
		kong.UsageOnError(),
	)

	// Dispatch the command
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
