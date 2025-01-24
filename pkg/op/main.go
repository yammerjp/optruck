// TODO: test
package op

import (
	"k8s.io/utils/exec"
)

type Target struct {
	AccountName string
	VaultName   string
	ItemName    string
}

type Client struct {
	exec exec.Interface
	Target
}

func NewClient(target Target) *Client {
	return &Client{
		exec:   exec.New(),
		Target: target,
	}
}
func (c *Client) BuildItemCommand(args ...string) exec.Cmd {
	if c.Target.AccountName != "" {
		args = append(args, "--account", c.Target.AccountName)
	}
	if c.Target.VaultName != "" {
		args = append(args, "--vault", c.Target.VaultName)
	}
	args = append(args, "--format", "json", "item")
	return c.exec.Command("op", args...)
}
