package output

import (
	"fmt"
	"os"

	"github.com/yammerjp/optruck/pkg/op"
)

func (c *Client) Write(resp *op.ItemCreateResponse, accountName string) error {
	// open
	envTemplateFile, err := os.OpenFile(c.EnvTemplateFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open env template file: %v", err)
	}
	defer envTemplateFile.Close()

	fmt.Fprintf(envTemplateFile, "#   - 1password vault: %s\n", resp.Vault.Name)
	fmt.Fprintf(envTemplateFile, "#   - 1password account: %s\n", accountName)
	fmt.Fprintf(envTemplateFile, "#   - 1password item: %s\n", resp.Title)
	fmt.Fprintf(envTemplateFile, "# To restore, run the following command:\n")
	fmt.Fprintf(envTemplateFile, "#   $ cat .env.1password | grep -v '^#' | op inject > .env\n")

	for _, field := range resp.Fields {
		if field.Purpose != "" {
			continue
		}
		if field.Type == "CONCEALED" {
			fmt.Fprintf(envTemplateFile, "%s={{op://%s/%s/%s}}\n", field.Label, resp.Vault.Name, resp.Title, field.ID)
		} else {
			fmt.Fprintf(envTemplateFile, "%s=%s\n", field.Label, field.Value)
		}
	}

}
