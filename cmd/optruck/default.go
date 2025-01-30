package optruck

import (
	"fmt"

	"github.com/yammerjp/optruck/pkg/op"
)

const (
	defaultEnvFilePath = ".env"
)

func (cli *CLI) SetDefaultIfEmpty() error {
	if err := cli.setDefaultTargetIfNotSet(); err != nil {
		return err
	}

	if err := cli.setDefaultDataSourceIfNotSet(); err != nil {
		return err
	}

	if err := cli.setDefaultOutputIfNotSet(); err != nil {
		return err
	}

	return nil
}

func (cli *CLI) setDefaultTargetIfNotSet() error {
	if cli.Account == "" {
		// if account is not specified and exist only one account, use it
		accounts, err := op.NewExecutableClient(nil).ListAccounts()
		if err != nil {
			return fmt.Errorf("failed to list accounts: %w", err)
		}
		if len(accounts) == 1 {
			cli.Account = accounts[0].URL
		}
	}
	if cli.Vault == "" {
		// if vault is not specified and exist only one vault, use it
		vaults, err := op.NewAccountClient(cli.Account, nil).ListVaults()
		if err != nil {
			return fmt.Errorf("failed to list vaults: %w", err)
		}
		if len(vaults) == 1 {
			cli.Vault = vaults[0].Name
		}
	}
	return nil
}

func (cli *CLI) setDefaultDataSourceIfNotSet() error {
	if cli.EnvFile == "" && cli.K8sSecret == "" {
		cli.EnvFile = defaultEnvFilePath
	}
	return nil
}

func (cli *CLI) setDefaultOutputIfNotSet() error {
	if cli.Output == "" {
		cli.Output = defaultOutputPath(cli.K8sSecret != "", cli.Item)
	}
	return nil
}

func defaultOutputPath(isK8s bool, itemName string) string {
	if isK8s {
		return fmt.Sprintf("%s-secret.yaml.1password", itemName)
	}
	return ".env.1password"
}
