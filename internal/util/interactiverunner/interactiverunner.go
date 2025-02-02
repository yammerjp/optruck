package interactiverunner

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func SelectTemplateBuilder(selectedPrefix string, mainField string, subField string) *promptui.SelectTemplates {
	active := fmt.Sprintf("▸ {{ .%s | cyan | underline }}", mainField)
	if subField != "" {
		active += fmt.Sprintf(` {{"("|faint}}{{ .%s | red | underline }}{{")"|faint}}`, subField)
	}

	inactive := fmt.Sprintf("  {{ .%s | cyan }}", mainField)
	if subField != "" {
		inactive += fmt.Sprintf(` {{"("|faint}}{{ .%s | red }}{{")"|faint}}`, subField)
	}

	selected := fmt.Sprintf(`{{ "✔" | green }} %-20s: {{ .%s }}`, selectedPrefix, mainField)
	if subField != "" {
		selected += fmt.Sprintf(` {{"("|faint}}{{ .%s }}{{")"|faint}}`, subField)
	}

	return &promptui.SelectTemplates{
		Label:    `{{ . | yellow }}`,
		Active:   active,
		Inactive: inactive,
		Selected: selected,
	}
}

func PromptTemplateBuilder(successPrefix string, mainField string) *promptui.PromptTemplates {
	return &promptui.PromptTemplates{
		Prompt:  `{{ . | yellow }}`,
		Valid:   fmt.Sprintf(`{{ "✔" | green }} {{ .%s | yellow }}`, mainField),
		Invalid: fmt.Sprintf(`{{ "✘" | red }} {{ .%s | yellow }}`, mainField),
		Success: fmt.Sprintf(`{{ "✔" | green }} %-20s: `, successPrefix),
	}
}

type InteractiveRunner interface {
	Select(promptui.Select) (int, string, error)
	Input(promptui.Prompt) (string, error)
}

type InteractiveRunnerImpl struct{}

func (r *InteractiveRunnerImpl) Select(prompt promptui.Select) (int, string, error) {
	return prompt.Run()
}

func (r *InteractiveRunnerImpl) Input(prompt promptui.Prompt) (string, error) {
	return prompt.Run()
}
