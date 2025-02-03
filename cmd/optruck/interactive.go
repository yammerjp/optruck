package optruck

import (
	"fmt"
	"log/slog"

	"github.com/yammerjp/optruck/internal/errors"
	"github.com/yammerjp/optruck/internal/interactive"
)

func (cli *CLI) SetOptionsInteractively(runner interactive.Runner) error {
	if err := cli.setDataSourceInteractively(runner); err != nil {
		return err
	}
	if err := cli.setTargetInteractively(runner); err != nil {
		return err
	}
	if err := cli.setDestInteractively(runner); err != nil {
		return err
	}
	return nil
}

func (cli *CLI) setDataSourceInteractively(runner interactive.Runner) error {
	if cli.EnvFile != "" || cli.K8sSecret != "" {
		slog.Debug("data source already set", "envFile", cli.EnvFile, "k8sSecret", cli.K8sSecret)
		// already set
		return nil
	}
	ds, err := runner.SelectDataSource()
	if err != nil {
		return errors.WrapError(
			err,
			"データソースの選択に失敗しました",
			"入力内容を確認し、再度実行してください",
		)
	}
	switch ds {
	case interactive.DataSourceEnvFile:
		slog.Debug("setting env file path")
		envFilePath, err := runner.PromptEnvFilePath()
		if err != nil {
			return errors.WrapError(
				err,
				"環境変数ファイルのパス指定に失敗しました",
				"指定したパスが正しいか、アクセス権限があるか確認してください",
			)
		}
		cli.EnvFile = envFilePath
	case interactive.DataSourceK8sSecret:
		slog.Debug("setting k8s secret")
		if cli.K8sNamespace == "" {
			namespace, err := runner.SelectKubeNamespace()
			if err != nil {
				return errors.WrapError(
					err,
					"Kubernetesネームスペースの選択に失敗しました",
					"kubectlが正しく設定されているか確認してください",
				)
			}
			cli.K8sNamespace = namespace
		}
		if cli.K8sSecret == "" {
			secret, err := runner.SelectKubeSecret(cli.K8sNamespace)
			if err != nil {
				return errors.WrapError(
					err,
					fmt.Sprintf("Kubernetesシークレットの選択に失敗しました (namespace: %s)", cli.K8sNamespace),
					"指定したネームスペースにシークレットが存在するか確認してください",
				)
			}
			cli.K8sSecret = secret
		}
	default:
		return errors.NewInvalidArgumentError(
			"データソース",
			fmt.Sprintf("不正な値です: %s", ds),
			"環境変数ファイルまたはKubernetesシークレットを選択してください",
		)
	}
	return nil
}

func (cli *CLI) setTargetInteractively(runner interactive.Runner) error {
	if cli.Account == "" {
		account, err := runner.SelectOpAccount()
		if err != nil {
			return errors.WrapError(
				err,
				"1Passwordアカウントの選択に失敗しました",
				"1Password CLIが正しくインストールされ、認証されていることを確認してください",
			)
		}
		cli.Account = account
	}

	if cli.Vault == "" {
		vault, err := runner.SelectOpVault(cli.Account)
		if err != nil {
			return errors.WrapError(
				err,
				"1Passwordボールトの選択に失敗しました",
				"指定したアカウントにアクセス権限があるか確認してください",
			)
		}
		cli.Vault = vault
	}
	if cli.Item == "" {
		if !cli.Overwrite {
			overwrite, err := runner.SelectOpItemOverwriteOrNot()
			if err != nil {
				return errors.WrapError(
					err,
					"上書きモードの選択に失敗しました",
					"入力内容を確認し、再度実行してください",
				)
			}
			cli.Overwrite = overwrite
		}
		if cli.Overwrite {
			itemName, err := runner.SelectOpItemName(cli.Account, cli.Vault)
			if err != nil {
				return errors.WrapError(
					err,
					"1Passwordアイテムの選択に失敗しました",
					"指定したボールトにアクセス権限があるか確認してください",
				)
			}
			cli.Item = itemName
		} else {
			itemName, err := runner.PromptOpItemName(cli.Account, cli.Vault, cli.K8sSecret)
			if err != nil {
				return errors.WrapError(
					err,
					"1Passwordアイテム名の入力に失敗しました",
					"入力内容を確認し、再度実行してください",
				)
			}
			cli.Item = itemName
		}
	}
	return nil
}

func (cli *CLI) setDestInteractively(runner interactive.Runner) error {
	if cli.Output != "" {
		// already set
		return nil
	}
	outputPath, err := runner.PromptOutputPath(cli.K8sSecret)
	if err != nil {
		return errors.WrapError(
			err,
			"出力先パスの指定に失敗しました",
			"指定したパスが正しいか、書き込み権限があるか確認してください",
		)
	}
	cli.Output = outputPath
	return nil
}
