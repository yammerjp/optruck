package optruck

import "github.com/yammerjp/optruck/internal/interactive"

func (cli CLI) buildResultCommand() ([]string, error) {
	cmds := []string{"optruck", cli.Item}
	// target options
	if cli.Overwrite {
		cmds = append(cmds, "--overwrite")
	}
	if cli.Account != "" {
		cmds = append(cmds, "--account", cli.Account)
	}
	if cli.Vault != "" {
		cmds = append(cmds, "--vault", cli.Vault)
	}

	// data source options
	if cli.EnvFile != "" {
		cmds = append(cmds, "--env-file", cli.EnvFile)
	} else if cli.K8sSecret != "" {
		cmds = append(cmds, "--k8s-secret", cli.K8sSecret)
		if cli.K8sNamespace != interactive.DefaultKubernetesNamespace {
			cmds = append(cmds, "--k8s-namespace", cli.K8sNamespace)
		}
	}

	// output options
	if cli.Output != "" {
		cmds = append(cmds, "--output", cli.Output)
	}
	return cmds, nil
}
