package optruck

import (
	"fmt"

	"github.com/yammerjp/optruck/internal/errors"
	"github.com/yammerjp/optruck/internal/interactive"
	"github.com/yammerjp/optruck/pkg/actions"
	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/kube"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

func (cli *CLI) buildAction(confirmation func() error) (actions.Action, error) {
	client, err := cli.buildOpItemClient(true)
	if err != nil {
		return nil, err
	}

	source, err := cli.buildDataSource()
	if err != nil {
		return nil, err
	}

	dest, err := cli.buildDest()
	if err != nil {
		return nil, err
	}

	return &actions.MirrorConfig{
		OpItemClient: *client,
		DataSource:   source,
		Dest:         dest,
		Overwrite:    cli.Overwrite,
		Confirmation: confirmation,
	}, nil
}

func (cli *CLI) buildOpItemClient(strict bool) (*op.ItemClient, error) {
	if strict {
		if cli.Account == "" {
			accounts, err := op.NewExecutableClient().ListAccounts()
			if err != nil {
				return nil, errors.NewOperationFailedError(
					"1Passwordアカウントの一覧取得",
					err,
					"1Password CLIが正しくインストールされ、認証されていることを確認してください",
				)
			}
			if len(accounts) != 1 {
				return nil, errors.NewMissingRequirementError(
					"1Passwordアカウントの指定",
					"--account オプションで1Passwordアカウントを指定するか、--interactive オプションを使用してください",
				)
			}
			cli.Account = accounts[0].URL
		}
		if cli.Vault == "" {
			vaults, err := op.NewAccountClient(cli.Account).ListVaults()
			if err != nil {
				return nil, errors.NewOperationFailedError(
					"1Passwordボールトの一覧取得",
					err,
					"1Password CLIが正しく設定されていることを確認してください",
				)
			}
			if len(vaults) != 1 {
				return nil, errors.NewMissingRequirementError(
					"1Passwordボールトの指定",
					"--vault オプションで1Passwordボールトを指定するか、--interactive オプションを使用してください",
				)
			}
			cli.Vault = vaults[0].Name
		}
		if cli.Item == "" {
			return nil, errors.NewMissingRequirementError(
				"1Passwordアイテム名またはID",
				"コマンド引数でアイテム名を指定するか、--interactive オプションを使用してください",
			)
		}
	}

	return op.NewItemClient(cli.Account, cli.Vault, cli.Item), nil
}

func (cli *CLI) buildDataSource() (datasources.Source, error) {
	if cli.K8sSecret != "" {
		if cli.K8sNamespace == "" {
			cli.K8sNamespace = interactive.DefaultKubernetesNamespace
		}
		client := kube.NewClient()
		// Validate that the namespace and secret exist
		namespaces, err := client.GetNamespaces()
		if err != nil {
			return nil, errors.NewOperationFailedError(
				"Kubernetesネームスペースの一覧取得",
				err,
				"kubectlが正しく設定されていることを確認してください",
			)
		}
		namespaceExists := false
		for _, ns := range namespaces {
			if ns == cli.K8sNamespace {
				namespaceExists = true
				break
			}
		}
		if !namespaceExists {
			return nil, errors.NewNotFoundError(
				fmt.Sprintf("Kubernetesネームスペース '%s'", cli.K8sNamespace),
				"--k8s-namespace オプションで正しいネームスペースを指定してください",
			)
		}

		return &datasources.K8sSecretSource{
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
			Client:     client,
		}, nil
	}
	if cli.EnvFile == "" {
		cli.EnvFile = interactive.DefaultEnvFilePath
	}
	return &datasources.EnvFileSource{Path: cli.EnvFile}, nil
}

func (cli *CLI) buildDest() (output.Dest, error) {
	if cli.Output == "" {
		cli.Output = interactive.DefaultOutputPath(cli.K8sSecret)
	}
	if cli.K8sSecret != "" {
		return &output.K8sSecretTemplateDest{
			Path:       cli.Output,
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
		}, nil
	}

	return &output.EnvTemplateDest{
		Path: cli.Output,
	}, nil
}
