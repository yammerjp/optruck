package op

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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

func (c *Client) BuildCreateItemRequest(envPairs map[string]string) (ItemCreateRequest, error) {
	ret := ItemCreateRequest{
		Title:    c.Target.ItemName,
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

type SecretResponse struct {
	AccountName string
	VaultName   string
	VaultID     string
	ItemName    string
	ItemID      string
	FieldLabels []string
}

type FieldRef struct {
	Label string
	Ref   string
}

func (sr *SecretResponse) GetFieldRefs() []FieldRef {
	ret := []FieldRef{}
	for _, field := range sr.FieldLabels {
		ret = append(ret, FieldRef{Label: field, Ref: fmt.Sprintf("{{op://%s/%s/%s}}", sr.VaultName, sr.ItemName, field)})
	}
	return ret
}

func (c *Client) GetSecrets(resp *ItemCreateResponse) (*SecretResponse, error) {
	ret := SecretResponse{
		AccountName: c.Target.AccountName,
		VaultName:   c.Target.VaultName,
		VaultID:     resp.Vault.ID,
		ItemName:    resp.Title,
		ItemID:      resp.ID,
	}

	for _, field := range resp.Fields {
		if field.Purpose == "" {
			ret.FieldLabels = append(ret.FieldLabels, field.Label)
		}
	}

	return &ret, nil
}

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

func (c *Client) CreateItem(envPairs map[string]string) (*SecretResponse, error) {
	req, err := c.BuildCreateItemRequest(envPairs)
	if err != nil {
		return nil, err
	}
	resp, err := c.CreateItemByRequest(req)
	if err != nil {
		return nil, err
	}
	return c.GetSecrets(resp)
}

func (c *Client) CreateItemByRequest(req ItemCreateRequest) (*ItemCreateResponse, error) {
	reqStr, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	cmd := c.BuildItemCommand("create")
	cmd.SetStdin(bytes.NewBuffer(reqStr))
	var stdout bytes.Buffer
	cmd.SetStdout(&stdout)
	cmd.SetStderr(os.Stderr)

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	var resp ItemCreateResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
