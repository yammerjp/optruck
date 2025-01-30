package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/kube"
	"github.com/yammerjp/optruck/pkg/op"
)

// TODO: test
func (b *ConfigBuilder) SetConfigInteractively() error {
	if err := b.validateSpecially(); err != nil {
		// Validate early before entering interactive mode, even though we'll check again later
		return err
	}

	if err := b.setDataSourceInteractively(); err != nil {
		return err
	}

	if err := b.setTargetAccountInteractively(); err != nil {
		return err
	}
	if err := b.setTargetVaultInteractively(); err != nil {
		return err
	}
	if err := b.setTargetItemInteractively(); err != nil {
		return err
	}

	if err := b.setDestInteractively(); err != nil {
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
	return b.confirmToProceed(cmds)
}

func (b *ConfigBuilder) confirmToProceed(cmds []string) error {
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
	}
	return cmds, nil
}

func (b *ConfigBuilder) setDataSourceInteractively() error {
	if b.envFile != "" || b.k8sSecret != "" {
		// already set
		return nil
	}
	prompt := promptui.Select{
		Label: "Select data source",
		Items: []string{"env file", "k8s secret"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	switch result {
	case "env file":
		if err := b.setEnvFilePathInteractively(); err != nil {
			return err
		}
	case "k8s secret":
		if err := b.setK8sSecretInteractively(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid data source: %s", result)
	}
	return nil
}

func (b *ConfigBuilder) setEnvFilePathInteractively() error {
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
	return nil
}

func (b *ConfigBuilder) setK8sSecretInteractively() error {
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
	prompt := promptui.Select{
		Label: fmt.Sprintf("Select kubernetes secret on namespace %s", b.k8sNamespace),
		Items: secrets,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	b.k8sSecret = result
	return nil
}

func (b *ConfigBuilder) setTargetAccountInteractively() error {
	if b.account != "" {
		// already set
		return nil
	}

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
	return nil
}

func (b *ConfigBuilder) setTargetVaultInteractively() error {
	if b.vault != "" {
		// already set
		return nil
	}

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
	return nil
}

func (b *ConfigBuilder) setTargetItemInteractively() error {
	if !b.overwrite {
		prompt := promptui.Select{
			Label: "Select overwrite mode",
			Items: []string{"overwrite existing", "create new"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.overwrite = result == "overwrite existing"
	}

	if b.item != "" {
		// already set
		return nil
	}

	opClient := b.buildOpTarget().BuildClient()
	items, err := opClient.ListItems()
	if err != nil {
		return err
	}

	if b.overwrite {
		return b.setItemBySelectExisting(items)
	}

	return b.setItemByInput(items)
}

func (b *ConfigBuilder) setItemBySelectExisting(currentItems []op.SecretReference) error {
	itemNames := make([]string, len(currentItems))
	for i, item := range currentItems {
		itemNames[i] = fmt.Sprintf("%s: %s", item.ItemID, item.ItemName)
	}
	prompt := promptui.Select{
		Label: "Select item name",
		Items: itemNames,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return err
	}
	b.item = currentItems[i].ItemID
	return nil
}

func (b *ConfigBuilder) defaultItemName() string {
	defaultName := ""
	// TODO: define default item name format
	if b.envFile != "" {
		defaultName = fmt.Sprintf("dotenv_%s", filepath.Base(filepath.Dir(b.envFile)))
	} else if b.k8sSecret != "" {
		defaultName = fmt.Sprintf("kubernetes_secret_%s_%s", b.k8sNamespace, b.k8sSecret)
	}
	return defaultName
}

func (b *ConfigBuilder) setItemByInput(currentItems []op.SecretReference) error {
	prompt := promptui.Prompt{
		Label:   "Enter item name",
		Default: b.defaultItemName(),
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("item name is required")
			}
			if len(input) > 100 {
				return fmt.Errorf("item name must be less than 100 characters")
			}
			if strings.Contains(input, " ") {
				return fmt.Errorf("item name must not contain spaces")
			}
			for _, n := range currentItems {
				if n.ItemName == input {
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
	return nil
}

func (b *ConfigBuilder) setDestInteractively() error {
	if b.envFile != "" {
		b.output = defaultOutputPathOnEnv
	} else if b.k8sSecret != "" {
		b.output = defaultOutputPathOnK8s(b.item)
	}
	return nil
}

func (b *ConfigBuilder) setOutputPathInteractively() error {
	if b.output != "" {
		// already set
		return nil
	}
	prompt := promptui.Prompt{
		Label:   "Enter output path",
		Default: b.output,
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("output path is required")
			}
			stat, err := os.Stat(input)
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			if stat.IsDir() {
				return fmt.Errorf("output path is already created as a directory")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		return err
	}
	b.output = result

	stat, err := os.Stat(result)
	if os.IsNotExist(err) {
		// if output path does not exist, not need to overwrite
		return nil
	}
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.New("output path is already created as a directory")
	}

	if b.overwrite {
		// allow overwrite
		return nil
	}

	confirmationPrompt := promptui.Select{
		Label: "Output path already exists. Do you want to overwrite?",
		Items: []string{"overwrite", "abort"},
	}
	_, result, err = confirmationPrompt.Run()
	if err != nil {
		return err
	}
	if result == "abort" {
		return fmt.Errorf("aborted, overwrite does not allowed")
	}
	b.overwrite = true
	return nil
}
