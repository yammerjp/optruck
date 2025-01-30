package optruck

import (
	"fmt"
)

// TODO: test

func (cli *CLI) validateConflictOptions() error {
	if len(cli.Item) > 100 {
		return fmt.Errorf("item must be less than 100 characters")
	}
	return nil
}
