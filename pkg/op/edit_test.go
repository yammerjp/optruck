package op

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"sort"
	"testing"

	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

var mockEditStdoutSuccess = `{
  "id": "test-id",
  "title": "test-item",
  "version": 1,
  "vault": {
    "id": "test-vault-id",
    "name": "test-vault-name"
  },
  "category": "LOGIN",
  "created_at": "2025-01-25T14:36:10.59238+09:00",
  "updated_at": "2025-01-25T14:36:10.59238+09:00",
  "additional_information": "—",
  "fields": [
    {
      "id": "FOO",
      "type": "CONCEALED",
      "label": "FOO",
      "value": "bar",
      "reference": "op://Private/test-item/FOO"
    },
    {
      "id": "BAR",
      "type": "CONCEALED",
      "label": "BAR",
      "value": "baz",
      "reference": "op://Private/test-item/BAR"
    }
  ]
}
`

var mockEditStdoutModified = `{
  "id": "test-id",
  "title": "test-item",
  "version": 1,
  "vault": {
    "id": "test-vault-id",
    "name": "test-vault-name"
  },
  "category": "LOGIN",
  "created_at": "2025-01-25T14:36:10.59238+09:00",
  "updated_at": "2025-01-25T14:36:10.59238+09:00",
  "additional_information": "—",
  "fields": [
    {
      "id": "FOO",
      "type": "CONCEALED",
      "label": "FOO",
      "value": "modified_bar",
      "reference": "op://Private/test-item/FOO"
    },
    {
      "id": "BAR",
      "type": "CONCEALED",
      "label": "BAR",
      "value": "modified_baz",
      "reference": "op://Private/test-item/BAR"
    }
  ]
}
`

var mockEditStdoutWithNewField = `{
  "id": "test-id",
  "title": "test-item",
  "version": 1,
  "vault": {
    "id": "test-vault-id",
    "name": "test-vault-name"
  },
  "category": "LOGIN",
  "created_at": "2025-01-25T14:36:10.59238+09:00",
  "updated_at": "2025-01-25T14:36:10.59238+09:00",
  "additional_information": "—",
  "fields": [
    {
      "id": "FOO",
      "type": "CONCEALED",
      "label": "FOO",
      "value": "bar",
      "reference": "op://Private/test-item/FOO"
    },
    {
      "id": "BAR",
      "type": "CONCEALED",
      "label": "BAR",
      "value": "baz",
      "reference": "op://Private/test-item/BAR"
    },
    {
      "id": "BAZ",
      "type": "CONCEALED",
      "label": "BAZ",
      "value": "qux",
      "reference": "op://Private/test-item/BAZ"
    }
  ]
}
`

