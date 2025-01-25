package op

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var ErrMoreThanOneItemMatches = errors.New("more than one item matches")
var ErrItemNotFound = errors.New("item not found")

type GetItemResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Vault struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"vault"`
	Fields []struct {
		ID        string `json:"id"`
		Label     string `json:"label"`
		Purpose   string `json:"purpose"`
		Reference string `json:"reference"`
	} `json:"fields"`
}

func (c *Client) GetItem() (*SecretReference, error) {
	cmd := c.BuildItemCommand("get", c.Target.ItemName)
	stdoutBuffer := bytes.NewBuffer(nil)
	stderrBuffer := bytes.NewBuffer(nil)
	cmd.SetStdout(stdoutBuffer)
	cmd.SetStderr(stderrBuffer)

	var resp GetItemResponse
	if err := cmd.Run(); err != nil {
		if strings.Contains(stderrBuffer.String(), " isn't an item. Specify the item with its UUID, name, or domain.") {
			return nil, ErrItemNotFound
		}
		if strings.Contains(stderrBuffer.String(), " More than one item matches ") {
			return nil, ErrMoreThanOneItemMatches
		}
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if err := json.Unmarshal(stdoutBuffer.Bytes(), &resp); err != nil {
		return nil, err
	}

	fieldLabels := []string{}
	for _, field := range resp.Fields {
		if field.Purpose == "" {
			fieldLabels = append(fieldLabels, field.Label)
		}
	}
	return &SecretReference{
		Account:     c.Target.Account,
		VaultName:   resp.Vault.Name,
		VaultID:     resp.Vault.ID,
		ItemName:    resp.Title,
		ItemID:      resp.ID,
		FieldLabels: fieldLabels,
	}, nil
}
