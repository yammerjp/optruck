package optruck

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/yammerjp/optruck/pkg/config"
)

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
		WithOverwrite(cli.Overwrite).
		WithLogLevel(cli.LogLevel)

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

	action, err := builder.Build()
	if err != nil {
		ctx.Fatalf("%v", err)
	}

	if err := action.Run(); err != nil {
		ctx.Fatalf("%v", err)
	}
}
