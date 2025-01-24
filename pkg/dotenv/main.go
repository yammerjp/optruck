package dotenv

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/yammerjp/optruck/pkg/op"
)

type Client struct {
	opClient *op.Client
	logger   *slog.Logger
}

func BuildClient(logger *slog.Logger) *Client {
	opClient := op.BuildClient()
	return &Client{opClient: opClient, logger: logger}
}

func NewClient(opClient *op.Client, logger *slog.Logger) *Client {
	return &Client{opClient: opClient, logger: logger}
}

type envTemplate struct {
	VaultName string
	VaultID   string
	ItemName  string
	ItemID    string

	Fields []struct {
		Label     string
		ID        string
		Value     string
		Concealed bool
	}
}

type uploadableSecret struct {
	AccountName string
	VaultName   string
	ItemName    string
	EnvPairs    map[string]string
}

func (c Client) ValidateUpload(accountName, vaultName, itemName string) error {
	exists, err := c.opClient.CheckItemExists(accountName, vaultName, itemName)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("item already exists")
	}
	return nil
}

func (c *Client) Upload(accountName, vaultName, itemName string, envFilePath string) (*op.ItemCreateResponse, error) {
	uploadableSecret, err := c.buildUploadableSecret(accountName, vaultName, itemName, envFilePath)
	if err != nil {
		return nil, err
	}
	return c.uploadSecret(uploadableSecret)
}

func (c *Client) buildUploadableSecret(accountName, vaultName, itemName, envFilePath string) (*uploadableSecret, error) {
	envPairs, err := c.readEnvFile(envFilePath)
	if err != nil {
		return nil, err
	}
	return &uploadableSecret{AccountName: accountName, VaultName: vaultName, ItemName: itemName, EnvPairs: envPairs}, nil
}

func (c *Client) readEnvFile(filePath string) (map[string]string, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return godotenv.Parse(reader)
}

func (c *Client) uploadSecret(uploadableSecret *uploadableSecret) (*op.ItemCreateResponse, error) {
	return c.opClient.CreateItem(uploadableSecret.AccountName, uploadableSecret.VaultName, uploadableSecret.ItemName, uploadableSecret.EnvPairs)
}

// check file exists
func (c *Client) CheckFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
