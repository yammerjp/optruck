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

func (cli *CLI) buildOrBuildWithInteractive() (actions actions.Action, err error) {
	if cli.Interactive {
		cli.runner = &InteractiveRunnerImpl{}
		if err = cli.SetOptionsInteractively(); err != nil {
			return nil, err
		}
		var cmds []string
		cmds, err = cli.buildResultCommand()
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
	fmt.Println("The selected options are same as below.")
	fmt.Print("    $ optruck")
	for _, cmd := range cmds {
		if strings.HasPrefix(cmd, "--") {
			// break line
			fmt.Printf(" \\\n      %s", cmd)
		} else {
			fmt.Printf(" %s", cmd)
		}
	}
	fmt.Println()

	i, _, err := cli.runner.Select(promptui.Select{
		Label:     "Do you want to proceed? (yes/no)",
		Items:     []string{"yes", "no"},
		Templates: selectTemplateBuilder("Do you want to proceed?", "", ""),
	})
	if err != nil {
		return err
	}
	if i != 0 {
		return fmt.Errorf("aborted")
	}
	return nil
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
