package op

type Account struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	UserUUID    string `json:"user_uuid"`
	AccountUUID string `json:"account_uuid"`
}

func (c *ExecutableClient) ListAccounts() ([]Account, error) {
	cmd := c.BuildCommand("account", "list")
	var resp []Account
	if err := cmd.RunWithJSON(nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
