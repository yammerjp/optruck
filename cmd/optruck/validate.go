package optruck

import (
	"fmt"
)

// TODO: test

func (cli *CLI) validateConflictOptions() error {
	if len(cli.Item) > 100 {
		return fmt.Errorf("item must be less than 100 characters")
	}
	if cli.EnvFile != "" && (cli.K8sSecret != "" || cli.K8sNamespace != "") {
		return fmt.Errorf("cannot use both --env-file and --k8s-secret or --k8s-namespace")
	}

	return nil
}
