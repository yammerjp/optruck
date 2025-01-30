package op

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
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

func (c *ItemClient) CreateItem(envPairs map[string]string) (*SecretReference, error) {
	req := ItemCreateRequest{
		Title:    c.ItemName,
		Category: "LOGIN",
		Fields:   make([]ItemCreateRequestField, 0, len(envPairs)),
	}

	// Sort keys to ensure consistent order
	keys := make([]string, 0, len(envPairs))
	for k := range envPairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Add fields in sorted order
	for _, k := range keys {
		req.Fields = append(req.Fields, ItemCreateRequestField{
			ID:    k,
			Type:  "CONCEALED",
			Label: k,
			Value: envPairs[k],
		})
	}

	reqStr, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	cmd := c.BuildCommand("item", "create")
	cmd.SetStdin(bytes.NewBuffer(reqStr))
	var stdout bytes.Buffer
	cmd.SetStdout(&stdout)
	cmd.SetStderr(os.Stderr)

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	var resp ItemResponse
	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, err
	}

	return c.BuildSecretReference(resp), nil
}
