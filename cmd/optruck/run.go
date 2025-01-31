package optruck

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/actions"
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

func (cli *CLI) buildOrBuildWithInteractive() (actions.Action, error) {
	if cli.Interactive {
		if err := cli.SetOptionsInteractively(); err != nil {
			return nil, err
		}
		cmds, err := cli.buildResultCommand()
		if err != nil {
			return nil, err
		}
		defer func() {
			if err == nil {
				err = cli.confirmToProceed(cmds)
			}
		}()
	}
	return cli.buildWithDefault()
}

func (cli *CLI) Run() error {
	action, err := cli.buildOrBuildWithInteractive()
	if err != nil {
		return err
	}
	return action.Run()
}

func (cli *CLI) confirmToProceed(cmds []string) error {
	fmt.Printf("The selected options are same as below.\n    $ %s\n", strings.Join(cmds, " "))
	fmt.Println("Do you want to proceed? (y/n)")
	prompt := promptui.Select{
		Label: "Proceed?",
		Items: []string{"y", "n"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	if result == "n" {
		return fmt.Errorf("aborted")
	}
	return nil
}

func (cli *CLI) buildResultCommand() ([]string, error) {
	cmds := []string{"optruck", cli.Item}
	if cli.Vault != "" {
		cmds = append(cmds, "--vault", cli.Vault)
	}
	if cli.Account != "" {
		cmds = append(cmds, "--account", cli.Account)
	}
	if cli.EnvFile != "" {
		cmds = append(cmds, "--env-file", cli.EnvFile)
	} else if cli.K8sSecret != "" {
		cmds = append(cmds, "--k8s-secret", cli.K8sSecret)
		if cli.K8sNamespace != "default" {
			cmds = append(cmds, "--k8s-namespace", cli.K8sNamespace)
		}
	}
	if cli.Output != "" {
		cmds = append(cmds, "--output", cli.Output)
	}
	if cli.Overwrite {
		cmds = append(cmds, "--overwrite")
	}
	return cmds, nil
}
