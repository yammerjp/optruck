package op

type Vault struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ContentVersion int    `json:"content_version"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	Items          int    `json:"items"`
}

func (c *AccountClient) ListVaults() ([]Vault, error) {
	cmd := c.BuildCommand("vault", "list")
	var resp []Vault
	if err := cmd.RunWithJSON(nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
