package op

import "fmt"

type SecretReference struct {
	Account     string
	VaultName   string
	VaultID     string
	ItemName    string
	ItemID      string
	FieldLabels []string
}

type FieldRef struct {
	Label string
	Ref   string
}

func (sr *SecretReference) GetFieldRefs() []FieldRef {
	ret := []FieldRef{}
	for _, field := range sr.FieldLabels {
		ret = append(ret, FieldRef{Label: field, Ref: fmt.Sprintf("{{op://%s/%s/%s}}", sr.VaultID, sr.ItemID, field)})
	}
	return ret
}

func buildSecretReferenceByItemCreateResponse(resp *ItemCreateResponse, account string) *SecretReference {
	fieldLabels := []string{}
	for _, field := range resp.Fields {
		if field.Purpose == "" {
			fieldLabels = append(fieldLabels, field.Label)
		}
	}
	return &SecretReference{
		Account:     account,
		VaultName:   resp.Vault.Name,
		VaultID:     resp.Vault.ID,
		ItemName:    resp.Title,
		ItemID:      resp.ID,
		FieldLabels: fieldLabels,
	}
}

func buildSecretReferenceByItemGetResponse(resp *GetItemResponse, account string) *SecretReference {
	fieldLabels := []string{}
	for _, field := range resp.Fields {
		if field.Purpose == "" {
			fieldLabels = append(fieldLabels, field.Label)
		}
	}
	return &SecretReference{
		Account:     account,
		VaultName:   resp.Vault.Name,
		VaultID:     resp.Vault.ID,
		ItemName:    resp.Title,
		ItemID:      resp.ID,
		FieldLabels: fieldLabels,
	}
}
