package optruck

import (
	"fmt"
	"os"

	"github.com/yammerjp/optruck/internal/errors"
	"github.com/yammerjp/optruck/internal/interactive"
	utilLogger "github.com/yammerjp/optruck/internal/util/logger"

	"github.com/alecthomas/kong"
)

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

var buildInfo BuildInfo

func Run(bi BuildInfo) {
	buildInfo = bi

	cli := CLI{}
	_ = kong.Parse(&cli,
		kong.Name("optruck"),
		kong.Description("A CLI tool for managing secrets and creating templates with 1Password."),
		kong.UsageOnError(),
		kong.Help(helpPrinter),
	)

	if err := cli.Run(); err != nil {
		// Check if it's our custom error type
		if userErr, ok := err.(*errors.UserError); ok {
			fmt.Fprintf(os.Stderr, "エラー: %s\n", userErr.Error())
			os.Exit(1)
		}
		// For unknown errors, provide a generic message
		fmt.Fprintf(os.Stderr, "予期せぬエラーが発生しました: %v\n提案: バグの可能性があります。--log-level debug オプションを付けて実行し、詳細なログを確認してください。\n", err)
		os.Exit(1)
	}
}

func (cli *CLI) Run() error {
	utilLogger.SetDefaultLogger(cli.LogLevel)

	var confirmation func() error

	if cli.Interactive {
		runner := *interactive.NewImplRunner()
		if err := cli.SetOptionsInteractively(runner); err != nil {
			return errors.WrapError(
				err,
				"対話モードでの設定に失敗しました",
				"入力内容を確認し、再度実行してください",
			)
		}
		cmds, err := cli.buildResultCommand()
		if err != nil {
			return errors.WrapError(
				err,
				"コマンドの構築に失敗しました",
				"バグの可能性があります。--log-level debug オプションを付けて実行し、詳細なログを確認してください",
			)
		}
		confirmation = func() error {
			return runner.Confirm(cmds)
		}
	} else {
		confirmation = func() error {
			// confirmed by default
			return nil
		}
	}

	action, err := cli.buildAction(confirmation)
	if err != nil {
		return err // buildActionで既に適切なエラーメッセージが設定されている
	}

	if err := action.Run(); err != nil {
		return errors.WrapError(
			err,
			"処理の実行に失敗しました",
			"エラーメッセージを確認し、必要な対応を行ってください",
		)
	}

	return nil
}
