package cmd

import "fmt"

const Version = "0.1.0"

// VersionCmd represents the `version` command.
type VersionCmd struct{}

// Run executes the logic for the `version` command.
func (cmd *VersionCmd) Run() error {
	fmt.Printf("optruck version %s\n", Version)
	return nil
}
