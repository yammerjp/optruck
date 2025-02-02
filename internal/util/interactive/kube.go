package interactive

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/kube"
)

func (r Runner) SelectKubeNamespace() (string, error) {
	kubeClient := kube.NewClient()
	namespaces, err := kubeClient.GetNamespaces()
	if err != nil {
		return "", err
	}

	i, _, err := r.Select(promptui.Select{
		Label:     "Select Kubernetes Namespace: ",
		Items:     namespaces,
		Templates: SelectTemplateBuilder("Kubernetes Namespace", "", ""),
	})
	if err != nil {
		return "", err
	}
	return namespaces[i], nil
}

func (r Runner) SelectKubeSecret(namespace string) (string, error) {
	kubeClient := kube.NewClient()
	secrets, err := kubeClient.GetSecrets(namespace)
	if err != nil {
		return "", err
	}
	i, _, err := r.Select(promptui.Select{
		Label:     fmt.Sprintf("Select kubernetes secret on namespace %s", namespace),
		Items:     secrets,
		Templates: SelectTemplateBuilder("Kubernetes Secret", "", ""),
	})
	if err != nil {
		return "", err
	}
	return secrets[i], nil
}
