package config

import (
	"fmt"
)

// TODO: test

func (b *ConfigBuilder) validateCommon() error {
	if err := b.validateSpecially(); err != nil {
		return err
	}
	if err := b.validateMissing(); err != nil {
		return err
	}
	return nil
}

func (b *ConfigBuilder) validateSpecially() error {
	if len(b.item) > 100 {
		return fmt.Errorf("item must be less than 100 characters")
	}
	if b.envFile != "" && (b.k8sSecret != "" || b.k8sNamespace != "") {
		return fmt.Errorf("cannot use both --env-file and --k8s-secret or --k8s-namespace")
	}

	return nil
}

func (b *ConfigBuilder) validateMissing() error {
	if b.envFile == "" && b.k8sSecret == "" {
		return fmt.Errorf("either --env-file or --k8s-secret is required")
	}
	if b.k8sSecret != "" && b.k8sNamespace == "" {
		return fmt.Errorf("k8s namespace is required")
	}

	if b.account == "" {
		return fmt.Errorf("account is required")
	}
	if b.vault == "" {
		return fmt.Errorf("vault is required")
	}
	if b.item == "" {
		return fmt.Errorf("item is required")
	}

	if b.output == "" {
		return fmt.Errorf("output is required")
	}
	return nil
}
