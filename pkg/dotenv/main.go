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

func (c *Client) StoreFromFile(ctx context.Context, accountName, vaultName, itemName, inputFilePath string) error {
	reader, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	envPairs, err := godotenv.Parse(reader)
	if err != nil {
		return err
	}

	return c.opClient.CreateItem(ctx, accountName, vaultName, itemName, envPairs)
}

func (c *Client) RestoreToFile(ctx context.Context, accountName, vaultName, itemName, outputFilePath string) error {
	fmt.Println("restore to file", accountName, vaultName, itemName, outputFilePath)
	return nil
}
