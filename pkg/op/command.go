package op

import "k8s.io/utils/exec"

type CommandOptions struct {
	AddAccount bool
	AddVault   bool
	Args       []string
}

func (c *ExecutableClient) BuildCommand(args ...string) exec.Cmd {
	args = append(args, "--format", "json")
	return c.exec.Command("op", args...)
}

func (c *AccountClient) BuildCommand(args ...string) exec.Cmd {
	args = append(args, "--account", c.Account)
	args = append(args, "--format", "json")
	return c.exec.Command("op", args...)
}

func (c *VaultClient) BuildCommand(args ...string) exec.Cmd {
	args = append(args, "--account", c.Account)
	args = append(args, "--vault", c.Vault)
	args = append(args, "--format", "json")
	return c.exec.Command("op", args...)
}
