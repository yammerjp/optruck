package kube

import (
	"fmt"
	"strings"

	optruckexec "github.com/yammerjp/optruck/pkg/exec"
)

type Client struct {
	optruckexec.CommandConfig
}

func NewClient(config optruckexec.CommandConfig) *Client {
	return &Client{
		CommandConfig: config,
	}
}

func (c *Client) GetSecret(namespace, secretName string) (map[string]string, error) {
	// TODO: use k8s.io/client-go

	command := c.Command("kubectl", "get", "secret", "-n", namespace, secretName, "-o", "jsonpath={.data}")
	// ex: {"AWS_ACCESS_KEY_ID":"YWJjZGVmZ2hpamtsbW5vcA==","AWS_SECRET_ACCESS_KEY":"YWJjZGVmZ2hpamtsbW5vcA=="}

	var secrets map[string]string
	err := command.RunWithJson(nil, &secrets)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	return secrets, nil
}

func (c *Client) GetNamespaces() ([]string, error) {
	command := c.Command("kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	stdout, err := command.Run(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}

	return strings.Split(stdout, " "), nil
}

func (c *Client) GetSecrets(namespace string) ([]string, error) {
	command := c.Command("kubectl", "get", "secrets", "-n", namespace, "--field-selector", "type=Opaque", "-o", "jsonpath={.items[*].metadata.name}")
	stdout, err := command.Run(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets: %w", err)
	}

	return strings.Split(stdout, " "), nil
}
