package op

import (
	"context"
	"encoding/json"
	"fmt"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

/*
{
  "title": "",
  "category": "PASSWORD",
  "fields": [
    {
      "id": "password",
      "type": "CONCEALED",
      "purpose": "PASSWORD",
      "label": "password",
      "password_details": {
        "strength": "TERRIBLE"
      },
      "value": ""
    },
    {
      "id": "notesPlain",
      "type": "STRING",
      "purpose": "NOTES",
      "label": "notesPlain",
      "value": ""
    }
  ]
}

*/

type ItemCreateRequest struct {
	Title    string
	Category string
	Fields   []ItemCreateRequestField
}

type ItemCreateRequestField struct {
	ID      string
	Type    string
	Purpose string
	Label   string
	Value   string
}

func (c *Client) BuildCreateItemRequest(ctx context.Context, vaultName, itemName string, envPairs map[string]string) (ItemCreateRequest, error) {
	ret := ItemCreateRequest{
		Title:    itemName,
		Category: "PASSWORD",
	}

	for k, v := range envPairs {
		ret.Fields = append(ret.Fields, ItemCreateRequestField{
			ID:      k,
			Type:    "CONCEALED",
			Purpose: "PASSWORD",
			Label:   k,
			Value:   v,
		})
	}

	return ret, nil
}

func (c *Client) CreateItem(ctx context.Context, vaultName, itemName string, envPairs map[string]string) error {
	req, err := c.BuildCreateItemRequest(ctx, vaultName, itemName, envPairs)
	if err != nil {
		return err
	}
	return c.CreateItemByRequest(ctx, req)
}

func (c *Client) CreateItemByRequest(ctx context.Context, req ItemCreateRequest) error {
	json, err := json.Marshal(req)
	if err != nil {
		return err
	}

	fmt.Println(string(json))

	return nil
}
