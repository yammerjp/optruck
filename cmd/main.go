package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/yammerjp/optruck/pkg/dotenv"
)

/*
optruck: build dotenv files or kubernetes secrets with 1password vaults

Usage:

	optruck [COMMAND] [OPTIONS]

Description:

	This CLI tool is for managing environment variable files (.env), 1Password, and Kubernetes Secrets.
	It allows you to store .env files or Kubernetes Secrets, generate templates from 1Password, and create snapshots for easy restoration.

Commands:

	dotenv              Operations for .env files
	  store             Store the current .env in 1Password.
	  generate          Generate a .env.1password template from 1Password.
	  snapshot          Take a snapshot of the current .env state and make it restorable.

	kube-secret         Operations for Kubernetes Secrets
	  store             Store the current Kubernetes Secret.
	  snapshot          Take a snapshot of the current Kubernetes Secret state and make it restorable.
	  generate          Generate a YAML template for Kubernetes Secrets with placeholders for `op inject`.

	interactive         Launch interactive mode to guide you through the process.

Options:

	--item <ITEM>         Specify the 1Password Item to use.
	                      Default: inferred from the .env file name or Kubernetes Secret name.
	--vault <VAULT>       Specify the 1Password Vault to use.
	                      Default: user's default Vault.
	--account <ACCOUNT>   Specify the 1Password account to use.
	                      Useful when working with multiple accounts.
	--k8s-secret <SECRET> Specify the Kubernetes Secret to use.
	                      Default: the default Kubernetes Secret name.
	--namespace <NAMESPACE> Kubernetes namespace to use.
	                      Default: the current namespace.
	--context <CONTEXT>   Kubernetes context to use.
	                      Default: the current context.
	--template-file <FILE> Specify the file path to the template file.
	                      Default: <name-secret>.1password.yaml
	--env-file <FILE>     Specify the file path to the .env file.
	                      Default: .env
	--help                Show help message.

Examples:

	# Store the current .env in 1Password
	optruck dotenv store --vault development --item my-project-env

	# Store the current Kubernetes Secret
	optruck kube-secret store --k8s-secret my-secret --namespace default

	# Generate a Kubernetes Secret YAML template with placeholders
	optruck kube-secret generate --k8s-secret my-secret --namespace default

	# Create a snapshot of the current .env state
	optruck dotenv snapshot --vault development --item my-project-env

	# Create a snapshot of the current Kubernetes Secret state
	optruck kube-secret snapshot --k8s-secret my-secret --namespace default

	# Launch interactive mode to guide through the process
	optruck interactive

```template-file-path
apiVersion: v1
kind: Secret
metadata:

	name: my-secret
	namespace: default

type: Opaque
data:

	my-key: {{ op://optruck-vault/my-secret/base64-my-key }}

```
*/
type OptruckCmd struct {
	Version kong.VersionFlag
	Dotenv  DotenvCmd `cmd:"" help:"store/restore .env file from 1password item"`
	Kube    KubeCmd   `cmd:"" help:"store/restore kubernetes secret from 1password item"`
}

type DotenvCmd struct {
	Store   DotenvStoreCmd   `cmd:"" help:"store .env file to 1password item"`
	Restore DotenvRestoreCmd `cmd:"" help:"restore .env file from 1password item"`
}

type DotenvStoreCmd struct {
	Account  string `required:"" help:"1password account name"`
	Vault    string `required:"" help:"1password vault name"`
	Name     string `required:"" help:"1password item name"`
	File     string `default:".env" help:"file path to store"`
	Template string `default:".env.1password" help:"file path to output template"`
}

type DotenvRestoreCmd struct {
	Account string `required:"" help:"1password account name"`
	Vault   string `required:"" help:"1password vault name"`
	Name    string `required:"" help:"1password item name"`
	File    string `required:"" help:"file path to restore"`
}

type KubeCmd struct {
	Store   KubeStoreCmd   `cmd:"" help:"store kubernetes secret to 1password item"`
	Restore KubeRestoreCmd `cmd:"" help:"restore kubernetes secret from 1password item"`
}

type KubeStoreCmd struct {
	Account string `required:"" help:"1password account name"`
	Vault   string `required:"" help:"1password vault name"`
	Name    string `required:"" help:"1password item name"`
	File    string `required:"" help:"file path to store"`
	Output  string `required:"" help:"file path to output template"`
}

type KubeRestoreCmd struct {
	Account string `required:"" help:"1password account name"`
	Vault   string `required:"" help:"1password vault name"`
	Name    string `required:"" help:"1password item name"`
	File    string `required:"" help:"file path to store"`
}

func (c *DotenvStoreCmd) Run() error {
	client := dotenv.BuildClient()
	return client.StoreFromFile(context.Background(), c.Account, c.Vault, c.Name, c.File, c.Template)
}

func (c *DotenvRestoreCmd) Run() error {
	client := dotenv.BuildClient()
	return client.RestoreToFile(context.Background(), c.Account, c.Vault, c.Name, c.File)
}

func (c *KubeStoreCmd) Run() error {
	fmt.Println("kube store")
	return nil
}

func (c *KubeRestoreCmd) Run() error {
	fmt.Println("kube restore")
	return nil
}

func Run() {
	ctx := kong.Parse(&OptruckCmd{})
	err := ctx.Run()
	if err != nil {
		fmt.Println(err, os.Stderr)
		os.Exit(1)
	}
}
