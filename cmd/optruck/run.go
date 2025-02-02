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

	action, err := cli.build(confirmation)
	if err != nil {
		return err
	}

	return action.Run()
}

func (cli *CLI) buildResultCommand() ([]string, error) {
	cmds := []string{"optruck", cli.Item}
	// target options
	if cli.Overwrite {
		cmds = append(cmds, "--overwrite")
	}
	if cli.Account != "" {
		cmds = append(cmds, "--account", cli.Account)
	}
	if cli.Vault != "" {
		cmds = append(cmds, "--vault", cli.Vault)
	}

	// data source options
	if cli.EnvFile != "" {
		cmds = append(cmds, "--env-file", cli.EnvFile)
	} else if cli.K8sSecret != "" {
		cmds = append(cmds, "--k8s-secret", cli.K8sSecret)
		if cli.K8sNamespace != "default" {
			cmds = append(cmds, "--k8s-namespace", cli.K8sNamespace)
		}
	}

	// output options
	if cli.Output != "" {
		cmds = append(cmds, "--output", cli.Output)
	}
	return cmds, nil
}
