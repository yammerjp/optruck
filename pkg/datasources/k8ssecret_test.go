package datasources

import (
	utilExec "github.com/yammerjp/optruck/internal/util/exec"

	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/yammerjp/optruck/pkg/kube"
	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

func TestValidateDNS1123Subdomain(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid name",
			input:   "my-secret-123",
			wantErr: false,
		},
		{
			name:    "valid name with dots",
			input:   "my.secret.123",
			wantErr: false,
		},
		{
			name:    "invalid: uppercase letters",
			input:   "My-Secret",
			wantErr: true,
		},
		{
			name:    "invalid: underscore",
			input:   "my_secret",
			wantErr: true,
		},
		{
			name:    "invalid: special characters",
			input:   "my@secret!123",
			wantErr: true,
		},
		{
			name:    "invalid: starts with hyphen",
			input:   "-mysecret",
			wantErr: true,
		},
		{
			name:    "invalid: ends with hyphen",
			input:   "mysecret-",
			wantErr: true,
		},
		{
			name:    "invalid: starts with dot",
			input:   ".mysecret",
			wantErr: true,
		},
		{
			name:    "invalid: ends with dot",
			input:   "mysecret.",
			wantErr: true,
		},
		{
			name:    "invalid: empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid: too long",
			input:   strings.Repeat("a", maxDNS1123SubdomainLength+1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDNS1123Subdomain(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDNS1123Subdomain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestK8sSecretSource_FetchSecrets(t *testing.T) {
	tests := []struct {
		name            string
		namespace       string
		secretName      string
		mockOutput      string
		mockErr         error
		wantErr         bool
		want            map[string]string
		expectedCommand string
		expectedArgs    []string
		skipCommandTest bool // flag to skip command test
	}{
		{
			name:       "success",
			namespace:  "default",
			secretName: "mysecret",
			mockOutput: `{"key1":"dmFsdWUx","key2":"dmFsdWUy"}`,
			mockErr:    nil,
			wantErr:    false,
			want: map[string]string{
				"key1": "dmFsdWUx",
				"key2": "dmFsdWUy",
			},
			expectedCommand: "kubectl",
			expectedArgs:    []string{"get", "secret", "-n", "default", "mysecret", "-o", "jsonpath={.data}"},
		},
		{
			name:            "invalid json",
			namespace:       "test-ns",
			secretName:      "test-secret",
			mockOutput:      `invalid json`,
			mockErr:         nil,
			wantErr:         true,
			want:            nil,
			expectedCommand: "kubectl",
			expectedArgs:    []string{"get", "secret", "-n", "test-ns", "test-secret", "-o", "jsonpath={.data}"},
		},
		{
			name:            "command execution error",
			namespace:       "default",
			secretName:      "mysecret",
			mockOutput:      "",
			mockErr:         errors.New("command failed: secret not found"),
			wantErr:         true,
			want:            nil,
			expectedCommand: "kubectl",
			expectedArgs:    []string{"get", "secret", "-n", "default", "mysecret", "-o", "jsonpath={.data}"},
		},
		{
			name:            "empty secret data",
			namespace:       "default",
			secretName:      "empty-secret",
			mockOutput:      `{}`,
			mockErr:         nil,
			wantErr:         false,
			want:            map[string]string{},
			expectedCommand: "kubectl",
			expectedArgs:    []string{"get", "secret", "-n", "default", "empty-secret", "-o", "jsonpath={.data}"},
		},
		{
			name:            "invalid namespace name",
			namespace:       "Invalid_Namespace@123",
			secretName:      "mysecret",
			mockOutput:      "",
			mockErr:         nil,
			wantErr:         true,
			want:            nil,
			skipCommandTest: true, // バリデーションエラーでコマンドは実行されない
		},
		{
			name:            "invalid secret name",
			namespace:       "default",
			secretName:      "Invalid_Secret@123",
			mockOutput:      "",
			mockErr:         nil,
			wantErr:         true,
			want:            nil,
			skipCommandTest: true, // バリデーションエラーでコマンドは実行されない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExec := &testingexec.FakeExec{}
			var capturedCommand string
			var capturedArgs []string

			fakeCmd := &testingexec.FakeCmd{
				RunScript: []testingexec.FakeAction{
					func() ([]byte, []byte, error) {
						return []byte(tt.mockOutput), nil, tt.mockErr
					},
				},
			}

			cmdAction := func(cmd string, args ...string) exec.Cmd {
				capturedCommand = cmd
				capturedArgs = args
				return fakeCmd
			}
			fakeExec.CommandScript = []testingexec.FakeCommandAction{cmdAction}

			client := &kube.Client{}
			utilExec.SetExec(fakeExec)

			source := &K8sSecretSource{
				Namespace:  tt.namespace,
				SecretName: tt.secretName,
				Client:     client,
			}

			got, err := source.FetchSecrets()

			// コマンドと引数の検証（バリデーションエラーの場合はスキップ）
			if !tt.skipCommandTest {
				if capturedCommand != tt.expectedCommand {
					t.Errorf("expected command %q, got %q", tt.expectedCommand, capturedCommand)
				}
				if !reflect.DeepEqual(capturedArgs, tt.expectedArgs) {
					t.Errorf("expected args %v, got %v", tt.expectedArgs, capturedArgs)
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("K8sSecretSource.FetchSecrets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("K8sSecretSource.FetchSecrets() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
