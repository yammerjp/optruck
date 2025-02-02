package op

import (
	optruckexec "github.com/yammerjp/optruck/pkg/exec"
)

func (c *ExecutableClient) BuildCommand(args ...string) *optruckexec.Command {
	args = append(args, "--format", "json")
	return c.Command("op", args...)
}

func (c *AccountClient) BuildCommand(args ...string) *optruckexec.Command {
	args = append(args, "--account", c.Account)
	args = append(args, "--format", "json")
	return c.Command("op", args...)
}

func (c *VaultClient) BuildCommand(args ...string) *optruckexec.Command {
	args = append(args, "--account", c.Account)
	args = append(args, "--vault", c.Vault)
	args = append(args, "--format", "json")
	return c.Command("op", args...)
}
