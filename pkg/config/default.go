package config

import "fmt"

const defaultEnvFilePath = ".env"
const defaultOutputFormat = "env"
const defaultOutputPathOnEnv = ".env.1password"

func defaultOutputPathOnK8s(item string) string {
	return fmt.Sprintf("%s-secret.yaml.1password", item)
}

func (b *ConfigBuilder) SetDefaultIfEmpty() error {
	if !b.isUpload && !b.isTemplate && !b.isMirror {
		b.isMirror = true
	}
	if b.envFile == "" && b.k8sSecret == "" {
		b.envFile = defaultEnvFilePath
	}
	if b.k8sSecret != "" && b.k8sNamespace == "" {
		return fmt.Errorf("k8s namespace is required when using k8s secret")
	}
	if b.outputFormat == "" {
		b.outputFormat = defaultOutputFormat
	}
	if b.output == "" {
		if b.outputFormat == "env" {
			b.output = defaultOutputPathOnEnv
		} else {
			b.output = defaultOutputPathOnK8s(b.item)
		}
	}
	if b.account == "" {
		// FIXME: not use op.Client directly in plg/config
		// if account is not specified and exist only one account, use it
		opClient := b.buildOpTarget().BuildClient()
		accounts, err := opClient.ListAccounts()
		if err != nil {
			return fmt.Errorf("failed to list accounts: %w", err)
		}
		if len(accounts) == 1 {
			b.account = accounts[0].URL
		}
	}
	if b.vault == "" {
		// FIXME: not use op.Client directly in plg/config
		// if vault is not specified and exist only one vault, use it
		opClient := b.buildOpTarget().BuildClient()
		vaults, err := opClient.ListVaults()
		if err != nil {
			return fmt.Errorf("failed to list vaults: %w", err)
		}
		if len(vaults) == 1 {
			b.vault = vaults[0].Name
		}
	}
	return nil
}