func TestEditItem(t *testing.T) {
	tests := []struct {
		name           string
		itemName       string
		account        string
		vault          string
		envPairs       map[string]string
		mockStdout     string
		mockStderr     string
		mockExitStatus int
		wantErr        error
		wantRef        *SecretReference
		wantArgs       []string
		wantStdin      string
	}{
		{
			name:     "success",
			itemName: "test-item",
			account:  "test-account",
			vault:    "test-vault-name",
			envPairs: map[string]string{
				"FOO": "bar",
				"BAR": "baz",
			},
			mockStdout:     mockEditStdoutSuccess,
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
			wantArgs:  []string{"item", "edit", "test-item", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
			wantStdin: `{"fields":[{"id":"FOO","type":"CONCEALED","label":"FOO","value":"bar"},{"id":"BAR","type":"CONCEALED","label":"BAR","value":"baz"}]}`,
		},
		{
			name:     "verify exact json",
			itemName: "test-item",
			account:  "test-account",
			vault:    "test-vault-name",
			envPairs: map[string]string{
				"FOO": "bar",
			},
			mockStdout:     mockEditStdoutSuccess,
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
			wantArgs:  []string{"item", "edit", "test-item", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
			wantStdin: `{"fields":[{"id":"FOO","type":"CONCEALED","label":"FOO","value":"bar"}]}`,
		},
		{
			name:     "modify existing fields",
			itemName: "test-item",
			account:  "test-account",
			vault:    "test-vault-name",
			envPairs: map[string]string{
				"FOO": "modified_bar",
				"BAR": "modified_baz",
			},
			mockStdout:     mockEditStdoutModified,
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
			wantArgs:  []string{"item", "edit", "test-item", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
			wantStdin: `{"fields":[{"id":"FOO","type":"CONCEALED","label":"FOO","value":"modified_bar"},{"id":"BAR","type":"CONCEALED","label":"BAR","value":"modified_baz"}]}`,
		},
		{
			name:     "add new field",
			itemName: "test-item",
			account:  "test-account",
			vault:    "test-vault-name",
			envPairs: map[string]string{
				"FOO": "bar",
				"BAR": "baz",
				"BAZ": "qux",
			},
			mockStdout:     mockEditStdoutWithNewField,
			mockExitStatus: 0,
			wantErr:        nil,
			wantRef: &SecretReference{
				Account:     "test-account",
				VaultName:   "test-vault-name",
				VaultID:     "test-vault-id",
				ItemName:    "test-item",
				ItemID:      "test-id",
				FieldLabels: []string{"FOO", "BAR", "BAZ"},
			},
			wantArgs:  []string{"item", "edit", "test-item", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
			wantStdin: `{"fields":[{"id":"FOO","type":"CONCEALED","label":"FOO","value":"bar"},{"id":"BAR","type":"CONCEALED","label":"BAR","value":"baz"},{"id":"BAZ","type":"CONCEALED","label":"BAZ","value":"qux"}]}`,
		},
		{
			name:     "remove field",
			itemName: "test-item",
			account:  "test-account",
			vault:    "test-vault-name",
			envPairs: map[string]string{
				"FOO": "bar",
			},
			mockStdout:     mockEditStdoutSuccess,
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
			wantArgs:  []string{"item", "edit", "test-item", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
			wantStdin: `{"fields":[{"id":"FOO","type":"CONCEALED","label":"FOO","value":"bar"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &fakeWriter{}
			stderr := &fakeWriter{}
			var fcmd *testingexec.FakeCmd
			fcmd = &testingexec.FakeCmd{
				Stdout: stdout,
				Stderr: stderr,
				RunScript: []testingexec.FakeAction{
					func() ([]byte, []byte, error) {
						// Check stdin
						if tt.wantStdin != "" {
							buf := new(bytes.Buffer)
							if _, err := buf.ReadFrom(fcmd.Stdin); err != nil {
								t.Errorf("failed to read stdin: %v", err)
							}
							gotJSON := buf.Bytes()
							var got, want struct {
								Fields []struct {
									ID    string `json:"id"`
									Type  string `json:"type"`
									Label string `json:"label"`
									Value string `json:"value"`
								} `json:"fields"`
							}
							if err := json.Unmarshal(gotJSON, &got); err != nil {
								t.Errorf("failed to unmarshal got JSON: %v", err)
							}
							if err := json.Unmarshal([]byte(tt.wantStdin), &want); err != nil {
								t.Errorf("failed to unmarshal want JSON: %v", err)
							}

							// Sort fields by ID for comparison
							sort.Slice(got.Fields, func(i, j int) bool {
								return got.Fields[i].ID < got.Fields[j].ID
							})
							sort.Slice(want.Fields, func(i, j int) bool {
								return want.Fields[i].ID < want.Fields[j].ID
							})

							gotStr, err := json.Marshal(got)
							if err != nil {
								t.Errorf("failed to marshal got JSON: %v", err)
							}
							wantStr, err := json.Marshal(want)
							if err != nil {
								t.Errorf("failed to marshal want JSON: %v", err)
							}
							if !bytes.Equal(gotStr, wantStr) {
								t.Errorf("stdin JSON = %s, want %s", string(gotJSON), tt.wantStdin)
							}
						}

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
			fakeLogger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

			client := NewExecutableClient(fakeExec, fakeLogger).BuildItemClient(tt.account, tt.vault, tt.itemName)

			got, err := client.EditItem(tt.envPairs)
			if err != tt.wantErr {
				t.Errorf("EditItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if got.Account != tt.wantRef.Account {
					t.Errorf("EditItem() Account = %v, want %v", got.Account, tt.wantRef.Account)
				}
				if got.VaultName != tt.wantRef.VaultName {
					t.Errorf("EditItem() VaultName = %v, want %v", got.VaultName, tt.wantRef.VaultName)
				}
				if got.VaultID != tt.wantRef.VaultID {
					t.Errorf("EditItem() VaultID = %v, want %v", got.VaultID, tt.wantRef.VaultID)
				}
				if got.ItemName != tt.wantRef.ItemName {
					t.Errorf("EditItem() ItemName = %v, want %v", got.ItemName, tt.wantRef.ItemName)
				}
				if got.ItemID != tt.wantRef.ItemID {
					t.Errorf("EditItem() ItemID = %v, want %v", got.ItemID, tt.wantRef.ItemID)
				}
				if len(got.FieldLabels) != len(tt.wantRef.FieldLabels) {
					t.Errorf("EditItem() FieldLabels length = %v, want %v", len(got.FieldLabels), len(tt.wantRef.FieldLabels))
				}
				for i, label := range got.FieldLabels {
					if label != tt.wantRef.FieldLabels[i] {
						t.Errorf("EditItem() FieldLabel[%d] = %v, want %v", i, label, tt.wantRef.FieldLabels[i])
					}
				}
			}
		})
	}
}
