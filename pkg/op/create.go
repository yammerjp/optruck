package op

import (
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

	var resp ItemResponse
	if err := c.BuildCommand("item", "create").RunWithJson(req, &resp); err != nil {
		return nil, err
	}

	return c.BuildSecretReference(resp), nil
}
