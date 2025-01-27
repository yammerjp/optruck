package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/kube"
)

// TODO: test
func (b *ConfigBuilder) SetConfigInteractively() error {
	if err := b.validateSpecially(); err != nil {
		// Validate early before entering interactive mode, even though we'll check again later
		return err
	}

	if err := b.setActionInteractively(); err != nil {
		return err
	}
	fmt.Println("action", b.isUpload, b.isTemplate, b.isMirror)
	if err := b.setDataSourceInteractively(); err != nil {
		return err
	}
	if err := b.setTargetInteractively(); err != nil {
		return err
	}
	if err := b.setDestInteractively(); err != nil {
		return err
	}
	if err := b.validateSpecially(); err != nil {
		return err
	}
	cmds, err := b.buildResultCommand()
	if err != nil {
		return err
	}
	if err := b.SetDefaultIfEmpty(); err != nil {
		return err
	}
	if err := b.validateCommon(); err != nil {
		return err
	}
	fmt.Printf("The selected options are same as below.\n    $ %s\n", strings.Join(cmds, " "))
	fmt.Println("Do you want to proceed? (y/n)")
	prompt := promptui.Select{
		Label: "Proceed?",
		Items: []string{"y", "n"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	if result == "n" {
		return fmt.Errorf("aborted")
	}
	return nil
}

func (b *ConfigBuilder) buildResultCommand() ([]string, error) {
	cmds := []string{"optruck", b.item}
	if b.isUpload {
		cmds = append(cmds, "--upload")
	}
	if b.isTemplate {
		cmds = append(cmds, "--template")
	}
	if b.isMirror {
		cmds = append(cmds, "--mirror")
	}
	if b.envFile != "" {
		cmds = append(cmds, "--env-file", b.envFile)
	}
	if b.output != "" {
		cmds = append(cmds, "--output", b.output)
	}
	if b.vault != "" {
		cmds = append(cmds, "--vault", b.vault)
	}
	if b.account != "" {
		cmds = append(cmds, "--account", b.account)
	}
	if b.overwrite {
		cmds = append(cmds, "--overwrite")
	} else {
		if b.overwriteTarget {
			cmds = append(cmds, "--overwrite-target")
		}
		if b.overwriteTemplate {
			cmds = append(cmds, "--overwrite-template")
		}
	}
	return cmds, nil
}

func (b *ConfigBuilder) setActionInteractively() error {
	if b.isUpload || b.isTemplate || b.isMirror {
		// already set
		return nil
	}

	prompt := promptui.Select{
		Label:     "Select action",
		Items:     []string{"upload", "template", "mirror"},
		CursorPos: 2,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	b.isUpload = result == "upload"
	b.isTemplate = result == "template"
	b.isMirror = result == "mirror"
	return nil
}

func (b *ConfigBuilder) setDataSourceInteractively() error {
	if b.isTemplate {
		// data source is not needed
		return nil
	}
	if b.envFile != "" || b.k8sSecret != "" {
		// already set
		return nil
	}
	prompt := promptui.Select{
		Label:     "Select data source",
		Items:     []string{"env file", "k8s secret"},
		CursorPos: 2,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	switch result {
	case "env file":
		// TODO: create new function setEnvFilePathInteractively()
		prompt := promptui.Prompt{
			Label:   "Enter env file path",
			Default: defaultEnvFilePath,
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("env file path is required")
				}
				if _, err := os.Stat(input); err != nil {
					return fmt.Errorf("env file does not exist")
				}
				return nil
			},
		}
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.envFile = result
	case "k8s secret":
		// TODO: create new function setK8sSecretInteractively()
		kubeClient := kube.NewClient()
		namespaces, err := kubeClient.GetNamespaces()
		if err != nil {
			return err
		}

		if b.k8sNamespace == "" {
			prompt := promptui.Select{
				Label: "Select kubernetes namespace",
				Items: namespaces,
			}
			_, result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.k8sNamespace = result
		}

		secrets, err := kubeClient.GetSecrets(b.k8sNamespace)
		if err != nil {
			return err
		}
		prompt = promptui.Select{
			Label: fmt.Sprintf("Select kubernetes secret on namespace %s", b.k8sNamespace),
			Items: secrets,
		}
		_, result, err = prompt.Run()
		if err != nil {
			return err
		}
		b.k8sSecret = result
	default:
		return fmt.Errorf("invalid data source: %s", result)
	}
	return nil
}

func (b *ConfigBuilder) setTargetInteractively() error {

	// TODO: create new function setAccountInteractively()
	if b.account == "" {
		opClient := b.buildOpTarget().BuildClient()
		accounts, err := opClient.ListAccounts()
		if err != nil {
			return err
		}
		accountNames := make([]string, len(accounts))
		for i, account := range accounts {
			accountNames[i] = account.URL
		}
		prompt := promptui.Select{
			Label: "Select account",
			Items: accountNames,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.account = result
	}

	// TODO: create new function setVaultInteractively()
	if b.vault == "" {
		// regenerate opClient with selected account
		opClient := b.buildOpTarget().BuildClient()
		vaults, err := opClient.ListVaults()
		if err != nil {
			return err
		}
		vaultNames := make([]string, len(vaults))
		for i, vault := range vaults {
			vaultNames[i] = vault.Name
		}
		prompt := promptui.Select{
			Label: "Select vault",
			Items: vaultNames,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.vault = result
	}

	// TODO: create new function setItemInteractively()
	if !b.overwriteTarget && !b.overwrite {
		prompt := promptui.Select{
			Label: "Select overwrite mode",
			Items: []string{"overwrite existing", "create new"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		switch result {
		case "overwrite existing":
			b.overwriteTarget = true
		case "create new":
			b.overwriteTarget = false
		default:
			return fmt.Errorf("invalid selection: %s", result)
		}
	}

	if b.item == "" {
		// regenerate opClient with selected account and vault
		opClient := b.buildOpTarget().BuildClient()
		items, err := opClient.ListItems()
		if err != nil {
			return err
		}
		itemNames := make([]string, len(items))
		for i, item := range items {
			itemNames[i] = item.ItemName
		}
		if b.overwrite || b.overwriteTarget {
			// TODO: return item id if item name is duplicated
			prompt := promptui.Select{
				Label: "Select item name",
				Items: itemNames,
			}
			_, result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.item = result
		} else {
			defaultName := ""
			// TODO: define default item name format
			if b.envFile != "" {
				defaultName = fmt.Sprintf("dotenv_%s", filepath.Base(filepath.Dir(b.envFile)))
			} else if b.k8sSecret != "" {
				defaultName = fmt.Sprintf("kubernetes_secret_%s_%s", b.k8sNamespace, b.k8sSecret)
			}
			prompt := promptui.Prompt{
				Label:   "Enter item name",
				Default: defaultName,
				Validate: func(input string) error {
					// TODO: define item name format
					if input == "" {
						return fmt.Errorf("item name is required")
					}
					if len(input) > 100 {
						return fmt.Errorf("item name must be less than 100 characters")
					}
					if strings.Contains(input, " ") {
						return fmt.Errorf("item name must not contain spaces")
					}
					for _, n := range itemNames {
						if n == input {
							return fmt.Errorf("item name must be unique")
						}
					}
					return nil
				},
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.item = result
		}
	}
	return nil
}

func (b *ConfigBuilder) setDestInteractively() error {
	if b.isUpload {
		// upload action does not need dest
		return nil
	}

	// TODO: create new function setOutputFormatInteractively()
	// set format if not set
	if b.isTemplate && b.outputFormat == "" {
		// if upload or mirror, detect format from data source
		prompt := promptui.Select{
			Label: "Select output format",
			Items: []string{"env", "k8s"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.outputFormat = result
	}

	// TODO: create new function setOutputPathInteractively()
	// FIXME: Having user select overwrite mode first is not user-friendly and should be changed
	// set overwrite if not set
	if !b.overwrite && !b.overwriteTemplate {
		prompt := promptui.Select{
			Label: "Select template overwrite mode",
			Items: []string{"overwrite existing", "create new"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		switch result {
		case "overwrite existing":
			b.overwriteTemplate = true
		case "create new":
			b.overwriteTemplate = false
		default:
			return fmt.Errorf("invalid selection: %s", result)
		}
	}

	defaultOutputPath := ""
	if b.outputFormat == "env" {
		defaultOutputPath = defaultOutputPathOnEnv
	} else {
		defaultOutputPath = defaultOutputPathOnK8s(b.item)
	}

	// set output path if not set
	if b.output == "" {
		if b.overwrite || b.overwriteTemplate {
			// validate existing file
			prompt := promptui.Prompt{
				Label:   "Enter output path",
				Default: defaultOutputPath,
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("output path is required")
					}
					if _, err := os.Stat(input); err != nil {
						return fmt.Errorf("output path does not exist")
					}
					return nil
				},
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.output = result
		} else {
			// validate new file
			prompt := promptui.Prompt{
				Label:   "Enter output path",
				Default: defaultOutputPath,
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("output path is required")
					}
					if _, err := os.Stat(input); err == nil {
						return fmt.Errorf("output path already exists")
					}
					return nil
				},
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.output = result
		}
	}
	return nil
}
