package kube

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yammerjp/optruck/internal/errors"
	utilExec "github.com/yammerjp/optruck/internal/util/exec"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetSecret(namespace, secretName string) (map[string]string, error) {
	// TODO: use k8s.io/client-go
	cmd := utilExec.NewCommand("kubectl", "get", "secret", "-n", namespace, secretName, "-o", "jsonpath={.data}")
	// ex: {"AWS_ACCESS_KEY_ID":"YWJjZGVmZ2hpamtsbW5vcA==","AWS_SECRET_ACCESS_KEY":"YWJjZGVmZ2hpamtsbW5vcA=="}
	secrets := make(map[string]string)
	if err := cmd.RunWithJSON(nil, &secrets); err != nil {
		return nil, errors.NewOperationFailedError(
			fmt.Sprintf("シークレット '%s' の取得", secretName),
			err,
			"kubectlが正しく設定されているか、指定したネームスペースとシークレットが存在するか確認してください",
		)
	}
	return secrets, nil
}

func (c *Client) GetNamespaces() ([]string, error) {
	cmd := utilExec.NewCommand("kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	stdout := &bytes.Buffer{}
	if err := cmd.Run(nil, stdout); err != nil {
		return nil, errors.NewOperationFailedError(
			"Kubernetesネームスペースの一覧取得",
			err,
			"kubectlが正しく設定されているか確認してください",
		)
	}
	output := stdout.String()
	return strings.Split(output, " "), nil
}

func (c *Client) GetSecrets(namespace string) ([]string, error) {
	cmd := utilExec.NewCommand("kubectl", "get", "secrets", "-n", namespace, "--field-selector", "type=Opaque", "-o", "jsonpath={.items[*].metadata.name}")
	stdout := &bytes.Buffer{}
	if err := cmd.Run(nil, stdout); err != nil {
		return nil, errors.NewOperationFailedError(
			fmt.Sprintf("ネームスペース '%s' のシークレット一覧取得", namespace),
			err,
			"kubectlが正しく設定されているか、指定したネームスペースが存在するか確認してください",
		)
	}
	output := stdout.String()
	return strings.Split(output, " "), nil
}
