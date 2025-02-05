package optruck

import (
	"fmt"

	"github.com/yammerjp/optruck/internal/interactive"
	"github.com/yammerjp/optruck/pkg/actions"
	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/kube"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

func (cli *CLI) buildAction(confirmation func() error) (actions.Action, error) {
	ds, err := cli.buildDataSource()
	if err != nil {
		return nil, err
	}

	dest, err := cli.buildDest()
	if err != nil {
		return nil, err
	}

	opItemClient, err := cli.buildOpItemClient(true)
	if err != nil {
		return nil, err
	}

	return &actions.MirrorConfig{
		OpItemClient: *opItemClient,
		DataSource:   ds,
		Dest:         dest,
		Overwrite:    cli.Overwrite,
		Confirmation: confirmation,
	}, nil
}

func (cli *CLI) buildOpItemClient(strict bool) (*op.ItemClient, error) {
	if strict {
		if cli.Account == "" {
			accounts, err := op.NewExecutableClient().ListAccounts()
			if err != nil {
				return nil, fmt.Errorf("failed to list accounts: %w. Please check your 1Password configuration and try again.", err)
			}
			if len(accounts) != 1 {
				return nil, fmt.Errorf("multiple accounts found, please specify the account with --account option")
			}
			cli.Account = accounts[0].URL
		}
		if cli.Vault == "" {
			vaults, err := op.NewAccountClient(cli.Account).ListVaults()
			if err != nil {
				return nil, fmt.Errorf("failed to list vaults: %w. Please check your 1Password configuration and try again.", err)
			}
			if len(vaults) != 1 {
				return nil, fmt.Errorf("multiple vaults found, please specify the vault with --vault option")
			}
			cli.Vault = vaults[0].Name
		}
		if cli.Item == "" {
			return nil, fmt.Errorf("item name or ID is required, please specify the item with argument like `$ optruck <item>`")
		}
	}

	return op.NewItemClient(cli.Account, cli.Vault, cli.Item), nil
}

func (cli *CLI) buildDataSource() (datasources.Source, error) {
	if cli.K8sSecret != "" {
		if cli.K8sNamespace == "" {
			cli.K8sNamespace = interactive.DefaultKubernetesNamespace
		}
		return &datasources.K8sSecretSource{
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
			Client:     kube.NewClient(),
		}, nil
	}
	if cli.EnvFile == "" {
		cli.EnvFile = interactive.DefaultEnvFilePath
	}
	return &datasources.EnvFileSource{Path: cli.EnvFile}, nil
}

func (cli *CLI) buildDest() (output.Dest, error) {
	if cli.Output == "" {
		cli.Output = interactive.DefaultOutputPath(cli.K8sSecret)
	}
	if cli.K8sSecret != "" {
		return &output.K8sSecretTemplateDest{
			Path:       cli.Output,
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
		}, nil
	}

	return &output.EnvTemplateDest{
		Path: cli.Output,
	}, nil
}
