package config

import (
	"fmt"
	"log/slog"

	"github.com/yammerjp/optruck/pkg/actions"
	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

// Action は実行可能なアクションを表すインターフェース
type Action interface {
	Run() error
}

// ConfigBuilder は設定の構築を担当する
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
	logger       *slog.Logger
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
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

func (b *ConfigBuilder) WithLogger(logger *slog.Logger) *ConfigBuilder {
	b.logger = logger
	return b
}

// validateCommon は共通のバリデーションルールを実行
func (b *ConfigBuilder) validateCommon() error {
	if b.item == "" {
		return fmt.Errorf("item is required")
	}
	if len(b.item) > 100 {
		return fmt.Errorf("item must be less than 100 characters")
	}
	if b.vault == "" {
		return fmt.Errorf("vault is required")
	}
	if b.account == "" {
		return fmt.Errorf("account is required")
	}
	return nil
}

// buildOpTarget は1Password操作のターゲットを構築
func (b *ConfigBuilder) buildOpTarget() op.Target {
	return op.Target{
		Account:  b.account,
		Vault:    b.vault,
		ItemName: b.item,
	}
}

// buildDataSource はデータソースを構築
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

// buildDest は出力先を構築
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

// BuildUpload はアップロードアクションを構築
func (b *ConfigBuilder) BuildUpload() (Action, error) {
	if err := b.validateCommon(); err != nil {
		return nil, err
	}

	ds, err := b.buildDataSource()
	if err != nil {
		return nil, err
	}

	return &actions.UploadConfig{
		Logger:     b.logger,
		Target:     b.buildOpTarget(),
		DataSource: ds,
		Overwrite:  b.overwrite,
	}, nil
}

// BuildTemplate はテンプレート生成アクションを構築
func (b *ConfigBuilder) BuildTemplate() (Action, error) {
	if err := b.validateCommon(); err != nil {
		return nil, err
	}

	dest, err := b.buildDest()
	if err != nil {
		return nil, err
	}

	return &actions.TemplateConfig{
		Logger:    b.logger,
		Target:    b.buildOpTarget(),
		Dest:      dest,
		Overwrite: b.overwrite,
	}, nil
}

// BuildMirror はミラーアクションを構築
func (b *ConfigBuilder) BuildMirror() (Action, error) {
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

	return &actions.MirrorConfig{
		Logger:     b.logger,
		Target:     b.buildOpTarget(),
		DataSource: ds,
		Dest:       dest,
		Overwrite:  b.overwrite,
	}, nil
}

// Build は指定されたアクションを構築
func (b *ConfigBuilder) Build(isUpload, isTemplate bool) (Action, error) {
	if isUpload && !isTemplate {
		return b.BuildUpload()
	}
	if !isUpload && isTemplate {
		return b.BuildTemplate()
	}
	return b.BuildMirror()
}
