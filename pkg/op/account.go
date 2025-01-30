package op

import (
	"bytes"
	"encoding/json"
)

type Account struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	UserUUID    string `json:"user_uuid"`
	AccountUUID string `json:"account_uuid"`
}

func (c *ExecutableClient) ListAccounts() ([]Account, error) {
	cmd := c.BuildCommand("account", "list")
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
