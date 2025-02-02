package interactive

import (
	"github.com/manifoldco/promptui"
)

type Runnable interface {
	Select(promptui.Select) (int, string, error)
	Input(promptui.Prompt) (string, error)
}

type RunnableImpl struct{}

func (r *RunnableImpl) Select(prompt promptui.Select) (int, string, error) {
	return prompt.Run()
}

func (r *RunnableImpl) Input(prompt promptui.Prompt) (string, error) {
	return prompt.Run()
}
