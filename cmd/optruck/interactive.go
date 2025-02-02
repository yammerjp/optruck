package optruck

import (
	"fmt"
	"log/slog"

	"github.com/yammerjp/optruck/internal/util/interactive"
)

func (cli *CLI) SetOptionsInteractively() error {
	if err := cli.setDataSourceInteractively(); err != nil {
		return err
	}
	if err := cli.setTargetInteractively(); err != nil {
		return err
	}
	if err := cli.setDestInteractively(); err != nil {
		return err
	}
	return nil
}

func (cli *CLI) setDataSourceInteractively() error {
	if cli.EnvFile != "" || cli.K8sSecret != "" {
		slog.Debug("data source already set", "envFile", cli.EnvFile, "k8sSecret", cli.K8sSecret)
		// already set
		return nil
	}
	ds, err := cli.runner.SelectDataSource()
	if err != nil {
		return err
	}
	switch ds {
	case interactive.DataSourceEnvFile:
		slog.Debug("setting env file path")
		envFilePath, err := cli.runner.PromptEnvFilePath()
		if err != nil {
			return err
		}
		cli.EnvFile = envFilePath
	case interactive.DataSourceK8sSecret:
		slog.Debug("setting k8s secret")
		if cli.K8sNamespace == "" {
			namespace, err := cli.runner.SelectKubeNamespace()
			if err != nil {
				return err
			}
			cli.K8sNamespace = namespace
		}
		if cli.K8sSecret == "" {
			secret, err := cli.runner.SelectKubeSecret(cli.K8sNamespace)
			if err != nil {
				return err
			}
			cli.K8sSecret = secret
		}
	default:
		return fmt.Errorf("invalid data source: %s", ds)
	}
	return nil
}

func (cli *CLI) setTargetInteractively() error {
	if cli.Account == "" {
		account, err := cli.runner.SelectOpAccount()
		if err != nil {
			return err
		}
		cli.Account = account
	}

	if cli.Vault == "" {
		vault, err := cli.runner.SelectOpVault(cli.Account)
		if err != nil {
			return err
		}
		cli.Vault = vault
	}
	if cli.Item == "" {
		if cli.Overwrite {
			overwrite, err := cli.runner.SelectOpItemOverwriteOrNot()
			if err != nil {
				return err
			}
			cli.Overwrite = overwrite
		} else {
			itemName, err := cli.runner.PromptOpItemName(cli.Account, cli.Vault)
			if err != nil {
				return err
			}
			cli.Item = itemName
		}
	}
	return nil
}

func (cli *CLI) setDestInteractively() error {
	if cli.Output != "" {
		// already set
		return nil
	}
	outputPath, err := cli.runner.PromptOutputPath(cli.K8sSecret)
	if err != nil {
		return err
	}
	cli.Output = outputPath
	return nil
}
