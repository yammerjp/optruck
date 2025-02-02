package op

type ExecutableClient struct {
}

func NewExecutableClient() *ExecutableClient {
	return &ExecutableClient{}
}

type AccountClient struct {
	ExecutableClient
	Account string
}

func NewAccountClient(account string) *AccountClient {
	return &AccountClient{
		ExecutableClient: *NewExecutableClient(),
		Account:          account,
	}
}

type VaultClient struct {
	AccountClient
	Vault string
}

func NewVaultClient(account, vault string) *VaultClient {
	return &VaultClient{
		AccountClient: *NewAccountClient(account),
		Vault:         vault,
	}
}

type ItemClient struct {
	VaultClient
	ItemName string
}

func NewItemClient(account, vault, itemName string) *ItemClient {
	return &ItemClient{
		VaultClient: *NewVaultClient(account, vault),
		ItemName:    itemName,
	}
}
