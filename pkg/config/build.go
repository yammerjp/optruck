package config

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

func (b *ConfigBuilder) BuildLogger() (*slog.Logger, func(), error) {
	var logLevel slog.Level
	switch b.logLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var f io.Writer
	var cleanup func()

	if b.logFile == "" {
		f = io.Discard
		cleanup = func() {}
	} else {
		var err error
		f, err = os.OpenFile(b.logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, nil, err
		}
		cleanup = func() {
			f.(io.WriteCloser).Close()
		}
	}
	return slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: logLevel})), cleanup, nil
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
	if b.outputFormat == "k8s" {
		if b.k8sSecret == "" {
			return nil, fmt.Errorf("k8s secret is required")
		}
		if b.k8sNamespace == "" {
			return nil, fmt.Errorf("k8s namespace is required")
		}
		return &output.K8sSecretTemplateDest{
			Path:       b.output,
			Namespace:  b.k8sNamespace,
			SecretName: b.k8sSecret,
		}, nil
	}

	if b.outputFormat == "env" || b.envFile != "" {
		return &output.EnvTemplateDest{
			Path: b.output,
		}, nil
	}

	if b.k8sSecret != "" {
		if b.k8sNamespace == "" {
			return nil, fmt.Errorf("k8s namespace is required")
		}
		return &output.K8sSecretTemplateDest{
			Path:       b.output,
			Namespace:  b.k8sNamespace,
			SecretName: b.k8sSecret,
		}, nil
	}

	return nil, fmt.Errorf("invalid output format: %s", b.outputFormat)
}

func (b *ConfigBuilder) BuildUpload() (actions.Action, func(), error) {
	if err := b.validateCommon(); err != nil {
		return nil, nil, err
	}

	ds, err := b.buildDataSource()
	if err != nil {
		return nil, nil, err
	}

	logger, cleanup, err := b.BuildLogger()
	if err != nil {
		return nil, nil, err
	}

	return &actions.UploadConfig{
		Logger:          logger,
		Target:          b.buildOpTarget(),
		DataSource:      ds,
		OverwriteTarget: b.overwriteTarget || b.overwrite,
	}, cleanup, nil
}

func (b *ConfigBuilder) BuildTemplate() (actions.Action, func(), error) {
	if err := b.validateCommon(); err != nil {
		return nil, nil, err
	}

	dest, err := b.buildDest()
	if err != nil {
		return nil, nil, err
	}

	logger, cleanup, err := b.BuildLogger()
	if err != nil {
		return nil, nil, err
	}

	return &actions.TemplateConfig{
		Logger:            logger,
		Target:            b.buildOpTarget(),
		Dest:              dest,
		OverwriteTemplate: b.overwriteTemplate || b.overwrite,
	}, cleanup, nil
}

func (b *ConfigBuilder) BuildMirror() (actions.Action, func(), error) {
	if err := b.validateCommon(); err != nil {
		return nil, nil, err
	}

	ds, err := b.buildDataSource()
	if err != nil {
		return nil, nil, err
	}

	dest, err := b.buildDest()
	if err != nil {
		return nil, nil, err
	}

	logger, cleanup, err := b.BuildLogger()
	if err != nil {
		return nil, nil, err
	}

	return &actions.MirrorConfig{
		Logger:            logger,
		Target:            b.buildOpTarget(),
		DataSource:        ds,
		Dest:              dest,
		OverwriteTarget:   b.overwriteTarget || b.overwrite,
		OverwriteTemplate: b.overwriteTemplate || b.overwrite,
	}, cleanup, nil
}

func (b *ConfigBuilder) Build() (actions.Action, func(), error) {
	if !b.isUpload && !b.isTemplate {
		// mirror is default actions.Action
		return b.BuildMirror()
	}
	if b.isUpload && !b.isTemplate && !b.isMirror {
		return b.BuildUpload()
	}
	if !b.isUpload && b.isTemplate && !b.isMirror {
		return b.BuildTemplate()
	}
	return nil, nil, fmt.Errorf("only one of --upload, --template, or --mirror can be specified")
}
