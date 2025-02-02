package op

func (c *VaultClient) ListItems() ([]SecretReference, error) {
	cmd := c.BuildCommand("item", "list")
	var resp []ItemResponse
	if err := cmd.RunWithJSON(nil, &resp); err != nil {
		return nil, err
	}

	refs := make([]SecretReference, len(resp))
	for i, item := range resp {
		refs[i] = *c.BuildSecretReference(item)
	}
	return refs, nil
}

func (c *VaultClient) FilterItems(itemSpecifier string) ([]SecretReference, error) {
	refs, err := c.ListItems()
	if err != nil {
		return nil, err
	}

	filtered := make([]SecretReference, 0)
	for _, ref := range refs {
		if itemSpecifier == ref.ItemName || itemSpecifier == ref.ItemID {
			filtered = append(filtered, ref)
		}
	}
	return filtered, nil
}
