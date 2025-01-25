// TODO: test
package op

import (
	"errors"

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

func NewClient(target Target) *Client {
	return &Client{
		exec:   exec.New(),
		Target: target,
	}
}
func (c *Client) BuildItemCommand(args ...string) exec.Cmd {
	if c.Target.Account != "" {
		args = append(args, "--account", c.Target.Account)
	}
	if c.Target.Vault != "" {
		args = append(args, "--vault", c.Target.Vault)
	}
	args = append(args, "--format", "json", "item")
	return c.exec.Command("op", args...)
}

var ErrItemAlreadyExists = errors.New("item already exists")

func (c *Client) UploadItem(envPairs map[string]string, overwrite bool) (*SecretReference, error) {
	_, err := c.GetItem(c.Target.ItemName)
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
	// update item
	return nil, errors.New("not implemented")
}
