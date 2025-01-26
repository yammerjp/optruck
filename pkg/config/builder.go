package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/yammerjp/optruck/pkg/actions"
	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

type ConfigBuilder struct {
	item         string
	vault        string
	account      string
	envFile      string
	k8sSecret    string
	k8sNamespace string
	output       string
	outputFormat string
	overwrite    bool
	interactive  bool
	logLevel     string
	logFile      string
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

func (b *ConfigBuilder) WithInteractive(interactive bool) *ConfigBuilder {
	b.interactive = interactive
	return b
}

func (b *ConfigBuilder) WithItem(item string) *ConfigBuilder {
	b.item = item
	return b
}

func (b *ConfigBuilder) WithVault(vault string) *ConfigBuilder {
	b.vault = vault
	return b
}

func (b *ConfigBuilder) WithAccount(account string) *ConfigBuilder {
	b.account = account
	return b
}

func (b *ConfigBuilder) WithEnvFile(envFile string) *ConfigBuilder {
	b.envFile = envFile
	return b
}

func (b *ConfigBuilder) WithK8sSecret(secret string) *ConfigBuilder {
	b.k8sSecret = secret
	return b
}

func (b *ConfigBuilder) WithK8sNamespace(namespace string) *ConfigBuilder {
	b.k8sNamespace = namespace
	return b
}

func (b *ConfigBuilder) WithOutput(output string) *ConfigBuilder {
	b.output = output
	return b
}

func (b *ConfigBuilder) WithOutputFormat(format string) *ConfigBuilder {
	b.outputFormat = format
	return b
}

func (b *ConfigBuilder) WithOverwrite(overwrite bool) *ConfigBuilder {
	b.overwrite = overwrite
	return b
}

func (b *ConfigBuilder) WithLogLevel(logLevel string) *ConfigBuilder {
	b.logLevel = logLevel
	return b
}

func (b *ConfigBuilder) WithLogFile(logFile string) *ConfigBuilder {
	b.logFile = logFile
	return b
}

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

func (b *ConfigBuilder) validateCommon() error {
	if b.item == "" {
		return fmt.Errorf("item is required")
	}
	if len(b.item) > 100 {
		return fmt.Errorf("item must be less than 100 characters")
	}
	if b.vault == "" {
		// FIXME: not use op.Client directly in plg/config
		// if vault is not specified and exist only one vault, use it
		opClient := b.buildOpTarget().BuildClient()
		vaults, err := opClient.ListVaults()
		if err != nil {
			return fmt.Errorf("failed to list vaults: %w", err)
		}
		if len(vaults) == 1 {
			b.vault = vaults[0].Name
		} else {
			return fmt.Errorf("vault is required")
		}
	}
	if b.account == "" {
		// FIXME: not use op.Client directly in plg/config
		// if account is not specified and exist only one account, use it
		opClient := b.buildOpTarget().BuildClient()
		accounts, err := opClient.ListAccounts()
		if err != nil {
			return fmt.Errorf("failed to list accounts: %w", err)
		}
		if len(accounts) == 1 {
			b.account = accounts[0].URL
		} else {
			return fmt.Errorf("account is required")
		}
	}
	return nil
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
		if b.envFile != "" {
			return nil, fmt.Errorf("cannot use both --k8s-secret and --env-file")
		}
		if b.k8sNamespace == "" {
			return nil, fmt.Errorf("k8s namespace is required")
		}
		return &datasources.K8sSecretSource{
			Namespace:  b.k8sNamespace,
			SecretName: b.k8sSecret,
			Client:     datasources.NewK8sClient(),
		}, nil
	}

	if b.envFile == "" {
		return nil, fmt.Errorf("env file is required")
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
		Logger:     logger,
		Target:     b.buildOpTarget(),
		DataSource: ds,
		Overwrite:  b.overwrite,
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
		Logger:    logger,
		Target:    b.buildOpTarget(),
		Dest:      dest,
		Overwrite: b.overwrite,
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
		Logger:     logger,
		Target:     b.buildOpTarget(),
		DataSource: ds,
		Dest:       dest,
		Overwrite:  b.overwrite,
	}, cleanup, nil
}

func (b *ConfigBuilder) Build(isUpload, isTemplate, isMirror bool) (actions.Action, func(), error) {
	if !isUpload && !isTemplate {
		// mirror is default actions.Action
		return b.BuildMirror()
	}
	if isUpload && !isTemplate && !isMirror {
		return b.BuildUpload()
	}
	if !isUpload && isTemplate && !isMirror {
		return b.BuildTemplate()
	}
	return nil, nil, fmt.Errorf("only one of --upload, --template, or --mirror can be specified")
}
