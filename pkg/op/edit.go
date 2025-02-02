package op

import (
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

	cmd := c.BuildCommand("item", "edit", c.ItemName)
	var resp ItemResponse
	if err := cmd.RunWithJSON(req, &resp); err != nil {
		return nil, err
	}

	return c.BuildSecretReference(resp), nil
}
