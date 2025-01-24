package op

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var ErrMoreThanOneItemMatches = errors.New("more than one item matches")

type GetItemResponse struct {
	ID    string `json:"id"`
	Vault struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"vault"`
	Fields []struct {
		ID        string `json:"id"`
		Label     string `json:"label"`
		Purpose   string `json:"purpose"`
		Reference string `json:"reference"`
	} `json:"fields"`
}

func (c *Client) GetItem(itemName string) (*GetItemResponse, error) {
	cmd := c.BuildItemCommand("get", itemName)
	stdoutBuffer := bytes.NewBuffer(nil)
	stderrBuffer := bytes.NewBuffer(nil)
	cmd.SetStdout(stdoutBuffer)
	cmd.SetStderr(stderrBuffer)

	var resp GetItemResponse
	if err := cmd.Run(); err != nil {
		if strings.Contains(stderrBuffer.String(), " isn't an item. Specify the item with its UUID, name, or domain.") {
			// not found
			return nil, nil
		}
		if strings.Contains(stderrBuffer.String(), " More than one item matches ") {
			return nil, ErrMoreThanOneItemMatches
		}
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if err := json.Unmarshal(stdoutBuffer.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
