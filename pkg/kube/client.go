package kube

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

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
		return nil, fmt.Errorf("failed to get secret with `$ kubectl get secret -n %s %s -o jsonpath={.data}`: %w", namespace, secretName, err)
	}
	return secrets, nil
}

func (c *Client) GetNamespaces() ([]string, error) {
	cmd := utilExec.NewCommand("kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	stdout := &bytes.Buffer{}
	if err := cmd.Run(nil, stdout); err != nil {
		return nil, fmt.Errorf("failed to get namespaces with `$ kubectl get namespaces -o jsonpath={.items[*].metadata.name}`: %w", err)
	}
	output := stdout.String()
	if output == "" {
		return nil, errors.New("no namespaces found")
	}
	return strings.Split(output, " "), nil
}

func (c *Client) GetSecrets(namespace string) ([]string, error) {
	cmd := utilExec.NewCommand("kubectl", "get", "secrets", "-n", namespace, "--field-selector", "type=Opaque", "-o", "jsonpath={.items[*].metadata.name}")
	stdout := &bytes.Buffer{}
	if err := cmd.Run(nil, stdout); err != nil {
		return nil, fmt.Errorf("failed to get secrets with `$ kubectl get secrets -n %s --field-selector type=Opaque -o jsonpath={.items[*].metadata.name}`: %w", namespace, err)
	}
	output := stdout.String()
	if output == "" {
		return nil, errors.New("no secrets found")
	}
	return strings.Split(output, " "), nil
}
