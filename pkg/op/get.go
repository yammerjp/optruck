package op

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

/*
$ op item get xxxx --format=json
{
  "id": "xxxx",
  "title": "first-item",
  "version": 1,
  "vault": {
    "id": "vvvv",
    "name": "optruck-development"
  },
  "category": "PASSWORD",
  "created_at": "2025-01-20T23:29:41.160193+09:00",
  "updated_at": "2025-01-20T23:29:41.160194+09:00",
  "additional_information": "Mon Jan 20 23:29:41 JST 2025",
  "fields": [
    {
      "id": "password",
      "type": "CONCEALED",
      "purpose": "PASSWORD",
      "label": "password",
      "value": "bar",
      "reference": "op://optruck-development/first-item/password",
      "password_details": {
        "strength": "TERRIBLE",
        "history": ["baz"]
      }
    },
    {
      "id": "notesPlain",
      "type": "STRING",
      "purpose": "NOTES",
      "label": "notesPlain",
      "reference": "op://optruck-development/first-item/notesPlain"
    }
  ]
}
*/

var ErrMoreThanOneItemMatches = errors.New("more than one item matches")

type GetItemResponse struct {
	ID    string `json:"id"`
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

func (c *Client) GetItem(ctx context.Context, accountName, vaultName, itemName string) (*GetItemResponse, error) {
	cmd := exec.Command("op", "get", "item", itemName, "--vault", vaultName, "--account", accountName)
	stdoutBuffer := bytes.NewBuffer(nil)
	stderrBuffer := bytes.NewBuffer(nil)
	cmd.Stdout = stdoutBuffer
	cmd.Stderr = stderrBuffer

	var resp GetItemResponse
	if err := cmd.Run(); err != nil {
		if strings.Contains(stderrBuffer.String(), " isn't an item. Specify the item with its UUID, name, or domain.") {
			// not found
			return nil, nil
		}
		if strings.Contains(stderrBuffer.String(), " More than one item matches ") {
			return nil, ErrMoreThanOneItemMatches
		}
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if err := json.Unmarshal(stdoutBuffer.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
