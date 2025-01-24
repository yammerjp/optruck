package output

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yammerjp/optruck/pkg/op"
)

type Client struct {
	AccountID           string
	VaultID             string
	ItemID              string
	ItemName            string
	logger              *slog.Logger
	EnvTemplateFilePath string
}

func (c *Client) Print(resp *op.ItemCreateResponse) {
	// open
	envTemplateFile, err := os.OpenFile(c.EnvTemplateFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to open env template file: %v", err))
		return
	}
	defer envTemplateFile.Close()

	// write
	fmt.Fprintf(envTemplateFile, "#   - 1password item: %s\n", c.ItemName)
	fmt.Fprintf(envTemplateFile, "# To restore, run the following command:\n")
	fmt.Fprintf(envTemplateFile, "#   $ cat .env.1password | grep -v '^#' | op inject > .env\n")

	for _, field := range resp.Fields {
		if field.Purpose != "" {
			continue
		}
		if field.Type == "CONCEALED" {
			fmt.Fprintf(envTemplateFile, "%s={{op://%s/%s/%s}}\n", field.Label, c.VaultID, c.ItemID, field.ID)
		} else {
			fmt.Fprintf(envTemplateFile, "%s=%s\n", field.Label, field.Value)
		}
	}

}
