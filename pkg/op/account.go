package op

type Account struct {
	URL         string `json:"url"`
	Email       string `json:"email"`
	UserUUID    string `json:"user_uuid"`
	AccountUUID string `json:"account_uuid"`
}

func (c *ExecutableClient) ListAccounts() ([]Account, error) {
	var resp []Account
	if err := c.BuildCommand("account", "list").RunWithJson(nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
