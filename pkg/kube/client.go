package kube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"k8s.io/utils/exec"
)

type Client struct {
	exec.Interface
}

func NewClient(exec exec.Interface) *Client {
	return &Client{exec}
}

func (c *Client) GetSecret(namespace, secretName string) (map[string]string, error) {
	// TODO: use k8s.io/client-go
	cmd := c.Command("kubectl", "get", "secret", "-n", namespace, secretName, "-o", "jsonpath={.data}")
	// ex: {"AWS_ACCESS_KEY_ID":"YWJjZGVmZ2hpamtsbW5vcA==","AWS_SECRET_ACCESS_KEY":"YWJjZGVmZ2hpamtsbW5vcA=="}
	stdout := &bytes.Buffer{}
	cmd.SetStdout(stdout)
	cmd.SetStderr(os.Stderr)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	secrets := make(map[string]string)
	if err := json.Unmarshal(stdout.Bytes(), &secrets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}
	return secrets, nil
}

func (c *Client) GetNamespaces() ([]string, error) {
	cmd := c.Command("kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	stdout := &bytes.Buffer{}
	cmd.SetStdout(stdout)
	cmd.SetStderr(os.Stderr)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get namespaces: %w", err)
	}
	namespaces := strings.Split(stdout.String(), " ")
	return namespaces, nil
}

func (c *Client) GetSecrets(namespace string) ([]string, error) {
	cmd := c.Command("kubectl", "get", "secrets", "-n", namespace, "--field-selector", "type=Opaque", "-o", "jsonpath={.items[*].metadata.name}")
	stdout := &bytes.Buffer{}
	cmd.SetStdout(stdout)
	cmd.SetStderr(os.Stderr)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to get secrets: %w", err)
	}
	return strings.Split(stdout.String(), " "), nil
}
