package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/alecthomas/kong"
	"github.com/yammerjp/optruck/internal/op"
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

type GetCmd struct {
	Vault string `required:"" help:"1password vault name"`
	ID    string `required:"" help:"1password item id"`
}

func (c *GetCmd) Run() error {
	itemRequest := op.ItemRequest{
		Vault: c.Vault,
		ID:    c.ID,
	}
	item, err := itemRequest.GetItem()
	if err != nil {
		return err
	}
	fmt.Println(item)
	return nil
}

type OptruckCmd struct {
	Version kong.VersionFlag
	Create  CreateCmd `cmd:"" help:"create kubernetes secret and 1password item"`
	Get     GetCmd    `cmd:"" help:"get 1password item"`
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
	// TODO: validate arguments
	// TODO: create 1password item
	if err := c.createMaskedSecretYaml(); err != nil {
		return err
	}

	// check if 1password item already exists
	cmd := exec.Command("op", "item", "create", "--vault", c.Vault)
	tmpl, err := c.createSecretTemplateJson()
	if err != nil {
		return err
	}
	cmd.Stdin = tmpl
	// stdoutとstderrはそのまま出力
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (c *CreateCmd) createSecretTemplateJson() (*bytes.Buffer, error) {
	var buf bytes.Buffer

	var secretTemplate struct {
		Title    string `json:"title"`
		Category string `json:"category"`
		Fields   []struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Purpose  string `json:"purpose"`
			Label    string `json:"label"`
			Password string `json:"password"`
		} `json:"fields"`
	}

	secretTemplate.Title = c.Name
	secretTemplate.Category = "PASSWORD"
	for key, value := range c.Secret {
		secretTemplate.Fields = append(secretTemplate.Fields, struct {
			ID       string `json:"id"`
			Type     string `json:"type"`
			Purpose  string `json:"purpose"`
			Label    string `json:"label"`
			Password string `json:"password"`
		}{
			ID:       key,
			Type:     "CONCEALED",
			Purpose:  "PASSWORD",
			Label:    key,
			Password: value,
		})
	}

	err := json.NewEncoder(&buf).Encode(secretTemplate)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func (c *CreateCmd) createMaskedSecretYaml() error {
	fp, err := os.Create(fmt.Sprintf("%s.%s.optruck.yaml", c.Namespace, c.Name))
	if err != nil {
		return err
	}
	defer fp.Close()
	fp.WriteString(fmt.Sprintf("apiVersion: v1\nkind: Secret\nmetadata:\n  name: %s\n  namespace: %s\ntype: %s\ndata:\n", c.Name, c.Namespace, c.Type))
	for key := range c.Secret {
		fp.WriteString(fmt.Sprintf("  %s: {{ op://%s/%s/%s }}\n", key, c.Vault, c.Name, key))
	}

	return nil
}

func (c *CreateCmd) createItemCommandString() (string, error) {
	return fmt.Sprintf("op item create --vault %s", c.Vault), nil
}

var CLI OptruckCmd

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&CLI)
	if err != nil {
		fmt.Println(err)
	}
}
