package cmd

import (
	"fmt"

	"github.com/yammerjp/optruck/pkg/dotenv"
	"github.com/yammerjp/optruck/pkg/op"
)

// MirrorCmd represents the `mirror` command.
type MirrorCmd struct {
	UploadFlags
	TemplateFlags
	SharedFlags
}

// Run executes the logic for the `mirror` command.
func (cmd *MirrorCmd) Run() error {
	if err := cmd.ValidateInputSource(); err != nil {
		return err
	}

	// kubernetes option is not supported yet
	if cmd.KubeSecret != "" {
		return fmt.Errorf("kubernetes option is not supported yet")
	}
	if cmd.Context != "default" {
		return fmt.Errorf("context option is not supported yet")
	}
	if cmd.Namespace != "default" {
		return fmt.Errorf("namespace option is not supported yet")
	}

	// interactive mode is not supported yet
	if cmd.Interactive {
		return fmt.Errorf("interactive mode is not supported yet")
	}

	dotenvClient := dotenv.BuildClient()
	opClient := op.BuildClient()

	if !cmd.Overwrite {
		if dotenvClient.CheckFileExists(cmd.EnvFile) {
			return fmt.Errorf("file already exists")
		}
		exists, err := opClient.CheckItemExists(cmd.Account, cmd.Vault, cmd.Item)
		if err != nil {
			return fmt.Errorf("failed to check item exists: %w", err)
		}
		if exists {
			return fmt.Errorf("item already exists")
		}
	}

	fmt.Println("Executing 'upload' command...")
	fmt.Printf("EnvFile: %s\n", cmd.EnvFile)
	fmt.Printf("Vault: %s\n", cmd.Vault)
	fmt.Printf("Account: %s\n", cmd.Account)
	fmt.Printf("Item: %s\n", cmd.Item)
	fmt.Printf("Overwrite: %t\n", cmd.Overwrite)
	fmt.Println("--------------------------------")

	// TODO: if overwrite option is true and item already exists, update the item
	resp, err := dotenvClient.Upload(cmd.Account, cmd.Vault, cmd.Item, cmd.EnvFile)
	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}

	fmt.Println("--------------------------------")
	fmt.Println("Uploaded successfully!")
	fmt.Println("--------------------------------")
	fmt.Printf("Vault: %s\n", resp.Vault.Name)
	fmt.Printf("Item: %s\n", resp.Title)
	fmt.Printf("Item ID: %s\n", resp.ID)
	fmt.Println("Fields:")
	for _, field := range resp.Fields {
		fmt.Printf("  %s: %s\n", field.Label, field.Value)
	}
	fmt.Println("--------------------------------")

	envPairs := make(map[string]string)
	for _, field := range resp.Fields {
		envPairs[field.Label] = field.Value
	}

	err = dotenvClient.WriteEnvTemplateFile(cmd.Output, resp)
	if err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	fmt.Println("--------------------------------")
	fmt.Println("Template file created successfully!")
	fmt.Println("--------------------------------")

	return nil
}
