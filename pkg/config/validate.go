package config

import "fmt"

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
	if err := b.validateActionSpecially(); err != nil {
		return err
	}
	if err := b.validateOvewriteSpecially(); err != nil {
		return err
	}
	if err := b.validateTargetSpecially(); err != nil {
		return err
	}
	if err := b.validateDataSourceSpecially(); err != nil {
		return err
	}
	if err := b.validateDestSpecially(); err != nil {
		return err
	}
	return nil
}

func (b *ConfigBuilder) validateMissing() error {
	if err := b.validateActionMissing(); err != nil {
		return err
	}
	if err := b.validateOverwriteMissing(); err != nil {
		return err
	}
	if err := b.validateTargetMissing(); err != nil {
		return err
	}
	if err := b.validateDataSourceMissing(); err != nil {
		return err
	}
	if err := b.validateDestMissing(); err != nil {
		return err
	}
	return nil
}

func (b *ConfigBuilder) validateActionSpecially() error {
	if b.isUpload && b.isTemplate {
		return fmt.Errorf("cannot use both --upload and --template")
	}
	if b.isUpload && b.isMirror {
		return fmt.Errorf("cannot use both --upload and --mirror")
	}
	if b.isTemplate && b.isMirror {
		return fmt.Errorf("cannot use both --template and --mirror")
	}
	return nil
}

func (b *ConfigBuilder) validateActionMissing() error {
	if !b.isUpload && !b.isTemplate && !b.isMirror {
		return fmt.Errorf("one of --upload, --template, or --mirror must be specified")
	}
	return nil
}

func (b *ConfigBuilder) validateTargetSpecially() error {
	if len(b.item) > 100 {
		return fmt.Errorf("item must be less than 100 characters")
	}
	return nil
	// TODO check the file does not exist --overwrite-target || --overwrite
	// TODO check the file exists --overwrite-target || --overwrite
}

func (b *ConfigBuilder) validateTargetMissing() error {
	if b.account == "" {
		return fmt.Errorf("account is required")
	}
	if b.vault == "" {
		return fmt.Errorf("vault is required")
	}
	if b.item == "" {
		return fmt.Errorf("item is required")
	}
	return nil
}

func (b *ConfigBuilder) validateOvewriteSpecially() error {
	if b.overwriteTarget && b.overwrite {
		return fmt.Errorf("cannot use both --overwrite-target and --overwrite")
	}
	if b.overwriteTemplate && b.overwrite {
		return fmt.Errorf("cannot use both --overwrite-template and --overwrite")
	}
	if b.isTemplate && b.overwriteTarget {
		return fmt.Errorf("cannot use --overwrite-target on template action")
	}
	if b.isUpload && b.overwriteTemplate {
		return fmt.Errorf("cannot use --overwrite-template on upload action")
	}
	return nil
}

func (b *ConfigBuilder) validateOverwriteMissing() error {
	// overwrite is not required
	return nil
}

func (b *ConfigBuilder) validateDataSourceSpecially() error {
	if b.isTemplate {
		// data source is not needed
		if b.outputFormat != "" && b.outputFormat != "k8s" {
			if b.k8sSecret != "" {
				return fmt.Errorf("cannot use --k8s-secret on template action with output format %s", b.outputFormat)
			}
			if b.k8sNamespace != "" {
				return fmt.Errorf("cannot use --k8s-namespace on template action with output format %s", b.outputFormat)
			}
		}
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
		if b.envFile != "" {
			return fmt.Errorf("cannot use --env-file on template action")
		}
		if b.outputFormat != "k8s" {
			if b.k8sSecret != "" {
				return fmt.Errorf("cannot use --k8s-secret on template action with output format %s", b.outputFormat)
			}
			if b.k8sNamespace != "" {
				return fmt.Errorf("cannot use --k8s-namespace on template action with output format %s", b.outputFormat)
			}
		}
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

func (b *ConfigBuilder) validateDestSpecially() error {
	if b.isUpload {
		// upload is not needed
		if b.outputFormat != "" {
			return fmt.Errorf("output format is not allowed when uploading")
		}
		if b.output != "" {
			return fmt.Errorf("output is not allowed when uploading")
		}
		return nil
	}
	if b.outputFormat != "k8s" && b.outputFormat != "env" && b.outputFormat != "" {
		return fmt.Errorf("invalid output format: %s", b.outputFormat)
	}
	// TODO check the file does not exist --overwrite-template || --overwrite
	// TODO check the file exists --overwrite-target
	return nil
}

func (b *ConfigBuilder) validateDestMissing() error {
	if b.isUpload {
		// upload is not needed
		return nil
	}
	if b.outputFormat != "k8s" && b.outputFormat != "env" {
		return fmt.Errorf("invalid output format: %s", b.outputFormat)
	}
	if b.outputFormat == "k8s" {
		if b.k8sSecret == "" {
			return fmt.Errorf("k8s secret is required")
		}
		if b.k8sNamespace == "" {
			return fmt.Errorf("k8s namespace is required")
		}
	}

	if b.output == "" {
		return fmt.Errorf("output is required")
	}
	return nil
}
