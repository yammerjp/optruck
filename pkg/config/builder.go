package config

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/actions"
	"github.com/yammerjp/optruck/pkg/datasources"
	"github.com/yammerjp/optruck/pkg/kube"
	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

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
	if b.overwriteTarget && b.overwrite {
		return fmt.Errorf("cannot use both --overwrite-target and --overwrite")
	}
	if b.overwriteTemplate && b.overwrite {
		return fmt.Errorf("cannot use both --overwrite-template and --overwrite")
	}
	if b.overwriteTarget && b.isTemplate {
		return fmt.Errorf("cannot use --overwrite-target on template action")
	}
	if b.overwriteTemplate && b.isUpload {
		return fmt.Errorf("cannot use --overwrite-template on upload action")
	}

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
			Client:     kube.NewClient(),
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

const defaultEnvFilePath = ".env"
const defaultOutputFormat = "env"
const defaultOutputPathOnEnv = ".env.1password"

func defaultOutputPathOnK8s(item string) string {
	return fmt.Sprintf("%s-secret.yaml.1password", item)
}

// TODO: create new file for default value setting
func (b *ConfigBuilder) SetDefaultIfEmpty() error {
	if !b.isUpload && !b.isTemplate && !b.isMirror {
		b.isMirror = true
	}
	if b.envFile == "" && b.k8sSecret == "" {
		b.envFile = defaultEnvFilePath
	}
	if b.k8sSecret != "" && b.k8sNamespace == "" {
		return fmt.Errorf("k8s namespace is required when using k8s secret")
	}
	if b.outputFormat == "" {
		b.outputFormat = defaultOutputFormat
	}
	if b.output == "" {
		if b.outputFormat == "env" {
			b.output = defaultOutputPathOnEnv
		} else {
			b.output = defaultOutputPathOnK8s(b.item)
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
		}
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
		}
	}
	return nil
}

// TODO: create new file for interactive mode
func (b *ConfigBuilder) SetConfigInteractively() error {
	// FIXME: check to validate before interactive

	if err := b.setActionInteractively(); err != nil {
		return err
	}
	fmt.Println("action", b.isUpload, b.isTemplate, b.isMirror)
	if err := b.setDataSourceInteractively(); err != nil {
		return err
	}
	if err := b.setTargetInteractively(); err != nil {
		return err
	}
	if err := b.setDestInteractively(); err != nil {
		return err
	}
	// TODO: print result with command line option format
	// ex:
	//    The selected options are same as below.
	//        $ optruck --upload --env-file .env --output .env.1password --vault Development --account my.1password.com --item my-item
	return nil
}

func (b *ConfigBuilder) setActionInteractively() error {
	if b.isUpload || b.isTemplate || b.isMirror {
		// already set
		return nil
	}

	prompt := promptui.Select{
		Label:     "Select action",
		Items:     []string{"upload", "template", "mirror"},
		CursorPos: 2,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	b.isUpload = result == "upload"
	b.isTemplate = result == "template"
	b.isMirror = result == "mirror"
	return nil
}

func (b *ConfigBuilder) setDataSourceInteractively() error {
	if b.isTemplate {
		// data source is not needed
		return nil
	}
	if b.envFile != "" || b.k8sSecret != "" {
		// already set
		return nil
	}
	prompt := promptui.Select{
		Label:     "Select data source",
		Items:     []string{"env file", "k8s secret"},
		CursorPos: 2,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return err
	}
	switch result {
	case "env file":
		// TODO: create new function setEnvFilePathInteractively()
		prompt := promptui.Prompt{
			Label:   "Enter env file path",
			Default: defaultEnvFilePath,
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("env file path is required")
				}
				if _, err := os.Stat(input); err != nil {
					return fmt.Errorf("env file does not exist")
				}
				return nil
			},
		}
		result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.envFile = result
	case "k8s secret":
		// TODO: create new function setK8sSecretInteractively()
		kubeClient := kube.NewClient()
		namespaces, err := kubeClient.GetNamespaces()
		if err != nil {
			return err
		}

		if b.k8sNamespace == "" {
			prompt := promptui.Select{
				Label: "Select kubernetes namespace",
				Items: namespaces,
			}
			_, result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.k8sNamespace = result
		}

		secrets, err := kubeClient.GetSecrets(b.k8sNamespace)
		if err != nil {
			return err
		}
		prompt = promptui.Select{
			Label: fmt.Sprintf("Select kubernetes secret on namespace %s", b.k8sNamespace),
			Items: secrets,
		}
		_, result, err = prompt.Run()
		if err != nil {
			return err
		}
		b.k8sSecret = result
	default:
		return fmt.Errorf("invalid data source: %s", result)
	}
	return nil
}

