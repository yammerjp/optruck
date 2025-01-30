package op

import execPackage "k8s.io/utils/exec"

type ExecutableClient struct {
	exec execPackage.Interface
}

func NewExecutableClient(exec execPackage.Interface) *ExecutableClient {
	if exec == nil {
		exec = execPackage.New()
	}
	return &ExecutableClient{
		exec: exec,
	}
}

type AccountClient struct {
	ExecutableClient
	Account string
}

func NewAccountClient(account string, exec execPackage.Interface) *AccountClient {
	if exec == nil {
		exec = execPackage.New()
	}
	return &AccountClient{
		ExecutableClient: *NewExecutableClient(exec),
		Account:          account,
	}
}

type VaultClient struct {
	AccountClient
	Vault string
}

func NewVaultClient(account, vault string, exec execPackage.Interface) *VaultClient {
	return &VaultClient{
		AccountClient: *NewAccountClient(account, exec),
		Vault:         vault,
	}
}

type ItemClient struct {
	VaultClient
	ItemName string
}

func NewItemClient(account, vault, itemName string, exec execPackage.Interface) *ItemClient {
	return &ItemClient{
		VaultClient: *NewVaultClient(account, vault, exec),
		ItemName:    itemName,
	}
}
