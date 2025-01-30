package config

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

func (b *ConfigBuilder) BuildLogger() (*slog.Logger, error) {
	var logLevel slog.Level
	var f io.Writer
	switch b.logLevel {
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

func (b *ConfigBuilder) buildOpTarget() op.Target {
	return op.Target{
		Account:  b.account,
		Vault:    b.vault,
		ItemName: b.item,
	}
}

func (b *ConfigBuilder) buildDataSource() (datasources.Source, error) {
	if b.k8sSecret != "" {
		return &datasources.K8sSecretSource{
			Namespace:  b.k8sNamespace,
			SecretName: b.k8sSecret,
			Client:     kube.NewClient(),
		}, nil
	}

	return &datasources.EnvFileSource{Path: b.envFile}, nil
}

func (b *ConfigBuilder) buildDest() (output.Dest, error) {
	if b.k8sSecret != "" {
		return &output.K8sSecretTemplateDest{
			Path:       b.output,
			Namespace:  b.k8sNamespace,
			SecretName: b.k8sSecret,
		}, nil
	}

	return &output.EnvTemplateDest{
		Path: b.output,
	}, nil
}

func (b *ConfigBuilder) Build() (actions.Action, error) {
	if err := b.validateCommon(); err != nil {
		return nil, err
	}

	ds, err := b.buildDataSource()
	if err != nil {
		return nil, err
	}

	dest, err := b.buildDest()
	if err != nil {
		return nil, err
	}

	logger, err := b.BuildLogger()
	if err != nil {
		return nil, err
	}

	return &actions.MirrorConfig{
		Logger:            logger,
		Target:            b.buildOpTarget(),
		DataSource:        ds,
		Dest:              dest,
		OverwriteTarget:   b.overwrite,
		OverwriteTemplate: b.overwrite,
	}, nil
}