func (b *ConfigBuilder) setTargetInteractively() error {
	opClient := b.buildOpTarget().BuildClient()

	// TODO: create new function setAccountInteractively()
	if b.account == "" {
		accounts, err := opClient.ListAccounts()
		if err != nil {
			return err
		}
		accountNames := make([]string, len(accounts))
		for i, account := range accounts {
			accountNames[i] = account.URL
		}
		prompt := promptui.Select{
			Label: "Select account",
			Items: accountNames,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.account = result
	}

	// TODO: create new function setVaultInteractively()
	if b.vault == "" {
		vaults, err := opClient.ListVaults()
		if err != nil {
			return err
		}
		vaultNames := make([]string, len(vaults))
		for i, vault := range vaults {
			vaultNames[i] = vault.Name
		}
		prompt := promptui.Select{
			Label: "Select vault",
			Items: vaultNames,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.vault = result
	}

	// TODO: create new function setItemInteractively()
	if !b.overwriteTarget && !b.overwrite {
		prompt := promptui.Select{
			Label: "Select overwrite mode",
			Items: []string{"overwrite existing", "create new"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		switch result {
		case "overwrite existing":
			b.overwriteTarget = true
		case "create new":
			b.overwriteTarget = false
		default:
			return fmt.Errorf("invalid selection: %s", result)
		}
	}

	if b.item == "" {
		items, err := opClient.ListItems()
		if err != nil {
			return err
		}
		itemNames := make([]string, len(items))
		for i, item := range items {
			itemNames[i] = item.ItemName
		}
		if b.overwrite || b.overwriteTarget {
			b.overwrite = true
			prompt := promptui.Select{
				Label: "Select item name",
				Items: itemNames,
			}
			_, result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.item = result
		} else {
			b.overwrite = false
			defaultName := ""
			// TODO: define default item name format
			if b.envFile != "" {
				defaultName = fmt.Sprintf("dotenv_%s", filepath.Base(filepath.Dir(b.envFile)))
			} else if b.k8sSecret != "" {
				defaultName = fmt.Sprintf("kubernetes_secret_%s_%s", b.k8sNamespace, b.k8sSecret)
			}
			prompt := promptui.Prompt{
				Label:   "Enter item name",
				Default: defaultName,
				Validate: func(input string) error {
					// TODO: define item name format
					if input == "" {
						return fmt.Errorf("item name is required")
					}
					if len(input) > 100 {
						return fmt.Errorf("item name must be less than 100 characters")
					}
					if strings.Contains(input, " ") {
						return fmt.Errorf("item name must not contain spaces")
					}
					for _, n := range itemNames {
						if n == input {
							return fmt.Errorf("item name must be unique")
						}
					}
					return nil
				},
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.item = result
		}
	}
	return nil
}

func (b *ConfigBuilder) setDestInteractively() error {
	if b.isUpload {
		// upload action does not need dest
		return nil
	}

	// TODO: create new function setOutputFormatInteractively()
	// set format if not set
	if b.isTemplate && b.outputFormat == "" {
		// if upload or mirror, detect format from data source
		prompt := promptui.Select{
			Label: "Select output format",
			Items: []string{"env", "k8s"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		b.outputFormat = result
	}

	// TODO: create new function setOutputPathInteractively()
	// FIXME: Having user select overwrite mode first is not user-friendly and should be changed
	// set overwrite if not set
	if !b.overwrite && !b.overwriteTemplate {
		prompt := promptui.Select{
			Label: "Select template overwrite mode",
			Items: []string{"overwrite existing", "create new"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		switch result {
		case "overwrite existing":
			b.overwriteTemplate = true
		case "create new":
			b.overwriteTemplate = false
		default:
			return fmt.Errorf("invalid selection: %s", result)
		}
	}

	defaultOutputPath := ""
	if b.outputFormat == "env" {
		defaultOutputPath = defaultOutputPathOnEnv
	} else {
		defaultOutputPath = defaultOutputPathOnK8s(b.item)
	}

	// set output path if not set
	if b.output == "" {
		if b.overwrite || b.overwriteTemplate {
			// validate existing file
			prompt := promptui.Prompt{
				Label:   "Enter output path",
				Default: defaultOutputPath,
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("output path is required")
					}
					if _, err := os.Stat(input); err != nil {
						return fmt.Errorf("output path does not exist")
					}
					return nil
				},
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.output = result
		} else {
			// validate new file
			prompt := promptui.Prompt{
				Label:   "Enter output path",
				Default: defaultOutputPath,
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("output path is required")
					}
					if _, err := os.Stat(input); err == nil {
						return fmt.Errorf("output path already exists")
					}
					return nil
				},
			}
			result, err := prompt.Run()
			if err != nil {
				return err
			}
			b.output = result
		}
	}
	return nil
}
