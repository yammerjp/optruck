package op

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

/*
{
  "title": "",
  "category": "PASSWORD",
  "fields": [
    {
      "id": "password",
      "type": "CONCEALED",
      "purpose": "PASSWORD",
      "label": "password",
      "password_details": {
        "strength": "TERRIBLE"
      },
      "value": ""
    },
    {
      "id": "notesPlain",
      "type": "STRING",
      "purpose": "NOTES",
      "label": "notesPlain",
      "value": ""
    }
  ]
}

*/

type ItemCreateRequest struct {
	Title    string
	Category string
	Fields   []ItemCreateRequestField
}

type ItemCreateRequestField struct {
	ID      string
	Type    string
	Purpose string
	Label   string
	Value   string
}

func (c *Client) BuildCreateItemRequest(itemName string, envPairs map[string]string) (ItemCreateRequest, error) {
	ret := ItemCreateRequest{
		Title:    itemName,
		Category: "LOGIN",
	}

	for k, v := range envPairs {
		ret.Fields = append(ret.Fields, ItemCreateRequestField{
			ID:    k,
			Type:  "CONCEALED",
			Label: k,
			Value: v,
		})
	}

	return ret, nil
}

/*
{
  "id": "xxxxxxxxxxxxxxxxxxxxxxxxxx",
  "title": "first-item",
  "version": 1,
  "vault": {
    "id": "xxxxxxxxxxxxxxxxxxxxxxxxxx",
    "name": "optruck-development"
  },
  "category": "PASSWORD",
  "created_at": "2025-01-20T14:00:27.841127+09:00",
  "updated_at": "2025-01-20T14:00:27.841127+09:00",
  "additional_information": "Mon Jan 20 14:00:27 JST 2025",
  "fields": [
    {
      "id": "password",
      "type": "CONCEALED",
      "purpose": "PASSWORD",
      "label": "password",
      "value": "BAR",
      "reference": "op://optruck-development/first-item/password",
      "password_details": {
        "strength": "TERRIBLE"
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

type ItemCreateResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Version int    `json:"version"`
	Vault   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"vault"`
	Category              string                    `json:"category"`
	CreatedAt             string                    `json:"created_at"`
	UpdatedAt             string                    `json:"updated_at"`
	AdditionalInformation string                    `json:"additional_information"`
	Fields                []ItemCreateResponseField `json:"fields"`
}

type ItemCreateResponseField struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	Purpose         string `json:"purpose"`
	Label           string `json:"label"`
	Value           string `json:"value"`
	Reference       string `json:"reference"`
	PasswordDetails struct {
		Strength string `json:"strength"`
	} `json:"password_details"`
}

func (c *Client) CreateItem(accountName, vaultName, itemName string, envPairs map[string]string) (*ItemCreateResponse, error) {
	req, err := c.BuildCreateItemRequest(itemName, envPairs)
	if err != nil {
		return nil, err
	}
	return c.CreateItemByRequest(accountName, vaultName, req)
}

func (c *Client) CreateItemByRequest(accountName, vaultName string, req ItemCreateRequest) (*ItemCreateResponse, error) {
	reqStr, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	cmdArgs := []string{"item", "create", "--format", "json"}

	if accountName != "" {
		cmdArgs = append(cmdArgs, "--account", accountName)
	}
	if vaultName != "" {
		cmdArgs = append(cmdArgs, "--vault", vaultName)
	}

	cmd := exec.Command("op", cmdArgs...)
	cmd.Stdin = bytes.NewBuffer(reqStr)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	os.Stderr.Write(stderr.Bytes())
	if err != nil {
		return nil, err
	}

	var resp ItemCreateResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
