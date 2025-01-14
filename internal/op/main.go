package op

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

type ItemGetRequest struct {
	Vault    string
	IDorName string
}

type Item struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Vault struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"vault"`
	Tags    []string `json:"tags"`
	Version int      `json:"version"`
	Fields  []struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		Purpose   string `json:"purpose"`
		Label     string `json:"label"`
		Password  string `json:"password"`
		Reference string `json:"reference"`
	} `json:"fields"`
}

type ItemGetResponse Item

type ItemCreateRequest Item

func (req *ItemGetRequest) GetItem() (*ItemGetResponse, error) {
	cmd := exec.Command("op", "item", "get", req.IDorName, "--vault", req.Vault, "--format", "json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	stderrStr := stderr.String()

	if err != nil && strings.Contains(stderrStr, "More than one item matches") {
		return nil, nil
	} else if err != nil && strings.Contains(stderrStr, "Item not found") {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else if stderr.Len() > 0 {
		return nil, errors.New("Unknown errror: stderr contains " + stderr.String())
	}

	var item ItemGetResponse
	if err := json.Unmarshal(stdout.Bytes(), &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (req *ItemCreateRequest) CreateItem() error {
	cmd := exec.Command("op", "item", "create")
	// TODO: stdinにjsonを渡す
}

/*
$ op item get xxxxxx --format json | cat
{
  "id": "xxxxxxxxxxxxxxxxxxxxxxxxxx",
  "title": "example.com",
  "tags": ["2000-01-01 に 1PIF でインポート"],
  "version": 1,
  "vault": {
    "id": "xxxxxxxxxxxxxxxxxxxxxxxxxx",
    "name": "vault-name"
  },
  "category": "LOGIN",
  "last_edited_by": "xxxxxxxxxxxxxxxxxxxxxxxxxx",
  "created_at": "2000-01-01T00:00:00Z",
  "updated_at": "2000-01-01T00:00:00Z",
  "additional_information": "xxxxxx",
  "urls": [
    {
      "label": "website",
      "primary": true,
      "href": "xxx"
    }
  ],
  "fields": [
    {
      "id": "username",
      "type": "STRING",
      "purpose": "USERNAME",
      "label": "username",
      "value": "xxxxxx",
      "reference": "op://vault-name/example.com/username"
    },
    {
      "id": "password",
      "type": "CONCEALED",
      "purpose": "PASSWORD",
      "label": "password",
      "value": "xxxxxxxxxxxxxx",
      "reference": "op://vault-name/example.com/password",
      "password_details": {
        "strength": "VERY_GOOD"
      }
    },
    {
      "id": "notesPlain",
      "type": "STRING",
      "purpose": "NOTES",
      "label": "notesPlain",
      "reference": "op://vault-name/example.com/notesPlain"
    }
  ]
}

func (i *Item) CheckIfItemExists() error {
	cmd := exec.Command("op", "item", "get", i.ID, "--vault", i.Vault)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (i *Item) CreateItem() error {

}
*/
