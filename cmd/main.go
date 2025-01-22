package cmd

import (
	"fmt"

	"github.com/alecthomas/kong"
)

const Version = "0.1.0"

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

	Interactive bool `help:"Interactive mode." default:"false"`

	Namespace string `help:"Namespace of the Kubernetes Secret." default:"default"`
}

// ValidateInputSource ensures that exactly one input source is specified
func (f *UploadFlags) ValidateInputSource() error {
	hasEnvFile := f.EnvFile != ""
	hasKubeSecret := f.KubeSecret != ""

	if !hasEnvFile && !hasKubeSecret {
		return fmt.Errorf("either --env-file or --kube-secret must be specified")
	}
	if hasEnvFile && hasKubeSecret {
		return fmt.Errorf("cannot specify both --env-file and --kube-secret")
	}
	return nil
}

type UploadFlags struct {
	EnvFile    string `help:"Path to the .env file containing secrets." type:"existingfile" short:"e" default:".env"`
	KubeSecret string `help:"Name of the Kubernetes Secret to process." short:"k"`
	Context    string `help:"Context of the Kubernetes Secret." default:"default"`
}

// TemplateFlags defines flags specific to the template command
type TemplateFlags struct {
	Output    string `help:"Path to save the generated template file. default value is .env.1password.tpl or <name-secret>.yaml.1password.tpl"`
	Overwrite bool   `help:"Overwrite existing entries in the Vault and a template file." default:"false"`
}

// UploadCmd represents the `upload` command.
type UploadCmd struct {
	UploadFlags
	SharedFlags
}

// TemplateCmd represents the `template` command.
type TemplateCmd struct {
	TemplateFlags
	SharedFlags
}

// MirrorCmd represents the `mirror` command.
type MirrorCmd struct {
	UploadFlags
	TemplateFlags
	SharedFlags
}

// VersionCmd represents the `version` command.
type VersionCmd struct{}

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

// Run executes the logic for the `upload` command.
func (cmd *UploadCmd) Run() error {
	if err := cmd.ValidateInputSource(); err != nil {
		return err
	}
	// kubernetes option is not supported yet
	if cmd.KubeSecret != "" {
		return fmt.Errorf("kubernetes option is not supported yet")
	}
	if cmd.Context != "default" {
		return fmt.Errorf("context option is not supported yet")
	}
	if cmd.Namespace != "default" {
		return fmt.Errorf("namespace option is not supported yet")
	}

	// interactive mode is not supported yet
	if cmd.Interactive {
		return fmt.Errorf("interactive mode is not supported yet")
	}

	fmt.Println("Executing 'upload' command...")
	fmt.Printf("EnvFile: %s\n", cmd.EnvFile)
	fmt.Printf("Vault: %s\n", cmd.Vault)
	fmt.Printf("Account: %s\n", cmd.Account)
	fmt.Printf("Item: %s\n", cmd.Item)
	return nil
}

// Run executes the logic for the `template` command.
func (cmd *TemplateCmd) Run() error {
	if cmd.Namespace != "default" {
		return fmt.Errorf("namespace option is not supported yet")
	}

	fmt.Println("Executing 'template' command...")

	fmt.Printf("Output: %s\n", cmd.Output)
	fmt.Printf("Account: %s\n", cmd.Account)
	fmt.Printf("Item: %s\n", cmd.Item)
	fmt.Printf("Overwrite: %t\n", cmd.Overwrite)
	return nil
}

// Run executes the logic for the `mirror` command.
func (cmd *MirrorCmd) Run() error {
	if err := cmd.ValidateInputSource(); err != nil {
		return err
	}

	// kubernetes option is not supported yet
	if cmd.KubeSecret != "" {
		return fmt.Errorf("kubernetes option is not supported yet")
	}
	if cmd.Context != "default" {
		return fmt.Errorf("context option is not supported yet")
	}
	if cmd.Namespace != "default" {
		return fmt.Errorf("namespace option is not supported yet")
	}

	// interactive mode is not supported yet
	if cmd.Interactive {
		return fmt.Errorf("interactive mode is not supported yet")
	}

	fmt.Printf("EnvFile: %s\n", cmd.EnvFile)

	fmt.Printf("Vault: %s\n", cmd.Vault)
	fmt.Printf("Account: %s\n", cmd.Account)
	fmt.Printf("Item: %s\n", cmd.Item)
	fmt.Printf("Overwrite: %t\n", cmd.Overwrite)
	fmt.Printf("Output: %s\n", cmd.Output)
	return nil
}

func (cmd *VersionCmd) Run() error {
	fmt.Printf("optruck version %s\n", Version)
	return nil
}
