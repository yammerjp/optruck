package op

import (
	"log/slog"
	"os"
	"testing"

	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

var mockVaultListSuccess = `[
  {
    "id": "test-vault-id",
    "name": "test-vault",
    "content_version": 505,
    "created_at": "2022-09-04T08:17:14Z",
    "updated_at": "2025-01-25T05:36:11Z",
    "items": 300
  }
]`

func TestListVaults(t *testing.T) {
	tests := []struct {
		name           string
		account        string
		mockStdout     string
		mockStderr     string
		mockExitStatus int
		wantErr        bool
		wantVaults     []Vault
		wantArgs       []string
	}{
		{
			name:           "success",
			account:        "test-account",
			mockStdout:     mockVaultListSuccess,
			mockExitStatus: 0,
			wantErr:        false,
			wantVaults: []Vault{
				{
					ID:             "test-vault-id",
					Name:           "test-vault",
					ContentVersion: 505,
					CreatedAt:      "2022-09-04T08:17:14Z",
					UpdatedAt:      "2025-01-25T05:36:11Z",
					Items:          300,
				},
			},
			wantArgs: []string{"vault", "list", "--account", "test-account", "--format", "json"},
		},
		{
			name:           "command failure",
			account:        "test-account",
			mockStderr:     "some error",
			mockExitStatus: 1,
			wantErr:        true,
			wantArgs:       []string{"vault", "list", "--account", "test-account", "--format", "json"},
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
			fakeLogger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

			client := NewAccountClient(tt.account, fakeExec, fakeLogger)

			got, err := client.ListVaults()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListVaults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.wantVaults) {
					t.Errorf("ListVaults() got = %v vaults, want %v", len(got), len(tt.wantVaults))
					return
				}
				for i, vault := range got {
					if vault.ID != tt.wantVaults[i].ID {
						t.Errorf("ListVaults() vault[%d].ID = %v, want %v", i, vault.ID, tt.wantVaults[i].ID)
					}
					if vault.Name != tt.wantVaults[i].Name {
						t.Errorf("ListVaults() vault[%d].Name = %v, want %v", i, vault.Name, tt.wantVaults[i].Name)
					}
					if vault.ContentVersion != tt.wantVaults[i].ContentVersion {
						t.Errorf("ListVaults() vault[%d].ContentVersion = %v, want %v", i, vault.ContentVersion, tt.wantVaults[i].ContentVersion)
					}
					if vault.CreatedAt != tt.wantVaults[i].CreatedAt {
						t.Errorf("ListVaults() vault[%d].CreatedAt = %v, want %v", i, vault.CreatedAt, tt.wantVaults[i].CreatedAt)
					}
					if vault.UpdatedAt != tt.wantVaults[i].UpdatedAt {
						t.Errorf("ListVaults() vault[%d].UpdatedAt = %v, want %v", i, vault.UpdatedAt, tt.wantVaults[i].UpdatedAt)
					}
					if vault.Items != tt.wantVaults[i].Items {
						t.Errorf("ListVaults() vault[%d].Items = %v, want %v", i, vault.Items, tt.wantVaults[i].Items)
					}
				}
			}
		})
	}
}
