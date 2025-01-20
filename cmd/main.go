package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/yammerjp/optruck/pkg/dotenv"
)

/*
optruck: build dotenv/kubernetes secrets from 1password vaults

# filepathをもとに、1passwordのitemを作成する
$ optruck dotenv store --account <account-name> --vault <vault-name> --name <name> --file <file-path>
# 1passwordのitemをもとに、.envファイルを作成する
$ optruck dotenv restore --account <account-name> --vault <vault-name> --name <name> --file <file-path>

# kubernetes secretをもとに、1passwordのitemを作成する。作成後、復元のためのtemplateを出力する
$ optruck kube secret store --account <account-name> --vault <vault-name> --name <name> --file <secret-file-path> --output <template-file-path>

# 1passwordのitemとtemplateをもとに、kubernetes secretを作成する
$ optruck kube secret restore --account <account-name> --vault <vault-name> --name <name>  --output <secret-file-path>
$ optruck kube secret restore --account <account-name> --file=<template-file-path> --output <secret-file-path>

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
	Account string `required:"" help:"1password account name"`
	Vault   string `required:"" help:"1password vault name"`
	Name    string `required:"" help:"1password item name"`
	File    string `default:".env" help:"file path to store"`
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
	return client.StoreFromFile(context.Background(), c.Account, c.Vault, c.Name, c.File)
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
