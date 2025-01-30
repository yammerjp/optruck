package optruck

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
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

	if cli.Interactive {
		err := cli.SetConfigInteractively()
		if err != nil {
			ctx.Fatalf("%v", err)
		}
	}
	if err := cli.SetDefaultIfEmpty(); err != nil {
		ctx.Fatalf("%v", err)
	}

	action, err := cli.Build()
	if err != nil {
		ctx.Fatalf("%v", err)
	}

	if err := action.Run(); err != nil {
		ctx.Fatalf("%v", err)
	}
}
