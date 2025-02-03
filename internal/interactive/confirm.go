package interactive

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func (r Runner) Confirm(cmds []string) error {
	fmt.Println("The selected options are same as below.")
	fmt.Print("    $ optruck")
	for _, cmd := range cmds {
		if strings.HasPrefix(cmd, "--") {
			// break line
			fmt.Printf(" \\\n      %s", cmd)
		} else {
			fmt.Printf(" %s", cmd)
		}
	}
	fmt.Println()

	i, _, err := r.Select(promptui.Select{
		Label:     "Do you want to proceed? (yes/no)",
		Items:     []string{"yes", "no"},
		Templates: SelectTemplateBuilder("Do you want to proceed?", "", ""),
	})
	if err != nil {
		return err
	}
	if i != 0 {
		return fmt.Errorf("aborted")
	}
	return nil
}
