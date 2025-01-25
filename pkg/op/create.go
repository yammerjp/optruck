package op

import (
	"bytes"
	"encoding/json"
	"os"
)

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

type ItemCreateResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Version int    `json:"version"`
	Vault   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"vault"`
	Category              string `json:"category"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
	AdditionalInformation string `json:"additional_information"`
	Fields                []struct {
		ID              string `json:"id"`
		Type            string `json:"type"`
		Purpose         string `json:"purpose"`
		Label           string `json:"label"`
		Value           string `json:"value"`
		Reference       string `json:"reference"`
		PasswordDetails struct {
			Strength string `json:"strength"`
		} `json:"password_details"`
	} `json:"fields"`
}

func (c *Client) CreateItem(envPairs map[string]string) (*SecretReference, error) {
	req := ItemCreateRequest{
		Title:    c.Target.ItemName,
		Category: "LOGIN",
	}

	for k, v := range envPairs {
		req.Fields = append(req.Fields, ItemCreateRequestField{
			ID:    k,
			Type:  "CONCEALED",
			Label: k,
			Value: v,
		})
	}

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
