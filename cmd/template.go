package cmd

import "fmt"

// TemplateFlags defines flags specific to the template command
type TemplateFlags struct {
	Output string `help:"Path to save the generated template file." default:".env.1password.tpl"`
}

// TemplateCmd represents the `template` command.
type TemplateCmd struct {
	TemplateFlags
	SharedFlags
}

// Run executes the logic for the `template` command.
func (cmd *TemplateCmd) Run() error {
	if cmd.Namespace != "default" {
		return fmt.Errorf("namespace option is not supported yet")
	}

	fmt.Println("Executing 'template' command...")

	fmt.Printf("Output: %s\n", cmd.Output)
	fmt.Printf("Account: %s\n", cmd.Account)
	fmt.Printf("Item: %s\n", cmd.Item)
	fmt.Printf("Overwrite: %t\n", cmd.Overwrite)

	// TODO fetch item from 1Password

	return fmt.Errorf("not implemented")
}
