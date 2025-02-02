package datasources

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/yammerjp/optruck/pkg/kube"
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
	Client     *kube.Client
	Logger     *slog.Logger
}

func (s *K8sSecretSource) FetchSecrets() (map[string]string, error) {
	s.Logger.Debug("Fetching secrets from Kubernetes", "namespace", s.Namespace, "secret", s.SecretName)

	if err := validateDNS1123Subdomain(s.Namespace); err != nil {
		s.Logger.Error("Invalid namespace name", "namespace", s.Namespace, "error", err)
		return nil, fmt.Errorf("invalid namespace name: %w", err)
	}
	if err := validateDNS1123Subdomain(s.SecretName); err != nil {
		s.Logger.Error("Invalid secret name", "secret", s.SecretName, "error", err)
		return nil, fmt.Errorf("invalid secret name: %w", err)
	}

	secrets, err := s.Client.GetSecret(s.Namespace, s.SecretName)
	if err != nil {
		s.Logger.Error("Failed to get secret from Kubernetes", "namespace", s.Namespace, "secret", s.SecretName, "error", err)
		return nil, err
	}

	s.Logger.Debug("Successfully fetched secrets from Kubernetes", "namespace", s.Namespace, "secret", s.SecretName, "count", len(secrets))
	return secrets, nil
}
