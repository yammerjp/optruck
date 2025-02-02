package op

import (
	utilExec "github.com/yammerjp/optruck/internal/util/exec"

	"testing"

	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

var mockListStdoutSuccess = `[
  {
    "id": "test-id-1",
    "title": "test-item-1",
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
        "id": "FOO",
        "type": "CONCEALED",
        "label": "FOO",
        "value": "bar",
        "reference": "op://test-account/test-vault-name/test-item-1/FOO"
      },
      {
        "id": "BAR",
        "type": "CONCEALED",
        "label": "BAR",
        "value": "baz",
        "reference": "op://test-account/test-vault-name/test-item-1/BAR"
      }
    ]
  },
  {
    "id": "test-id-2",
    "title": "test-item-2",
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
        "id": "BAZ",
        "type": "CONCEALED",
        "label": "BAZ",
        "value": "qux",
        "reference": "op://test-account/test-vault-name/test-item-2/BAZ"
      }
    ]
  }
]`

func TestListItems(t *testing.T) {
	tests := []struct {
		name           string
		account        string
		vault          string
		mockStdout     string
		mockStderr     string
		mockExitStatus int
		wantErr        error
		wantRefs       []SecretReference
		wantArgs       []string
	}{
		{
			name:           "success",
			account:        "test-account",
			vault:          "test-vault-name",
			mockStdout:     mockListStdoutSuccess,
			mockExitStatus: 0,
			wantErr:        nil,
			wantRefs: []SecretReference{
				{
					Account:     "test-account",
					VaultName:   "test-vault-name",
					VaultID:     "test-vault-id",
					ItemName:    "test-item-1",
					ItemID:      "test-id-1",
					FieldLabels: []string{"FOO", "BAR"},
				},
				{
					Account:     "test-account",
					VaultName:   "test-vault-name",
					VaultID:     "test-vault-id",
					ItemName:    "test-item-2",
					ItemID:      "test-id-2",
					FieldLabels: []string{"BAZ"},
				},
			},
			wantArgs: []string{"item", "list", "--account", "test-account", "--vault", "test-vault-name", "--format", "json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcmd := &testingexec.FakeCmd{
				RunScript: []testingexec.FakeAction{
					func() ([]byte, []byte, error) {
						if tt.mockExitStatus != 0 {
							return []byte(tt.mockStdout), []byte(tt.mockStderr), &testingexec.FakeExitError{Status: tt.mockExitStatus}
						}
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

			client := NewVaultClient(tt.account, tt.vault)
			utilExec.SetExec(fakeExec)

			got, err := client.ListItems()
			if err != tt.wantErr {
				t.Errorf("ListItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr == nil {
				if len(got) != len(tt.wantRefs) {
					t.Errorf("ListItems() got %d items, want %d", len(got), len(tt.wantRefs))
					return
				}

				for i, ref := range got {
					if ref.Account != tt.wantRefs[i].Account {
						t.Errorf("ListItems()[%d] Account = %v, want %v", i, ref.Account, tt.wantRefs[i].Account)
					}
					if ref.VaultName != tt.wantRefs[i].VaultName {
						t.Errorf("ListItems()[%d] VaultName = %v, want %v", i, ref.VaultName, tt.wantRefs[i].VaultName)
					}
					if ref.VaultID != tt.wantRefs[i].VaultID {
						t.Errorf("ListItems()[%d] VaultID = %v, want %v", i, ref.VaultID, tt.wantRefs[i].VaultID)
					}
					if ref.ItemName != tt.wantRefs[i].ItemName {
						t.Errorf("ListItems()[%d] ItemName = %v, want %v", i, ref.ItemName, tt.wantRefs[i].ItemName)
					}
					if ref.ItemID != tt.wantRefs[i].ItemID {
						t.Errorf("ListItems()[%d] ItemID = %v, want %v", i, ref.ItemID, tt.wantRefs[i].ItemID)
					}
					if len(ref.FieldLabels) != len(tt.wantRefs[i].FieldLabels) {
						t.Errorf("ListItems()[%d] FieldLabels length = %v, want %v", i, len(ref.FieldLabels), len(tt.wantRefs[i].FieldLabels))
					}
					for j, label := range ref.FieldLabels {
						if label != tt.wantRefs[i].FieldLabels[j] {
							t.Errorf("ListItems()[%d] FieldLabel[%d] = %v, want %v", i, j, label, tt.wantRefs[i].FieldLabels[j])
						}
					}
				}
			}
		})
	}
}
