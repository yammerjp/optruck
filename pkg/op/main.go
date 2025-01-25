// TODO: test
package op

import (
	"errors"
	"fmt"

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

var ErrItemAlreadyExists = errors.New("item already exists")

func (c *Client) UploadItem(envPairs map[string]string, overwrite bool) (*SecretReference, error) {
	_, err := c.GetItem()
	if err != nil {
		if err == ErrItemNotFound {
			return c.CreateItem(envPairs)
		}
		if err == ErrMoreThanOneItemMatches {
			return nil, ErrMoreThanOneItemMatches
		}
		return nil, err
	}

	if !overwrite {
		return nil, ErrItemAlreadyExists
	}

	return c.EditItem(envPairs)
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
