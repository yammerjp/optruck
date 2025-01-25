// TODO: test
package op

import (
	"errors"
	"fmt"
	"log/slog"

	"k8s.io/utils/exec"
)

type Target struct {
	Account  string
	Vault    string
	ItemName string
}

type Client struct {
	exec exec.Interface
	Target
}

func (target Target) BuildClient() *Client {
	return &Client{
		exec:   exec.New(),
		Target: target,
	}
}

var ErrMoreThanOneItemFound = errors.New("more than one item found, please specify another item name")
var ErrItemAlreadyExists = errors.New("item already exists, use --overwrite to update")

func (c *Client) BuildItemCommand(args ...string) exec.Cmd {
	args = append([]string{"item"}, args...)

	if c.Target.Account != "" {
		args = append(args, "--account", c.Target.Account)
	}
	if c.Target.Vault != "" {
		args = append(args, "--vault", c.Target.Vault)
	}
	args = append(args, "--format", "json")
	return c.exec.Command("op", args...)
}

func (c *Client) UploadItem(envPairs map[string]string, overwrite bool) (*SecretReference, error) {
	refs, err := c.FilterItems(c.Target.ItemName)
	if err != nil {
		return nil, fmt.Errorf("failed to filter items: %w", err)
	}
	if len(refs) == 0 {
		slog.Debug("item not found, creating new item", "item", c.Target.ItemName)
		return c.CreateItem(envPairs)
	}
	if len(refs) == 1 {
		slog.Debug("item found, updating existing item", "item", c.Target.ItemName)
		if !overwrite {
			return nil, ErrItemAlreadyExists
		}
		return c.EditItem(envPairs)
	}
	return nil, ErrMoreThanOneItemFound
}

type SecretReference struct {
	Account     string
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

func (sr *SecretReference) GetFieldRefs() []FieldRef {
	ret := []FieldRef{}
	for _, field := range sr.FieldLabels {
		ret = append(ret, FieldRef{Label: field, Ref: fmt.Sprintf("{{op://%s/%s/%s}}", sr.VaultID, sr.ItemID, field)})
	}
	return ret
}
