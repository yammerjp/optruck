package op

import (
	utilExec "github.com/yammerjp/optruck/internal/util/exec"
)

type CommandOptions struct {
	AddAccount bool
	AddVault   bool
	Args       []string
}

func (c *ExecutableClient) BuildCommand(args ...string) utilExec.Command {
	args = append(args, "--format", "json")
	return utilExec.NewCommand("op", args...)
}

func (c *AccountClient) BuildCommand(args ...string) utilExec.Command {
	args = append(args, "--account", c.Account)
	args = append(args, "--format", "json")
	return utilExec.NewCommand("op", args...)
}

func (c *VaultClient) BuildCommand(args ...string) utilExec.Command {
	args = append(args, "--account", c.Account)
	args = append(args, "--vault", c.Vault)
	args = append(args, "--format", "json")
	return utilExec.NewCommand("op", args...)
}
