package op

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var ErrMoreThanOneItemMatches = errors.New("more than one item matches")
var ErrItemNotFound = errors.New("item not found")

func (c *Client) GetItem() (*SecretReference, error) {
	cmd := c.BuildItemCommand("get", c.Target.ItemName)
	stdoutBuffer := bytes.NewBuffer(nil)
	stderrBuffer := bytes.NewBuffer(nil)
	cmd.SetStdout(stdoutBuffer)
	cmd.SetStderr(stderrBuffer)

	var resp ItemResponse
	if err := cmd.Run(); err != nil {
		if strings.Contains(stderrBuffer.String(), " isn't an item. Specify the item with its UUID, name, or domain.") {
			return nil, ErrItemNotFound
		}
		if strings.Contains(stderrBuffer.String(), " More than one item matches ") {
			return nil, ErrMoreThanOneItemMatches
		}
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if err := json.Unmarshal(stdoutBuffer.Bytes(), &resp); err != nil {
		return nil, err
	}

	return c.BuildSecretReference(resp), nil
}
