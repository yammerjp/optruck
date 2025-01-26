package op

import (
	"bytes"
	"encoding/json"
)

/*
$ op account list --format json
[
  {
    "url": "my.1password.com",
    "email": "mail@example.com",
    "user_uuid": "0123456789ABCDEFGH",
    "account_uuid": "ABCDEFGH0123456789"
  }
]
*/

type Account struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	UserUUID    string `json:"user_uuid"`
	AccountUUID string `json:"account_uuid"`
}

func (c *Client) ListAccounts() ([]Account, error) {
	cmd := c.BuildCommand(CommandOptions{
		AddAccount: false,
		AddVault:   false,
		Args:       []string{"account", "list"},
	})
	stdoutBuffer := bytes.NewBuffer(nil)
	stderrBuffer := bytes.NewBuffer(nil)
	cmd.SetStdout(stdoutBuffer)
	cmd.SetStderr(stderrBuffer)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var resp []Account
	if err := json.Unmarshal(stdoutBuffer.Bytes(), &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
