package interactiverunner

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/kube"
)

type KubeNamespaceSelector struct {
	runner InteractiveRunner
}

func NewKubeNamespaceSelector(runner InteractiveRunner) *KubeNamespaceSelector {
	return &KubeNamespaceSelector{runner: runner}
}

func (k *KubeNamespaceSelector) Select() (string, error) {
	kubeClient := kube.NewClient()
	namespaces, err := kubeClient.GetNamespaces()
	if err != nil {
		return "", err
	}

	i, _, err := k.runner.Select(promptui.Select{
		Label:     "Select Kubernetes Namespace: ",
		Items:     namespaces,
		Templates: SelectTemplateBuilder("Kubernetes Namespace", "", ""),
	})
	if err != nil {
		return "", err
	}
	return namespaces[i], nil
}

type KubeSecretSelector struct {
	runner    InteractiveRunner
	namespace string
}

func NewKubeSecretSelector(runner InteractiveRunner, namespace string) *KubeSecretSelector {
	return &KubeSecretSelector{runner: runner, namespace: namespace}
}

func (k *KubeSecretSelector) Select() (string, error) {
	kubeClient := kube.NewClient()
	secrets, err := kubeClient.GetSecrets(k.namespace)
	if err != nil {
		return "", err
	}
	i, _, err := k.runner.Select(promptui.Select{
		Label:     fmt.Sprintf("Select kubernetes secret on namespace %s", k.namespace),
		Items:     secrets,
		Templates: SelectTemplateBuilder("Kubernetes Secret", "", ""),
	})
	if err != nil {
		return "", err
	}
	return secrets[i], nil
}
