package optruck

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/yammerjp/optruck/pkg/config"
)

func Run() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("optruck"),
		kong.Description("A CLI tool for managing secrets and creating templates with 1Password."),
		kong.UsageOnError(),
		kong.Help(helpPrinter),
	)

	logger, err := cli.buildLogger()
	if err != nil {
		ctx.Fatalf("%v", err)
	}

	// Handle version flag
	if cli.Version {
		fmt.Printf("optruck version %s\n", version)
		os.Exit(0)
	}

	builder := config.NewConfigBuilder().
		WithItem(cli.Item).
		WithVault(cli.Vault).
		WithAccount(cli.Account).
		WithEnvFile(cli.EnvFile).
		WithK8sSecret(cli.K8sSecret).
		WithK8sNamespace(cli.K8sNamespace).
		WithOutput(cli.Output).
		WithOverwrite(cli.Overwrite)

	if cli.Interactive {
		err := builder.SetConfigInteractively()
		if err != nil {
			ctx.Fatalf("%v", err)
		}
	}
	if err := builder.SetDefaultIfEmpty(); err != nil {
		ctx.Fatalf("%v", err)
	}

	action, err := builder.Build(logger)
	if err != nil {
		ctx.Fatalf("%v", err)
	}

	if err := action.Run(); err != nil {
		ctx.Fatalf("%v", err)
	}
}

func (cli *CLI) buildLogger() (*slog.Logger, error) {
	var logLevel slog.Level
	var f io.Writer
	switch cli.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
		f = os.Stderr
	case "info":
		logLevel = slog.LevelInfo
		f = os.Stderr
	case "warn":
		logLevel = slog.LevelWarn
		f = os.Stderr
	case "error":
		logLevel = slog.LevelError
		f = os.Stderr
	default:
		logLevel = slog.LevelInfo
		f = io.Discard
	}

	return slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: logLevel})), nil
}
