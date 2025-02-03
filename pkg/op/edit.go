package op

import (
	"fmt"
	"sort"

	"github.com/yammerjp/optruck/internal/errors"
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
	if len(envPairs) == 0 {
		return nil, errors.NewInvalidArgumentError(
			"環境変数",
			"環境変数が空です",
			"環境変数ファイルまたはKubernetesシークレットに値が設定されていることを確認してください",
		)
	}

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
		return nil, errors.NewOperationFailedError(
			fmt.Sprintf("1Passwordアイテム '%s' の編集", c.ItemName),
			err,
			"1Password CLIが正しく設定されているか、アカウントとボールトにアクセス権限があるか確認してください",
		)
	}

	return c.BuildSecretReference(resp), nil
}
