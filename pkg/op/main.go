package op

import (
	"k8s.io/utils/exec"
)

type Client struct {
	exec        exec.Interface
	AccountName string
	VaultName   string
}

func NewClient(accountName string, vaultName string) *Client {
	return &Client{
		exec:        exec.New(),
		AccountName: accountName,
		VaultName:   vaultName,
	}
}
