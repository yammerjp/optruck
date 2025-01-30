package optruck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/kube"
	"github.com/yammerjp/optruck/pkg/op"
)

// TODO: test

func (i InteractiveFlag) BeforeApply(ctx *kong.Context) error {
	cli := CLI{}
	if err := cli.SetOptionsInteractively(); err != nil {
		return err
	}
	// confirm
	// set options
	return nil
}

func (cli *CLI) SetOptionsInteractively() error {
	if err := cli.setDataSourceInteractively(); err != nil {
		return err
	}
	if err := cli.setTargetAccountInteractively(); err != nil {
		return err
	}
	if err := cli.setTargetVaultInteractively(); err != nil {
		return err
	}
	if err := cli.setTargetItemInteractively(); err != nil {
		return err
	}
	if err := cli.setDestInteractively(); err != nil {
		return err
	}
	return nil
}

func (cli *CLI) setDataSourceInteractively() error {
	if cli.EnvFile != "" || cli.K8sSecret != "" {
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
		if err := cli.setEnvFilePathInteractively(); err != nil {
			return err
		}
	case "k8s secret":
		if err := cli.setK8sSecretInteractively(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid data source: %s", result)
	}
	return nil
}

func (cli *CLI) setEnvFilePathInteractively() error {
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
	cli.EnvFile = result
	return nil
}

func (cli *CLI) setK8sSecretInteractively() error {
	kubeClient := kube.NewClient()
	namespaces, err := kubeClient.GetNamespaces()
	if err != nil {
		return err
	}

	if cli.K8sNamespace == "" {
		prompt := promptui.Select{
			Label: "Select kubernetes namespace",
			Items: namespaces,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		cli.K8sNamespace = result
	}

	secrets, err := kubeClient.GetSecrets(cli.K8sNamespace)
	if err != nil {
		return err
	}
	prompt := promptui.Select{
		Label: fmt.Sprintf("Select kubernetes secret on namespace %s", cli.K8sNamespace),
		Items: secrets,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	cli.K8sSecret = result
	return nil
}

func (cli *CLI) setTargetAccountInteractively() error {
	if cli.Account != "" {
		// already set
		return nil
	}

	accounts, err := op.NewExecutableClient(nil).ListAccounts()
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
	cli.Account = result
	return nil
}

func (cli *CLI) setTargetVaultInteractively() error {
	if cli.Vault != "" {
		// already set
		return nil
	}

	vaults, err := op.NewAccountClient(cli.Account, nil).ListVaults()
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
	cli.Vault = result
	return nil
}

func (cli *CLI) setTargetItemInteractively() error {
	if !cli.Overwrite {
		prompt := promptui.Select{
			Label: "Select overwrite mode",
			Items: []string{"overwrite existing", "create new"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		cli.Overwrite = result == "overwrite existing"
	}

	if cli.Item != "" {
		// already set
		return nil
	}

	items, err := op.NewVaultClient(cli.Account, cli.Vault, nil).ListItems()
	if err != nil {
		return err
	}

	if cli.Overwrite {
		return cli.setItemBySelectExisting(items)
	}

	return cli.setItemByInput(items)
}

func (cli *CLI) setItemBySelectExisting(currentItems []op.SecretReference) error {
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
	cli.Item = currentItems[i].ItemID
	return nil
}

func (cli *CLI) defaultItemName() string {
	defaultName := ""
	// TODO: define default item name format
	if cli.EnvFile != "" {
		defaultName = fmt.Sprintf("dotenv_%s", filepath.Base(filepath.Dir(cli.EnvFile)))
	} else if cli.K8sSecret != "" {
		defaultName = fmt.Sprintf("kubernetes_secret_%s_%s", cli.K8sNamespace, cli.K8sSecret)
	}
	return defaultName
}

func (cli *CLI) setItemByInput(currentItems []op.SecretReference) error {
	prompt := promptui.Prompt{
		Label:   "Enter item name",
		Default: cli.defaultItemName(),
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
	cli.Item = result
	return nil
}

func validateOutputPath(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if stat.IsDir() {
		return errors.New("output path is already created as a directory")
	}
	return nil
}

func (cli *CLI) setDestInteractively() error {
	if cli.Output != "" {
		// already set
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Enter output path",
		Default:  defaultOutputPath(cli.K8sSecret != "", cli.Item),
		Validate: validateOutputPath,
	}
	result, err := prompt.Run()
	if err != nil {
		return err
	}
	cli.Output = result

	if cli.Overwrite {
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
	cli.Overwrite = true
	return nil
}
