package op

import (
	"testing"

	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

var mockAccountListSuccess = `[
  {
    "url": "my.1password.com",
    "email": "mail@example.com",
    "user_uuid": "0123456789ABCDEFGH",
    "account_uuid": "ABCDEFGH0123456789"
  }
]`

func TestListAccounts(t *testing.T) {
	tests := []struct {
		name           string
		mockStdout     string
		mockStderr     string
		mockExitStatus int
		wantErr        bool
		wantAccounts   []Account
		wantArgs       []string
	}{
		{
			name:           "success",
			mockStdout:     mockAccountListSuccess,
			mockExitStatus: 0,
			wantErr:        false,
			wantAccounts: []Account{
				{
					URL:         "my.1password.com",
					Email:       "mail@example.com",
					UserUUID:    "0123456789ABCDEFGH",
					AccountUUID: "ABCDEFGH0123456789",
				},
			},
			wantArgs: []string{"account", "list", "--format", "json"},
		},
		{
			name:           "command failure",
			mockStderr:     "some error",
			mockExitStatus: 1,
			wantErr:        true,
			wantArgs:       []string{"account", "list", "--format", "json"},
		},
		{
			name:           "invalid json",
			mockStdout:     `{"invalid": "json"`,
			mockExitStatus: 0,
			wantErr:        true,
			wantArgs:       []string{"account", "list", "--format", "json"},
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

			client := NewExecutableClient(fakeExec)

			got, err := client.ListAccounts()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.wantAccounts) {
					t.Errorf("ListAccounts() got = %v accounts, want %v", len(got), len(tt.wantAccounts))
					return
				}
				for i, account := range got {
					if account.URL != tt.wantAccounts[i].URL {
						t.Errorf("ListAccounts() account[%d].URL = %v, want %v", i, account.URL, tt.wantAccounts[i].URL)
					}
					if account.Email != tt.wantAccounts[i].Email {
						t.Errorf("ListAccounts() account[%d].Email = %v, want %v", i, account.Email, tt.wantAccounts[i].Email)
					}
					if account.UserUUID != tt.wantAccounts[i].UserUUID {
						t.Errorf("ListAccounts() account[%d].UserUUID = %v, want %v", i, account.UserUUID, tt.wantAccounts[i].UserUUID)
					}
					if account.AccountUUID != tt.wantAccounts[i].AccountUUID {
						t.Errorf("ListAccounts() account[%d].AccountUUID = %v, want %v", i, account.AccountUUID, tt.wantAccounts[i].AccountUUID)
					}
				}
			}
		})
	}
}
