package dotenv

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/yammerjp/optruck/pkg/op"
)

type Client struct {
	opClient *op.Client
}

func BuildClient() *Client {
	opClient := op.NewClient()
	return &Client{opClient: opClient}
}

func NewClient(opClient *op.Client) *Client {
	return &Client{opClient: opClient}
}

func (c *Client) readEnvFile(filePath string) (map[string]string, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return godotenv.Parse(reader)
}

func (c *Client) writeEnvTemplateFile(filePath string, resp *op.ItemCreateResponse) error {
	writer, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer writer.Close()

	fmt.Fprintf(writer, "# This file is generated by optruck.\n")
	fmt.Fprintf(writer, "#   - 1password vault: %s\n", resp.Vault.Name)
	fmt.Fprintf(writer, "#   - 1password item: %s\n", resp.Title)
	fmt.Fprintf(writer, "# To restore, run the following command:\n")
	fmt.Fprintf(writer, "#   $ cat .env.1password | grep -v '^#' | op inject > .env\n")

	for _, field := range resp.Fields {
		if field.Purpose == "" {
			if field.Type == "STRING" {
				fmt.Fprintf(writer, "%s=%s\n", field.Label, field.Value)
			}
			if field.Type == "CONCEALED" {
				fmt.Fprintf(writer, "%s={{op://%s/%s/%s}}\n", field.Label, resp.Vault.ID, resp.ID, field.ID)
			}
		}
	}

	return nil
}

func (c *Client) StoreFromFile(ctx context.Context, accountName, vaultName, itemName, inputFilePath, templateFilePath string) error {
	envPairs, err := c.readEnvFile(inputFilePath)
	if err != nil {
		return err
	}

	resp, err := c.opClient.CreateItem(ctx, accountName, vaultName, itemName, envPairs)
	if err != nil {
		return err
	}

	return c.writeEnvTemplateFile(templateFilePath, resp)
}

func (c *Client) RestoreToFile(ctx context.Context, accountName, vaultName, itemName, outputFilePath string) error {
	fmt.Println("restore to file", accountName, vaultName, itemName, outputFilePath)
	return nil
}
