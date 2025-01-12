package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

/*
	optruck: build kubernetes secrets from 1password vaults

	// 引数をもとに、1passwordのitemと、commit可能なファイル -secret.optruck.yaml を作成する
	$ optruck create --namespace <namespace> --name <name> --type Opaque --secret <key>=<value>
	// example
	$ optruck create --namespace default --name my-secret --type Opaque --secret my-key=this-is-secret-string
	// =>
	```default.my-secret.optruck.yaml
	apiVersion: v1
	kind: Secret
	metadata:
	  name: my-secret
	  namespace: default
	type: Opaque
	data:
	  my-key: {{ op://optruck-vault/my-secret/base64-my-key }}
	```

	// commit可能な -secret.optruck.yaml をもとに、1passwordのitemを作成する
	$ optruck apply -f <secret-file>

	// secretのかかれたファイルから、1passwordのitemを作成する
	$ optruck create-item -f <secret-file> -v <vault-name>



*/

type OptruckCmd struct {
	Version kong.VersionFlag
	Create  CreateCmd `cmd:"" help:"create kubernetes secret and 1password item"`
}

func (c *OptruckCmd) Run() error {
	// help
	return nil
}

type CreateCmd struct {
	Namespace string            `default:"default" help:"kubernetes namespace for created secret"`
	Name      string            `required:"" help:"kubernetes secret name (will be used as 1password item name)"`
	Type      string            `default:"Opaque" help:"kubernetes secret type"`
	Vault     string            `default:"optruck" help:"1password vault name"`
	Secret    map[string]string `required:"" set:"K=V" help:"kubernetes secret name (will be used as 1password item key and value)"`
}

func (c *CreateCmd) Run() error {
	namespace := c.Namespace
	fmt.Println(namespace)
	fmt.Println(c.Name)
	fmt.Println(c.Type)
	fmt.Println(c.Vault)
	for key, secret := range c.Secret {
		fmt.Println(key, secret)
	}
	fp, err := os.Create(fmt.Sprintf("%s.%s.optruck.yaml", namespace, c.Name))
	if err != nil {
		return err
	}
	defer fp.Close()
	fp.WriteString(fmt.Sprintf("apiVersion: v1\nkind: Secret\nmetadata:\n  name: %s\n  namespace: %s\ntype: %s\ndata:\n", c.Name, namespace, c.Type))
	for key := range c.Secret {
		fp.WriteString(fmt.Sprintf("  %s: {{ op://%s/%s/%s }}\n", key, c.Vault, c.Name, key))
	}

	return nil
	// build cli args for 1password
}

var CLI OptruckCmd

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&CLI)
	if err != nil {
		fmt.Println(err)
	}
}
