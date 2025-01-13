package op

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

type ItemRequest struct {
	Vault string
	ID    string
}

type ItemResponse struct {
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

func (req *ItemRequest) GetItem() (*ItemResponse, error) {
	cmd := exec.Command("op", "item", "get", req.ID, "--vault", req.Vault, "--format", "json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	// TODO: stderrの中の文字列もチェック
	// exit codeが1かつstderrに"More than one item matches"が含まれているときは、itemが複数存在する
	// exit codeが1かつstderrに"Item not found"が含まれているときは、itemが存在しない

	// $ op item get example.com --format json
	// [ERROR] 2025/01/14 08:22:56 "example.com" isn't an item. Specify the item with its UUID, name, or domain.
	// $ op item get xxxxxxxxxxxxx.jp --format json
	// [ERROR] 2025/01/14 08:23:10 More than one item matches "xxxxxxxxxxxxxxx.jp". Try again and specify the item by its ID:
	//    * for the item "xxxxxxxxxx xxxxxxxxxxxxxxxxxxxxxxxxxxx" in vault vault-name: xxxxxxxxxxxxxxxxxxxxxxxxxx
	//    * for the item "xxxxxxxxxxxxxxxx xxxx" in vault vault-name: xxxxxxxxxxxxxxxxxxxxxxxxxx

	//

	var item ItemResponse
	if err := json.Unmarshal(stdout.Bytes(), &item); err != nil {
		return nil, err
	}
	return &item, nil
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
