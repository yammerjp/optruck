package optruck

import (
	"github.com/yammerjp/optruck/internal/util/interactive"
	utilLogger "github.com/yammerjp/optruck/internal/util/logger"

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
	if err := cli.Run(); err != nil {
		ctx.Fatalf("%v", err)
	}
}

func (cli *CLI) Run() error {
	utilLogger.SetDefaultLogger(cli.LogLevel)

	var confirmation func() error

	if cli.Interactive {
		runner := *interactive.NewImplRunner()
		if err := cli.SetOptionsInteractively(runner); err != nil {
			return err
		}
		cmds, err := cli.buildResultCommand()
		if err != nil {
			return err
		}
		confirmation = func() error {
			return runner.Confirm(cmds)
		}
	} else {
		confirmation = func() error {
			// confirmed by default
			return nil
		}
	}

	action, err := cli.buildAction(confirmation)
	if err != nil {
		return err
	}

	return action.Run()
}
