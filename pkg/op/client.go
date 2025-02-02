package op

import (
	"log/slog"

	execPackage "k8s.io/utils/exec"
)

type ExecutableClient struct {
	exec   execPackage.Interface
	logger *slog.Logger
}

func NewExecutableClient(exec execPackage.Interface, logger *slog.Logger) *ExecutableClient {
	return &ExecutableClient{
		exec:   exec,
		logger: logger,
	}
}

type AccountClient struct {
	ExecutableClient
	Account string
}

func NewAccountClient(account string, exec execPackage.Interface, logger *slog.Logger) *AccountClient {
	return &AccountClient{
		ExecutableClient: *NewExecutableClient(exec, logger),
		Account:          account,
	}
}

type VaultClient struct {
	AccountClient
	Vault string
}

func NewVaultClient(account, vault string, exec execPackage.Interface, logger *slog.Logger) *VaultClient {
	return &VaultClient{
		AccountClient: *NewAccountClient(account, exec, logger),
		Vault:         vault,
	}
}

type ItemClient struct {
	VaultClient
	ItemName string
}

func NewItemClient(account, vault, itemName string, exec execPackage.Interface, logger *slog.Logger) *ItemClient {
	return &ItemClient{
		VaultClient: *NewVaultClient(account, vault, exec, logger),
		ItemName:    itemName,
	}
}
