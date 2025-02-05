package op

import (
	"errors"
	"fmt"
	"log/slog"
)

var ErrMoreThanOneItemFound = errors.New("more than one item found, please specify another item name")
var ErrItemAlreadyExists = errors.New("item already exists, use --overwrite to update")

func (c *ItemClient) UploadItem(envPairs map[string]string, overwrite bool) (*SecretReference, error) {
	refs, err := c.FilterItems(c.ItemName)
	if err != nil {
		return nil, fmt.Errorf("failed to filter items: %w. Please check the item name and try again.", err)
	}
	if len(refs) == 0 {
		slog.Debug("item not found, creating new item", "item", c.ItemName)
		return c.CreateItem(envPairs)
	}
	if len(refs) == 1 {
		slog.Debug("item found, updating existing item", "item", c.ItemName)
		if !overwrite {
			return nil, ErrItemAlreadyExists
		}
		return c.EditItem(envPairs)
	}
	return nil, ErrMoreThanOneItemFound
}
