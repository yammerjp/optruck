package op

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

var ErrMoreThanOneItemMatches = errors.New("more than one item matches")
var ErrItemNotFound = errors.New("item not found")

func (c *Client) GetItem() (*SecretReference, error) {
	cmd := c.BuildCommand(CommandOptions{
		AddAccount: true,
		AddVault:   true,
		Args:       []string{"item", "get", c.Target.ItemName},
	})
	stdoutBuffer := bytes.NewBuffer(nil)
	stderrBuffer := bytes.NewBuffer(nil)
	cmd.SetStdout(stdoutBuffer)
	cmd.SetStderr(stderrBuffer)

	var resp ItemResponse
	if err := cmd.Run(); err != nil {
		slog.Error("failed to get item", "error", err)
		stderrStr := stderrBuffer.String()
		slog.Error("stderr", "stderr", stderrStr)
		if strings.Contains(stderrStr, " isn't an item") && strings.Contains(stderrStr, " Specify the item with its UUID, name, or domain.") {
			slog.Error("item not found")
			return nil, ErrItemNotFound
		}
		if strings.Contains(stderrBuffer.String(), " More than one item matches ") {
			slog.Error("more than one item matches")
			return nil, ErrMoreThanOneItemMatches
		}
		slog.Error("failed to get item", "error", err)
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if err := json.Unmarshal(stdoutBuffer.Bytes(), &resp); err != nil {
		return nil, err
	}

	return c.BuildSecretReference(resp), nil
}
