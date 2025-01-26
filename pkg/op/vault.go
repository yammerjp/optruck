package op

import (
	"bytes"
	"encoding/json"
)

type Vault struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ContentVersion int    `json:"content_version"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	Items          int    `json:"items"`
}

func (c *Client) ListVaults() ([]Vault, error) {
	cmd := c.BuildCommand(CommandOptions{
		AddAccount: true,
		AddVault:   false,
		Args:       []string{"vault", "list"},
	})
	stdoutBuffer := bytes.NewBuffer(nil)
	stderrBuffer := bytes.NewBuffer(nil)
	cmd.SetStdout(stdoutBuffer)
	cmd.SetStderr(stderrBuffer)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var resp []Vault
	if err := json.Unmarshal(stdoutBuffer.Bytes(), &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
