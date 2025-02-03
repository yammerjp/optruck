package interactive

import (
	"errors"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func (r Runner) PromptOutputPath(k8sSecret string) (string, error) {
	result, err := r.Input(promptui.Prompt{
		Label:     "Enter output path: ",
		Validate:  validateOutputPath,
		Templates: PromptTemplateBuilder("Output Path", ""),
		Default:   DefaultOutputPath(k8sSecret),
	})
	if err != nil {
		return "", err
	}

	// Check if file exists
	if _, err := os.Stat(result); err == nil {
		_, result, err := r.Select(promptui.Select{
			Label:     fmt.Sprintf("File %s already exists. Do you want to overwrite it?", result),
			Items:     []string{"overwrite", "cancel"},
			Templates: SelectTemplateBuilder("Overwrite", "", ""),
		})
		if err != nil {
			return "", err
		}
		if result == "cancel" {
			return "", fmt.Errorf("cancelled by user because file %s already exists, please specify another path", result)
		}
	}

	return result, nil
}

func validateOutputPath(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if stat.IsDir() {
		return errors.New("output path is already created as a directory, please specify another path")
	}
	return nil
}

func DefaultOutputPath(k8sSecret string) string {
	if k8sSecret != "" {
		return fmt.Sprintf("%s-secret.yaml.1password", k8sSecret)
	}
	return ".env.1password"
}
