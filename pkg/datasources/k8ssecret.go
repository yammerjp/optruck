package datasources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"k8s.io/utils/exec"
)

var dns1123SubdomainRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)

const maxDNS1123SubdomainLength = 253

func validateDNS1123Subdomain(name string) error {
	if len(name) > maxDNS1123SubdomainLength {
		return fmt.Errorf("must be no more than %d characters", maxDNS1123SubdomainLength)
	}
	if !dns1123SubdomainRegex.MatchString(name) {
		return fmt.Errorf("must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character")
	}
	return nil
}

type K8sSecretSource struct {
	Namespace  string
	SecretName string
	Client     *K8sClient
}

func (s *K8sSecretSource) FetchSecrets() (map[string]string, error) {
	if err := validateDNS1123Subdomain(s.Namespace); err != nil {
		return nil, fmt.Errorf("invalid namespace name: %w", err)
	}
	if err := validateDNS1123Subdomain(s.SecretName); err != nil {
		return nil, fmt.Errorf("invalid secret name: %w", err)
	}

	secrets, err := s.Client.GetSecret(s.Namespace, s.SecretName)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

type K8sClient struct {
	exec.Interface
}

func NewK8sClient() *K8sClient {
	return &K8sClient{exec.New()}
}

func (c *K8sClient) GetSecret(namespace, secretName string) (map[string]string, error) {
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
