package optruck

import (
	"io"
	"log/slog"
	"os"

	"github.com/yammerjp/optruck/pkg/actions"
	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/kube"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

func (cli *CLI) buildLogger() *slog.Logger {
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

	return slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: logLevel}))
}

func (cli *CLI) buildOpTarget() op.Target {
	return op.Target{
		Account:  cli.Account,
		Vault:    cli.Vault,
		ItemName: cli.Item,
	}
}

func (cli *CLI) buildDataSource() (datasources.Source, error) {
	if cli.K8sSecret != "" {
		return &datasources.K8sSecretSource{
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
			Client:     kube.NewClient(),
		}, nil
	}

	return &datasources.EnvFileSource{Path: cli.EnvFile}, nil
}

func (cli *CLI) buildDest() (output.Dest, error) {
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

func (cli *CLI) Build() (actions.Action, error) {
	if err := cli.validateCommon(); err != nil {
		return nil, err
	}

	ds, err := cli.buildDataSource()
	if err != nil {
		return nil, err
	}

	dest, err := cli.buildDest()
	if err != nil {
		return nil, err
	}

	return &actions.MirrorConfig{
		Logger:     cli.buildLogger(),
		Target:     cli.buildOpTarget(),
		DataSource: ds,
		Dest:       dest,
		Overwrite:  cli.Overwrite,
	}, nil
}
