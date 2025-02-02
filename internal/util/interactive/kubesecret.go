package interactive

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/yammerjp/optruck/pkg/kube"
)

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
