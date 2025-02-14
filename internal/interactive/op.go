package interactive

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/op"
)

func (r Runner) SelectOpAccount() (string, error) {
	accounts, err := op.NewExecutableClient().ListAccounts()
	if err != nil {
		return "", err
	}
	if len(accounts) == 0 {
		return "", fmt.Errorf("no 1Password accounts found")
	}
	i, _, err := r.Select(promptui.Select{
		Label:     "Select 1Password account: ",
		Items:     accounts,
		Templates: SelectTemplateBuilder("1Password Account", "Email", "URL"),
	})
	if err != nil {
		return "", err
	}
	return accounts[i].URL, nil
}

func (r Runner) SelectOpVault(account string) (string, error) {
	vaults, err := op.NewAccountClient(account).ListVaults()
	if err != nil {
		return "", err
	}
	if len(vaults) == 0 {
		return "", fmt.Errorf("no vaults found in account %s, please create a 1Password vault first", account)
	}
	i, _, err := r.Select(promptui.Select{
		Label:     "Select 1Password vault: ",
		Items:     vaults,
		Templates: SelectTemplateBuilder("1Password Vault", "Name", "ID"),
	})
	if err != nil {
		return "", err
	}
	return vaults[i].ID, nil
}

func (r Runner) SelectOpItemOverwriteOrNot() (bool, error) {
	_, result, err := r.Select(promptui.Select{
		Label:     "Select overwrite mode: ",
		Items:     []string{"overwrite existing", "create new"},
		Templates: SelectTemplateBuilder("Overwrite mode", "", ""),
	})
	if err != nil {
		return false, err
	}
	return result == "overwrite existing", nil
}

func (r Runner) SelectOpItemName(account, vault string) (string, error) {
	items, err := op.NewVaultClient(account, vault).ListItems()
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", fmt.Errorf("no items found in vault %s", vault)
	}

	i, _, err := r.Select(promptui.Select{
		Label:     "Select 1Password item name: ",
		Items:     items,
		Templates: SelectTemplateBuilder("1Password Item", "ItemName", "ItemID"),
	})
	if err != nil {
		return "", err
	}
	return items[i].ItemID, nil
}

func (r Runner) PromptOpItemName(account, vault, k8sSecret string) (string, error) {
	items, err := op.NewVaultClient(account, vault).ListItems()
	if err != nil {
		return "", err
	}

	defaultItemName, err := defaultItemName(k8sSecret)
	if err != nil {
		return "", err
	}

	result, err := r.Input(promptui.Prompt{
		Label:   "Enter 1Password item name: ",
		Default: defaultItemName,
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
			for _, n := range items {
				if n.ItemName == input {
					return fmt.Errorf("item name must be unique")
				}
			}
			return nil
		},
		Templates: PromptTemplateBuilder("1Password Item Name", ""),
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

func defaultItemName(k8sSecret string) (string, error) {
	// current directory name
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	basename := filepath.Base(dir)

	if k8sSecret == "" {
		return fmt.Sprintf("%s/.env", basename), nil
	} else {
		return fmt.Sprintf("%s/%s-secret.yaml", basename, k8sSecret), nil
	}
}
