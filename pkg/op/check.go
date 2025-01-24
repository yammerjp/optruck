package op

func (c *Client) CheckItemExists(itemName string) (bool, error) {
	item, err := c.GetItem(itemName)
	if err != nil {
		return false, err
	}

	return item != nil, nil
}
