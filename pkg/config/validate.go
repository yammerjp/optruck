package config

import "fmt"

func (b *ConfigBuilder) validateSpecially() error {
	if b.overwriteTarget && b.overwrite {
		return fmt.Errorf("cannot use both --overwrite-target and --overwrite")
	}
	if b.overwriteTemplate && b.overwrite {
		return fmt.Errorf("cannot use both --overwrite-template and --overwrite")
	}
	if b.overwriteTarget && b.isTemplate {
		return fmt.Errorf("cannot use --overwrite-target on template action")
	}
	if b.overwriteTemplate && b.isUpload {
		return fmt.Errorf("cannot use --overwrite-template on upload action")
	}

	if len(b.item) > 100 {
		return fmt.Errorf("item must be less than 100 characters")
	}

	if err := b.validateDataSourceSpecially(); err != nil {
		return err
	}
	return nil
}

func (b *ConfigBuilder) validateMissing() error {
	if b.item == "" {
		return fmt.Errorf("item is required")
	}

	if b.vault == "" {
		return fmt.Errorf("vault is required")
	}
	if b.account == "" {
		return fmt.Errorf("account is required")
	}
	if err := b.validateDataSourceMissing(); err != nil {
		return err
	}

	return nil
}

func (b *ConfigBuilder) validateCommon() error {
	if err := b.validateSpecially(); err != nil {
		return err
	}
	if err := b.validateMissing(); err != nil {
		return err
	}
	return nil
}

func (b *ConfigBuilder) validateDataSourceSpecially() error {
	if b.isTemplate {
		// data source is not needed
		return nil
	}
	if b.envFile != "" && b.k8sSecret != "" {
		return fmt.Errorf("cannot use both --env-file and --k8s-secret")
	}
	return nil
}

func (b *ConfigBuilder) validateDataSourceMissing() error {
	if b.isTemplate {
		// data source is not needed
		return nil
	}
	if b.envFile == "" && b.k8sSecret == "" {
		return fmt.Errorf("either --env-file or --k8s-secret is required")
	}
	if b.k8sSecret != "" && b.k8sNamespace == "" {
		return fmt.Errorf("k8s namespace is required")
	}
	return nil
}
