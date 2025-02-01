package optruck

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/yammerjp/optruck/pkg/actions"
	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/kube"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

func (cli *CLI) buildWithDefault() (actions.Action, error) {
	if err := cli.SetDefaultIfEmpty(); err != nil {
		return nil, err
	}
	return cli.build()
}

func (cli *CLI) build() (actions.Action, error) {
	ds, err := cli.buildDataSource()
	if err != nil {
		return nil, err
	}

	dest, err := cli.buildDest()
	if err != nil {
		return nil, err
	}

	opItemClient, err := cli.buildOpItemClient(true)
	if err != nil {
		return nil, err
	}

	return &actions.MirrorConfig{
		Logger:       cli.buildLogger(),
		OpItemClient: *opItemClient,
		DataSource:   ds,
		Dest:         dest,
		Overwrite:    cli.Overwrite,
	}, nil
}

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

func (cli *CLI) buildOpItemClient(strict bool) (*op.ItemClient, error) {
	if strict {
		if cli.Account == "" {
			return nil, fmt.Errorf("--account is required")
		}
		if cli.Vault == "" {
			return nil, fmt.Errorf("--vault is required")
		}
		if cli.Item == "" {
			return nil, fmt.Errorf("item specification is required")
		}
	}

	return op.NewItemClient(cli.Account, cli.Vault, cli.Item, cli.exec), nil
}

func (cli *CLI) buildDataSource() (datasources.Source, error) {
	if cli.EnvFile != "" {
		return &datasources.EnvFileSource{Path: cli.EnvFile}, nil
	}
	if cli.K8sSecret != "" {
		if cli.K8sNamespace == "" {
			return nil, fmt.Errorf("--k8s-namespace is required when using --k8s-secret")
		}
		return &datasources.K8sSecretSource{
			Namespace:  cli.K8sNamespace,
			SecretName: cli.K8sSecret,
			Client:     kube.NewClient(cli.exec),
		}, nil
	}
	return nil, fmt.Errorf("either --env-file or --k8s-secret is required")
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
