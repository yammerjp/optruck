package op

import (
	"bytes"
	"testing"

	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

var mockCreateStdoutSuccess = `
{
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
      "id": "username",
      "type": "STRING",
      "purpose": "USERNAME",
      "label": "username",
      "reference": "op://Private/test-item/username"
    },
    {
      "id": "password",
      "type": "CONCEALED",
      "purpose": "PASSWORD",
      "label": "password",
      "reference": "op://Private/test-item/password",
      "password_details": {}
    },
    {
      "id": "notesPlain",
      "type": "STRING",
      "purpose": "NOTES",
      "label": "notesPlain",
      "reference": "op://Private/test-item/notesPlain"
    },
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

func TestCreateItem(t *testing.T) {
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
			mockStdout:     mockCreateStdoutSuccess,
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
			wantArgs:  []string{"item", "create", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
			wantStdin: `{"Title":"test-item","Category":"LOGIN","Fields":[{"ID":"FOO","Type":"CONCEALED","Purpose":"","Label":"FOO","Value":"bar"},{"ID":"BAR","Type":"CONCEALED","Purpose":"","Label":"BAR","Value":"baz"}]}`,
		},
		{
			name:     "without account and vault",
			itemName: "test-item",
			envPairs: map[string]string{
				"FOO": "bar",
				"BAR": "baz",
			},
			mockStdout:     mockCreateStdoutSuccess,
			mockExitStatus: 0,
			wantErr:        nil,
			wantRef: &SecretReference{
				VaultName:   "test-vault-name",
				VaultID:     "test-vault-id",
				ItemName:    "test-item",
				ItemID:      "test-id",
				FieldLabels: []string{"FOO", "BAR"},
			},
			wantArgs:  []string{"item", "create", "--format", "json"},
			wantStdin: `{"Title":"test-item","Category":"LOGIN","Fields":[{"ID":"FOO","Type":"CONCEALED","Purpose":"","Label":"FOO","Value":"bar"},{"ID":"BAR","Type":"CONCEALED","Purpose":"","Label":"BAR","Value":"baz"}]}`,
		},
		{
			name:     "verify exact json",
			itemName: "test-item",
			account:  "test-account",
			vault:    "test-vault-name",
			envPairs: map[string]string{
				"FOO": "bar",
			},
			mockStdout:     mockCreateStdoutSuccess,
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
			wantArgs:  []string{"item", "create", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
			wantStdin: `{"Title":"test-item","Category":"LOGIN","Fields":[{"ID":"FOO","Type":"CONCEALED","Purpose":"","Label":"FOO","Value":"bar"}]}`,
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
							gotJSON := fcmd.Stdin.(*bytes.Buffer).Bytes()
							if !bytes.Equal(gotJSON, []byte(tt.wantStdin)) {
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

			client := &Client{
				exec: fakeExec,
				Target: Target{
					Account:  tt.account,
					Vault:    tt.vault,
					ItemName: tt.itemName,
				},
			}

			got, err := client.CreateItem(tt.envPairs)
			if err != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if got.Account != tt.wantRef.Account {
					t.Errorf("CreateItem() Account = %v, want %v", got.Account, tt.wantRef.Account)
				}
				if got.VaultName != tt.wantRef.VaultName {
					t.Errorf("CreateItem() VaultName = %v, want %v", got.VaultName, tt.wantRef.VaultName)
				}
				if got.VaultID != tt.wantRef.VaultID {
					t.Errorf("CreateItem() VaultID = %v, want %v", got.VaultID, tt.wantRef.VaultID)
				}
				if got.ItemName != tt.wantRef.ItemName {
					t.Errorf("CreateItem() ItemName = %v, want %v", got.ItemName, tt.wantRef.ItemName)
				}
				if got.ItemID != tt.wantRef.ItemID {
					t.Errorf("CreateItem() ItemID = %v, want %v", got.ItemID, tt.wantRef.ItemID)
				}
				if len(got.FieldLabels) != len(tt.wantRef.FieldLabels) {
					t.Errorf("CreateItem() FieldLabels length = %v, want %v", len(got.FieldLabels), len(tt.wantRef.FieldLabels))
				}
				for i, label := range got.FieldLabels {
					if label != tt.wantRef.FieldLabels[i] {
						t.Errorf("CreateItem() FieldLabel[%d] = %v, want %v", i, label, tt.wantRef.FieldLabels[i])
					}
				}
			}
		})
	}
}
