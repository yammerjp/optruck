package config

type ConfigBuilder struct {
	item              string
	vault             string
	account           string
	envFile           string
	k8sSecret         string
	k8sNamespace      string
	output            string
	outputFormat      string
	overwrite         bool
	overwriteTarget   bool
	overwriteTemplate bool
	logLevel          string
	logFile           string
	isUpload          bool
	isTemplate        bool
	isMirror          bool
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

func (b *ConfigBuilder) WithOverwriteTarget(overwriteTarget bool) *ConfigBuilder {
	b.overwriteTarget = overwriteTarget
	return b
}

func (b *ConfigBuilder) WithOverwriteTemplate(overwriteTemplate bool) *ConfigBuilder {
	b.overwriteTemplate = overwriteTemplate
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

func (b *ConfigBuilder) WithUpload(isUpload bool) *ConfigBuilder {
	b.isUpload = isUpload
	return b
}

func (b *ConfigBuilder) WithTemplate(isTemplate bool) *ConfigBuilder {
	b.isTemplate = isTemplate
	return b
}

func (b *ConfigBuilder) WithMirror(isMirror bool) *ConfigBuilder {
	b.isMirror = isMirror
	return b
}
