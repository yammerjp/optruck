package optruck

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/kube"
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

func selectTemplateBuilder(selectedPrefix string, mainField string, subField string) *promptui.SelectTemplates {
	active := fmt.Sprintf("▸ {{ .%s | cyan | underline }}", mainField)
	if subField != "" {
		active += fmt.Sprintf(` {{"("|faint}}{{ .%s | red | underline }}{{")"|faint}}`, subField)
	}

	inactive := fmt.Sprintf("  {{ .%s | cyan }}", mainField)
	if subField != "" {
		inactive += fmt.Sprintf(` {{"("|faint}}{{ .%s | red }}{{")"|faint}}`, subField)
	}

	selected := fmt.Sprintf(`{{ "✔" | green }} %-20s: {{ .%s }}`, selectedPrefix, mainField)
	if subField != "" {
		selected += fmt.Sprintf(` {{"("|faint}}{{ .%s }}{{")"|faint}}`, subField)
	}

	return &promptui.SelectTemplates{
		Label:    `{{ . | yellow }}`,
		Active:   active,
		Inactive: inactive,
		Selected: selected,
	}
}

func promptTemplateBuilder(successPrefix string, mainField string) *promptui.PromptTemplates {
	return &promptui.PromptTemplates{
		Prompt:  `{{ . | yellow }}`,
		Valid:   fmt.Sprintf(`{{ "✔" | green }} {{ .%s | yellow }}`, mainField),
		Invalid: fmt.Sprintf(`{{ "✘" | red }} {{ .%s | yellow }}`, mainField),
		Success: fmt.Sprintf(`{{ "✔" | green }} %-20s: `, successPrefix),
	}
}

type InteractiveRunner interface {
	Select(promptui.Select) (int, string, error)
	Input(promptui.Prompt) (string, error)
}

type InteractiveRunnerImpl struct{}

func (r *InteractiveRunnerImpl) Select(prompt promptui.Select) (int, string, error) {
	return prompt.Run()
}

func (r *InteractiveRunnerImpl) Input(prompt promptui.Prompt) (string, error) {
	return prompt.Run()
}

func (cli *CLI) setDataSourceInteractively() error {
	if cli.EnvFile != "" || cli.K8sSecret != "" {
		slog.Debug("data source already set", "envFile", cli.EnvFile, "k8sSecret", cli.K8sSecret)
		// already set
		return nil
	}
	_, result, err := cli.runner.Select(promptui.Select{
		Label:     "Select data source: ",
		Items:     []string{"env file", "k8s secret"},
		Templates: selectTemplateBuilder("Data Source", "", ""),
	})
	if err != nil {
		return err
	}
	slog.Debug("selected data source", "result", result)
	switch result {
	case "env file":
		slog.Debug("setting env file path")
		if err := cli.setEnvFilePathInteractively(); err != nil {
			return err
		}
	case "k8s secret":
		slog.Debug("setting k8s secret")
		if err := cli.setK8sSecretInteractively(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid data source: %s", result)
	}
	return nil
}

func (cli *CLI) setEnvFilePathInteractively() error {
	result, err := cli.runner.Input(promptui.Prompt{
		Label:   "Enter env file path: ",
		Default: defaultEnvFilePath,
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("env file path is required")
			}
			stat, err := os.Stat(input)
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			if stat.IsDir() {
				return fmt.Errorf("env file path is already created as a directory")
			}
			return nil
		},
		Templates: promptTemplateBuilder("Env File Path", ""),
	})
	if err != nil {
		return err
	}
	cli.EnvFile = result
	return nil
}

func (cli *CLI) setK8sSecretInteractively() error {
	kubeClient := kube.NewClient(cli.exec)
	namespaces, err := kubeClient.GetNamespaces()
	if err != nil {
		return err
	}

	if cli.K8sNamespace == "" {
		_, result, err := cli.runner.Select(promptui.Select{
			Label:     "Select Kubernetes Namespace: ",
			Items:     namespaces,
			Templates: selectTemplateBuilder("Kubernetes Namespace", "", ""),
		})
		if err != nil {
			return err
		}
		cli.K8sNamespace = result
	}

	secrets, err := kubeClient.GetSecrets(cli.K8sNamespace)
	if err != nil {
		return err
	}
	_, result, err := cli.runner.Select(promptui.Select{
		Label:     fmt.Sprintf("Select kubernetes secret on namespace %s", cli.K8sNamespace),
		Items:     secrets,
		Templates: selectTemplateBuilder("Kubernetes Secret", "", ""),
	})
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

	accounts, err := op.NewExecutableClient(cli.exec).ListAccounts()
	if err != nil {
		return err
	}
	if len(accounts) == 0 {
		return fmt.Errorf("no 1Password accounts found")
	}
	i, _, err := cli.runner.Select(promptui.Select{
		Label:     "Select 1Password account: ",
		Items:     accounts,
		Templates: selectTemplateBuilder("1Password Account", "Email", "URL"),
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

	vaults, err := op.NewAccountClient(cli.Account, cli.exec).ListVaults()
	if err != nil {
		return err
	}
	if len(vaults) == 0 {
		return fmt.Errorf("no vaults found in account %s", cli.Account)
	}
	i, _, err := cli.runner.Select(promptui.Select{
		Label:     "Select 1Password vault: ",
		Items:     vaults,
		Templates: selectTemplateBuilder("1Password Vault", "Name", "ID"),
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
			Templates: selectTemplateBuilder("Overwrite mode", "", ""),
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

	items, err := op.NewVaultClient(cli.Account, cli.Vault, cli.exec).ListItems()
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
		Templates: selectTemplateBuilder("1Password Item", "ItemName", "ItemID"),
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
		Templates: promptTemplateBuilder("1Password Item Name", ""),
	})
	if err != nil {
		return err
	}
	cli.Item = result
	return nil
}

func validateOutputPath(path string) error {
	if path == "" {
		return errors.New("output path is required")
	}
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
		Label:     "Enter output path: ",
		Validate:  validateOutputPath,
		Default:   defaultOutputPath(cli.K8sSecret),
		Templates: promptTemplateBuilder("Output Path", ""),
	})
	if err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(result); err == nil {
		_, overwrite, err := cli.runner.Select(promptui.Select{
			Label:     fmt.Sprintf("File %s already exists. Do you want to overwrite it?", result),
			Items:     []string{"overwrite", "cancel"},
			Templates: selectTemplateBuilder("Overwrite", "", ""),
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
