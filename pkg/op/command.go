package op

import "k8s.io/utils/exec"

type CommandOptions struct {
	AddAccount bool
	AddVault   bool
	Args       []string
}

func (c *Client) BuildCommand(opts CommandOptions) exec.Cmd {
	args := opts.Args
	if opts.AddAccount && c.Target.Account != "" {
		args = append(args, "--account", c.Target.Account)
	}
	if opts.AddVault && c.Target.Vault != "" {
		args = append(args, "--vault", c.Target.Vault)
	}
	args = append(args, "--format", "json")
	return c.exec.Command("op", args...)
}
