package config

type ConfigBuilder struct {
	item         string
	vault        string
	account      string
	envFile      string
	k8sSecret    string
	k8sNamespace string
	output       string
	overwrite    bool
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

func (b *ConfigBuilder) WithOverwrite(overwrite bool) *ConfigBuilder {
	b.overwrite = overwrite
	return b
}
