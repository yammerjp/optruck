package optruck

import (
	"fmt"
)

// TODO: test

func (cli *CLI) validateCommon() error {
	if err := cli.validateSpecially(); err != nil {
		return err
	}
	if err := cli.validateMissing(); err != nil {
		return err
	}
	return nil
}

func (cli *CLI) validateSpecially() error {
	if len(cli.Item) > 100 {
		return fmt.Errorf("item must be less than 100 characters")
	}
	if cli.EnvFile != "" && (cli.K8sSecret != "" || cli.K8sNamespace != "") {
		return fmt.Errorf("cannot use both --env-file and --k8s-secret or --k8s-namespace")
	}

	return nil
}

func (cli *CLI) validateMissing() error {
	if cli.EnvFile == "" && cli.K8sSecret == "" {
		return fmt.Errorf("either --env-file or --k8s-secret is required")
	}
	if cli.K8sSecret != "" && cli.K8sNamespace == "" {
		return fmt.Errorf("k8s namespace is required")
	}

	if cli.Account == "" {
		return fmt.Errorf("account is required")
	}
	if cli.Vault == "" {
		return fmt.Errorf("vault is required")
	}
	if cli.Item == "" {
		return fmt.Errorf("item is required")
	}

	if cli.Output == "" {
		return fmt.Errorf("output is required")
	}
	return nil
}
