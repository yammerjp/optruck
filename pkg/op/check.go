package op

func (c *Client) CheckItemExists(accountName, vaultName, itemName string) (bool, error) {
	item, err := c.GetItem(accountName, vaultName, itemName)
	if err != nil {
		return false, err
	}

	return item != nil, nil
}
