package datasources

import (
	"fmt"
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
}

func (s *K8sSecretSource) FetchSecrets() (map[string]string, error) {
	if err := validateDNS1123Subdomain(s.Namespace); err != nil {
		return nil, fmt.Errorf("invalid namespace name, please specify a valid namespace name with --k8s-namespace option: %w", err)
	}
	if err := validateDNS1123Subdomain(s.SecretName); err != nil {
		return nil, fmt.Errorf("invalid secret name, please specify a valid secret name with --k8s-secret option: %w", err)
	}

	secrets, err := s.Client.GetSecret(s.Namespace, s.SecretName)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}
