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

type ItemResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Version int    `json:"version"`
	Vault   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"vault"`
	Category              string `json:"category"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
	AdditionalInformation string `json:"additional_information"`
	Fields                []struct {
		ID              string `json:"id"`
		Type            string `json:"type"`
		Purpose         string `json:"purpose"`
		Label           string `json:"label"`
		Value           string `json:"value"`
		Reference       string `json:"reference"`
		PasswordDetails struct {
			Strength string `json:"strength"`
		} `json:"password_details"`
	} `json:"fields"`
}

func (sr *SecretReference) GetFieldRefs() []FieldRef {
	ret := []FieldRef{}
	for _, field := range sr.FieldLabels {
		ret = append(ret, FieldRef{Label: field, Ref: fmt.Sprintf("{{op://%s/%s/%s}}", sr.VaultID, sr.ItemID, field)})
	}
	return ret
}

func (c *AccountClient) BuildSecretReference(resp ItemResponse) *SecretReference {
	fieldLabels := []string{}
	for _, field := range resp.Fields {
		if field.Purpose == "" {
			fieldLabels = append(fieldLabels, field.Label)
		}
	}
	return &SecretReference{
		Account:     c.Account,
		VaultName:   resp.Vault.Name,
		VaultID:     resp.Vault.ID,
		ItemName:    resp.Title,
		ItemID:      resp.ID,
		FieldLabels: fieldLabels,
	}
}
