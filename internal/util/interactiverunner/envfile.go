package interactiverunner

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

const (
	DefaultEnvFilePath = ".env"
)

type EnvFilePrompter struct {
	runner InteractiveRunner
}

func NewEnvFilePrompter(runner InteractiveRunner) *EnvFilePrompter {
	return &EnvFilePrompter{runner: runner}
}

func (e *EnvFilePrompter) Prompt() (string, error) {
	result, err := e.runner.Input(promptui.Prompt{
		Label:   "Enter env file path: ",
		Default: DefaultEnvFilePath,
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("env file path is required")
			}
			stat, err := os.Stat(input)
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			if stat.IsDir() {
				return fmt.Errorf("env file path is already created as a directory")
			}
			return nil
		},
		Templates: PromptTemplateBuilder("Env File Path", ""),
	})
	if err != nil {
		return "", err
	}
	return result, nil

}
