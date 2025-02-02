package interactive

import (
	"github.com/manifoldco/promptui"
)

type Runner interface {
	Select(promptui.Select) (int, string, error)
	Input(promptui.Prompt) (string, error)
}

type RunnerImpl struct{}

func (r *RunnerImpl) Select(prompt promptui.Select) (int, string, error) {
	return prompt.Run()
}

func (r *RunnerImpl) Input(prompt promptui.Prompt) (string, error) {
	return prompt.Run()
}
