package optruck

import (
	"fmt"
	"log/slog"

	"github.com/yammerjp/optruck/internal/interactive"
)

func (cli *CLI) SetOptionsInteractively(runner interactive.Runner) error {
	if err := cli.setDataSourceInteractively(runner); err != nil {
		return err
	}
	if err := cli.setTargetInteractively(runner); err != nil {
		return err
	}
	if err := cli.setDestInteractively(runner); err != nil {
		return err
	}
	return nil
}

func (cli *CLI) setDataSourceInteractively(runner interactive.Runner) error {
	if cli.EnvFile != "" || cli.K8sSecret != "" {
		slog.Debug("data source already set", "envFile", cli.EnvFile, "k8sSecret", cli.K8sSecret)
		// already set
		return nil
	}
	ds, err := runner.SelectDataSource()
	if err != nil {
		return fmt.Errorf("failed to select data source: %w. Please select a valid data source and try again.", err)
	}
	switch ds {
	case interactive.DataSourceEnvFile:
		slog.Debug("setting env file path")
		envFilePath, err := runner.PromptEnvFilePath()
		if err != nil {
			return fmt.Errorf("failed to set env file path: %w. Please provide a valid path and try again.", err)
		}
		cli.EnvFile = envFilePath
	case interactive.DataSourceK8sSecret:
		slog.Debug("setting k8s secret")
		if cli.K8sNamespace == "" {
			namespace, err := runner.SelectKubeNamespace()
			if err != nil {
				return fmt.Errorf("failed to select Kubernetes namespace: %w. Please select a valid namespace and try again.", err)
			}
			cli.K8sNamespace = namespace
		}
		if cli.K8sSecret == "" {
			secret, err := runner.SelectKubeSecret(cli.K8sNamespace)
			if err != nil {
				return fmt.Errorf("failed to select Kubernetes secret: %w. Please select a valid secret and try again.", err)
			}
			cli.K8sSecret = secret
		}
	default:
		return fmt.Errorf("invalid data source: %s, select .env file or kubernetes secret", ds)
	}
	return nil
}

func (cli *CLI) setTargetInteractively(runner interactive.Runner) error {
	if cli.Account == "" {
		account, err := runner.SelectOpAccount()
		if err != nil {
			return fmt.Errorf("failed to select 1Password account: %w. Please select a valid account and try again.", err)
		}
		cli.Account = account
	}

	if cli.Vault == "" {
		vault, err := runner.SelectOpVault(cli.Account)
		if err != nil {
			return fmt.Errorf("failed to select 1Password vault: %w. Please select a valid vault and try again.", err)
		}
		cli.Vault = vault
	}
	if cli.Item == "" {
		if !cli.Overwrite {
			overwrite, err := runner.SelectOpItemOverwriteOrNot()
			if err != nil {
				return fmt.Errorf("failed to select overwrite mode: %w. Please select a valid mode and try again.", err)
			}
			cli.Overwrite = overwrite
		}
		if cli.Overwrite {
			itemName, err := runner.SelectOpItemName(cli.Account, cli.Vault)
			if err != nil {
				return fmt.Errorf("failed to select 1Password item name: %w. Please select a valid item name and try again.", err)
			}
			cli.Item = itemName
		} else {
			itemName, err := runner.PromptOpItemName(cli.Account, cli.Vault, cli.K8sSecret)
			if err != nil {
				return fmt.Errorf("failed to prompt 1Password item name: %w. Please provide a valid item name and try again.", err)
			}
			cli.Item = itemName
		}
	}
	return nil
}

func (cli *CLI) setDestInteractively(runner interactive.Runner) error {
	if cli.Output != "" {
		// already set
		return nil
	}
	outputPath, err := runner.PromptOutputPath(cli.K8sSecret)
	if err != nil {
		return fmt.Errorf("failed to prompt output path: %w. Please provide a valid path and try again.", err)
	}
	cli.Output = outputPath
	return nil
}
