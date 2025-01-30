package op

import (
	"testing"

	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

type fakeWriter struct {
	data []byte
}

func (w *fakeWriter) Write(p []byte) (n int, err error) {
	w.data = append(w.data, p...)
	return len(p), nil
}

var mockStdoutSuccess = `{
  "id": "test-id",
  "title": "test-item",
  "version": 1,
  "vault": {
    "id": "test-vault-id",
    "name": "test-vault-name"
  },
  "category": "LOGIN",
  "created_at": "2025-01-21T00:11:02.842101+09:00",
  "updated_at": "2025-01-21T00:11:02.842102+09:00",
  "additional_information": "—",
  "fields": [
    {
      "id": "username",
      "type": "STRING",
      "purpose": "USERNAME",
      "label": "username",
      "reference": "op://test-account/test-vault-name/test-item/username"
    },
    {
      "id": "password",
      "type": "CONCEALED",
      "purpose": "PASSWORD",
      "label": "password",
      "reference": "op://test-account/test-vault-name/test-item/password",
      "password_details": {}
    },
    {
      "id": "notesPlain",
      "type": "STRING",
      "purpose": "NOTES",
      "label": "notesPlain",
      "reference": "op://test-account/test-vault-name/test-item/notesPlain"
    },
    {
      "id": "FOO",
      "type": "CONCEALED",
      "label": "FOO",
      "value": "bar",
      "reference": "op://test-account/test-vault-name/test-item/FOO"
    },
    {
      "id": "BAR",
      "type": "CONCEALED",
      "label": "BAR",
      "value": "baz",
      "reference": "op://optruck-development/first-item/BAR"
    }
  ]
}`

func TestGetItem(t *testing.T) {
	tests := []struct {
		name           string
		itemName       string
		account        string
		vault          string
		mockStdout     string
		mockStderr     string
		mockExitStatus int
		wantErr        error
		wantRef        *SecretReference
		wantArgs       []string
	}{
		{
			name:           "success",
			itemName:       "test-item",
			account:        "test-account",
			vault:          "test-vault-name",
			mockStdout:     mockStdoutSuccess,
			mockExitStatus: 0,
			wantErr:        nil,
			wantRef: &SecretReference{
				Account:     "test-account",
				VaultName:   "test-vault-name",
				VaultID:     "test-vault-id",
				ItemName:    "test-item",
				ItemID:      "test-id",
				FieldLabels: []string{"FOO", "BAR"},
			},
			wantArgs: []string{"item", "get", "test-item", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
		},
		{
			name:           "item not found",
			itemName:       "non-existent",
			account:        "test-account",
			vault:          "test-vault-name",
			mockStderr:     `[ERROR] 2025/01/25 14:04:51 "unknown" isn't an item. Specify the item with its UUID, name, or domain.`,
			mockExitStatus: 1,
			wantErr:        ErrItemNotFound,
			wantArgs:       []string{"item", "get", "non-existent", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
		},
		{
			name:     "multiple items found",
			itemName: "ambiguous",
			account:  "test-account",
			vault:    "test-vault-name",
			mockStderr: `[ERROR] 2025/01/25 14:03:20 More than one item matches "ambiguous". Try again and specify the item by its ID: 
* for the item "ambiguous" in vault test-vault-name: xxxxxxxxxxxxxxxxxxxxxxxxxx
* for the item "ambiguous" in vault test-vault-name: xxxxxxxxxxxxxxxxxxxxxxxxxx`,
			mockExitStatus: 1,
			wantErr:        ErrMoreThanOneItemMatches,
			wantArgs:       []string{"item", "get", "ambiguous", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &fakeWriter{}
			stderr := &fakeWriter{}
			fcmd := &testingexec.FakeCmd{
				Stdout: stdout,
				Stderr: stderr,
				RunScript: []testingexec.FakeAction{
					func() ([]byte, []byte, error) {
						if tt.mockExitStatus != 0 {
							stderr.Write([]byte(tt.mockStderr))
							return []byte(tt.mockStdout), []byte(tt.mockStderr), &testingexec.FakeExitError{Status: tt.mockExitStatus}
						}
						stdout.Write([]byte(tt.mockStdout))
						return []byte(tt.mockStdout), []byte(tt.mockStderr), nil
					},
				},
			}

			fakeExec := &testingexec.FakeExec{
				CommandScript: []testingexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						if cmd != "op" {
							t.Errorf("expected command 'op', got %s", cmd)
						}
						if len(args) != len(tt.wantArgs) {
							t.Errorf("expected %d arguments, got %d: %v", len(tt.wantArgs), len(args), args)
						} else {
							for i, arg := range args {
								if arg != tt.wantArgs[i] {
									t.Errorf("argument %d: expected %s, got %s", i, tt.wantArgs[i], arg)
								}
							}
						}
						return fcmd
					},
				},
			}

			client := NewItemClient(tt.account, tt.vault, tt.itemName, fakeExec)

			got, err := client.GetItem()
			if err != tt.wantErr {
				t.Errorf("GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if got.Account != tt.wantRef.Account {
					t.Errorf("GetItem() Account = %v, want %v", got.Account, tt.wantRef.Account)
				}
				if got.VaultName != tt.wantRef.VaultName {
					t.Errorf("GetItem() VaultName = %v, want %v", got.VaultName, tt.wantRef.VaultName)
				}
				if got.VaultID != tt.wantRef.VaultID {
					t.Errorf("GetItem() VaultID = %v, want %v", got.VaultID, tt.wantRef.VaultID)
				}
				if got.ItemName != tt.wantRef.ItemName {
					t.Errorf("GetItem() ItemName = %v, want %v", got.ItemName, tt.wantRef.ItemName)
				}
				if got.ItemID != tt.wantRef.ItemID {
					t.Errorf("GetItem() ItemID = %v, want %v", got.ItemID, tt.wantRef.ItemID)
				}
				if len(got.FieldLabels) != len(tt.wantRef.FieldLabels) {
					t.Errorf("GetItem() FieldLabels length = %v, want %v", len(got.FieldLabels), len(tt.wantRef.FieldLabels))
				}
				for i, label := range got.FieldLabels {
					if label != tt.wantRef.FieldLabels[i] {
						t.Errorf("GetItem() FieldLabel[%d] = %v, want %v", i, label, tt.wantRef.FieldLabels[i])
					}
				}
			}
		})
	}
}
