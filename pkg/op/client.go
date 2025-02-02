package op

import (
	"log/slog"

	optruckexec "github.com/yammerjp/optruck/pkg/exec"
	"k8s.io/utils/exec"
)

type ExecutableClient struct {
	optruckexec.CommandConfig
}

func NewExecutableClient(exec exec.Interface, logger *slog.Logger) *ExecutableClient {
	return &ExecutableClient{
		CommandConfig: optruckexec.NewCommandConfig(exec, logger),
	}
}

func NewExecutableClientFromConfig(exec optruckexec.CommandConfig) *ExecutableClient {
	return &ExecutableClient{
		CommandConfig: exec,
	}
}

type AccountClient struct {
	ExecutableClient
	Account string
}

func (c *ExecutableClient) BuildAccountClient(account string) *AccountClient {
	return &AccountClient{
		ExecutableClient: *c,
		Account:          account,
	}
}

type VaultClient struct {
	AccountClient
	Vault string
}

func (c *ExecutableClient) BuildVaultClient(account, vault string) *VaultClient {
	return &VaultClient{
		AccountClient: *c.BuildAccountClient(account),
		Vault:         vault,
	}
}

type ItemClient struct {
	VaultClient
	ItemName string
}

func (c *ExecutableClient) BuildItemClient(account, vault, itemName string) *ItemClient {
	return &ItemClient{
		VaultClient: *c.BuildVaultClient(account, vault),
		ItemName:    itemName,
	}
}
