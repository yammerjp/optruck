package op

import (
	"fmt"
	"sort"

	"github.com/yammerjp/optruck/internal/errors"
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
	if len(envPairs) == 0 {
		return nil, errors.NewInvalidArgumentError(
			"環境変数",
			"環境変数が空です",
			"環境変数ファイルまたはKubernetesシークレットに値が設定されていることを確認してください",
		)
	}

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

	cmd := c.BuildCommand("item", "create")
	var resp ItemResponse
	if err := cmd.RunWithJSON(req, &resp); err != nil {
		return nil, errors.NewOperationFailedError(
			fmt.Sprintf("1Passwordアイテム '%s' の作成", c.ItemName),
			err,
			"1Password CLIが正しく設定されているか、アカウントとボールトにアクセス権限があるか確認してください",
		)
	}
	return c.BuildSecretReference(resp), nil
}
