package optruck

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/internal/util/interactiverunner"
	"github.com/yammerjp/optruck/pkg/op"
)

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
		slog.Debug("data source already set", "envFile", cli.EnvFile, "k8sSecret", cli.K8sSecret)
		// already set
		return nil
	}
	ds, err := interactiverunner.NewDataSourceSelector(cli.runner).Select()
	if err != nil {
		return err
	}
	switch ds {
	case interactiverunner.DataSourceEnvFile:
		slog.Debug("setting env file path")
		envFilePath, err := interactiverunner.NewEnvFilePrompter(cli.runner).Prompt()
		if err != nil {
			return err
		}
		cli.EnvFile = envFilePath
	case interactiverunner.DataSourceK8sSecret:
		slog.Debug("setting k8s secret")
		if cli.K8sNamespace == "" {
			namespace, err := interactiverunner.NewKubeNamespaceSelector(cli.runner).Select()
			if err != nil {
				return err
			}
			cli.K8sNamespace = namespace
		}
		if cli.K8sSecret == "" {
			secret, err := interactiverunner.NewKubeSecretSelector(cli.runner, cli.K8sNamespace).Select()
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

func (cli *CLI) setTargetAccountInteractively() error {
	if cli.Account != "" {
		// already set
		return nil
	}

	accounts, err := op.NewExecutableClient().ListAccounts()
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		return fmt.Errorf("no 1Password accounts found")
	}
	i, _, err := cli.runner.Select(promptui.Select{
		Label:     "Select 1Password account: ",
		Items:     accounts,
		Templates: interactiverunner.SelectTemplateBuilder("1Password Account", "Email", "URL"),
	})
	if err != nil {
		return err
	}
	cli.Account = accounts[i].URL
	return nil
}

func (cli *CLI) setTargetVaultInteractively() error {
	if cli.Vault != "" {
		// already set
		return nil
	}

	if cli.Account == "" {
		return fmt.Errorf("account must be set before selecting vault")
	}

	vaults, err := op.NewAccountClient(cli.Account).ListVaults()
	if err != nil {
		return err
	}
	if len(vaults) == 0 {
		return fmt.Errorf("no vaults found in account %s", cli.Account)
	}
	i, _, err := cli.runner.Select(promptui.Select{
		Label:     "Select 1Password vault: ",
		Items:     vaults,
		Templates: interactiverunner.SelectTemplateBuilder("1Password Vault", "Name", "ID"),
	})
	if err != nil {
		return err
	}
	cli.Vault = vaults[i].ID
	return nil
}

func (cli *CLI) setTargetItemInteractively() error {
	if !cli.Overwrite {
		_, result, err := cli.runner.Select(promptui.Select{
			Label:     "Select overwrite mode: ",
			Items:     []string{"overwrite existing", "create new"},
			Templates: interactiverunner.SelectTemplateBuilder("Overwrite mode", "", ""),
		})
		if err != nil {
			return err
		}
		cli.Overwrite = result == "overwrite existing"
	}

	if cli.Item != "" {
		// already set
		return nil
	}

	if cli.Account == "" {
		return fmt.Errorf("account must be set before selecting item")
	}

	if cli.Vault == "" {
		return fmt.Errorf("vault must be set before selecting item")
	}

	items, err := op.NewVaultClient(cli.Account, cli.Vault).ListItems()
	if err != nil {
		return err
	}

	if cli.Overwrite {
		if len(items) == 0 {
			return fmt.Errorf("no items found in vault %s", cli.Vault)
		}
		return cli.setItemBySelectExisting(items)
	}

	return cli.setItemByInput(items)
}

func (cli *CLI) setItemBySelectExisting(currentItems []op.SecretReference) error {
	i, _, err := cli.runner.Select(promptui.Select{
		Label:     "Select 1Password item name: ",
		Items:     currentItems,
		Templates: interactiverunner.SelectTemplateBuilder("1Password Item", "ItemName", "ItemID"),
	})
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
	result, err := cli.runner.Input(promptui.Prompt{
		Label:   "Enter 1Password item name: ",
		Default: cli.defaultItemName(),
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("item name is required")
			}
			if len(input) < 1 {
				return fmt.Errorf("item name must be at least 1 character")
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
		Templates: interactiverunner.PromptTemplateBuilder("1Password Item Name", ""),
	})
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

	result, err := cli.runner.Input(promptui.Prompt{
		Label: "Enter output path: ",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("output path is required")
			}
			if err := validateOutputPath(input); err != nil {
				return err
			}
			// Check if the directory exists
			dir := filepath.Dir(input)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return fmt.Errorf("directory %s does not exist", dir)
			}
			return nil
		},
		Templates: interactiverunner.PromptTemplateBuilder("Output Path", ""),
		Default:   defaultOutputPath(cli.K8sSecret),
	})
	if err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(result); err == nil {
		_, overwrite, err := cli.runner.Select(promptui.Select{
			Label:     fmt.Sprintf("File %s already exists. Do you want to overwrite it?", result),
			Items:     []string{"overwrite", "cancel"},
			Templates: interactiverunner.SelectTemplateBuilder("Overwrite", "", ""),
		})
		if err != nil {
			return err
		}
		if overwrite == "cancel" {
			return fmt.Errorf("cancelled by user")
		}
	}

	cli.Output = result
	return nil
}

func (cli *CLI) confirmToProceed(cmds []string) error {
	fmt.Println("The selected options are same as below.")
	fmt.Print("    $ optruck")
	for _, cmd := range cmds {
		if strings.HasPrefix(cmd, "--") {
			// break line
			fmt.Printf(" \\\n      %s", cmd)
		} else {
			fmt.Printf(" %s", cmd)
		}
	}
	fmt.Println()

	i, _, err := cli.runner.Select(promptui.Select{
		Label:     "Do you want to proceed? (yes/no)",
		Items:     []string{"yes", "no"},
		Templates: interactiverunner.SelectTemplateBuilder("Do you want to proceed?", "", ""),
	})
	if err != nil {
		return err
	}
	if i != 0 {
		return fmt.Errorf("aborted")
	}
	return nil
}
