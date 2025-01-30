package op

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
)

type ItemEditRequest struct {
	Fields []ItemEditRequestField `json:"fields"`
}

type ItemEditRequestField struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
}

func (c *ItemClient) EditItem(envPairs map[string]string) (*SecretReference, error) {
	req := ItemEditRequest{
		Fields: make([]ItemEditRequestField, 0, len(envPairs)),
	}

	// Sort keys to ensure consistent order
	keys := make([]string, 0, len(envPairs))
	for k := range envPairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Add fields in sorted order
	for _, k := range keys {
		req.Fields = append(req.Fields, ItemEditRequestField{
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

	var stdout bytes.Buffer
	cmd := c.BuildCommand("item", "edit", c.ItemName)
	cmd.SetStdin(bytes.NewBuffer(reqStr))
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
