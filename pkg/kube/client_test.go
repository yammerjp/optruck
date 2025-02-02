package kube

import (
	"errors"
	"reflect"
	"testing"

	utilExec "github.com/yammerjp/optruck/internal/util/exec"

	"k8s.io/utils/exec"
	testingexec "k8s.io/utils/exec/testing"
)

func TestGetSecret(t *testing.T) {
	tests := []struct {
		name         string
		namespace    string
		secretName   string
		mockStdout   string
		mockStderr   string
		exitStatus   int
		cmdError     error
		expectedErr  bool
		expectedData map[string]string
	}{
		{
			name:       "success",
			namespace:  "default",
			secretName: "mysecret",
			mockStdout: `{"key1":"dmFsdWUx","key2":"dmFsdWUy"}`,
			exitStatus: 0,
			expectedData: map[string]string{
				"key1": "dmFsdWUx",
				"key2": "dmFsdWUy",
			},
		},
		{
			name:        "invalid json",
			namespace:   "default",
			secretName:  "mysecret",
			mockStdout:  `invalid json`,
			exitStatus:  0,
			expectedErr: true,
		},
		{
			name:        "secret not found",
			namespace:   "default",
			secretName:  "nonexistent",
			mockStderr:  "Error from server (NotFound): secrets \"nonexistent\" not found",
			exitStatus:  1,
			expectedErr: true,
		},
		{
			name:        "kubectl not found",
			namespace:   "default",
			secretName:  "mysecret",
			cmdError:    errors.New("executable file not found in $PATH"),
			expectedErr: true,
		},
		{
			name:        "empty namespace",
			namespace:   "",
			secretName:  "mysecret",
			mockStderr:  "error: required flag(s) \"namespace\" not set",
			exitStatus:  1,
			expectedErr: true,
		},
		{
			name:        "empty secret name",
			namespace:   "default",
			secretName:  "",
			mockStderr:  "error: resource name may not be empty",
			exitStatus:  1,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcmd := &testingexec.FakeCmd{
				RunScript: []testingexec.FakeAction{
					func() ([]byte, []byte, error) {
						if tt.cmdError != nil {
							return nil, nil, tt.cmdError
						}
						if tt.exitStatus != 0 {
							return []byte(tt.mockStdout), []byte(tt.mockStderr), &testingexec.FakeExitError{Status: tt.exitStatus}
						}
						return []byte(tt.mockStdout), []byte(tt.mockStderr), nil
					},
				},
			}

			fakeExec := &testingexec.FakeExec{
				CommandScript: []testingexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						expectedCmd := "kubectl"
						expectedArgs := []string{"get", "secret", "-n", tt.namespace, tt.secretName, "-o", "jsonpath={.data}"}
						if cmd != expectedCmd {
							t.Errorf("expected command %q but got %q", expectedCmd, cmd)
						}
						if !reflect.DeepEqual(args, expectedArgs) {
							t.Errorf("expected args %v but got %v", expectedArgs, args)
						}
						return fcmd
					},
				},
			}

			client := NewClient()
			utilExec.SetExec(fakeExec)
			data, err := client.GetSecret(tt.namespace, tt.secretName)

			if tt.expectedErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectedErr && !reflect.DeepEqual(data, tt.expectedData) {
				t.Errorf("expected %v, got %v", tt.expectedData, data)
			}
		})
	}
}

func TestGetNamespaces(t *testing.T) {
	tests := []struct {
		name          string
		mockStdout    string
		mockStderr    string
		exitStatus    int
		cmdError      error
		expectedErr   bool
		expectedNames []string
	}{
		{
			name:          "success",
			mockStdout:    "default kube-system kube-public",
			exitStatus:    0,
			expectedNames: []string{"default", "kube-system", "kube-public"},
		},
		{
			name:          "empty",
			mockStdout:    "",
			exitStatus:    0,
			expectedNames: []string{},
		},
		{
			name:        "error",
			mockStderr:  "error getting namespaces",
			exitStatus:  1,
			expectedErr: true,
		},
		{
			name:        "kubectl not found",
			cmdError:    errors.New("executable file not found in $PATH"),
			expectedErr: true,
		},
		{
			name:        "permission denied",
			mockStderr:  "error: forbidden: User \"system:anonymous\" cannot list resource \"namespaces\"",
			exitStatus:  1,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcmd := &testingexec.FakeCmd{
				RunScript: []testingexec.FakeAction{
					func() ([]byte, []byte, error) {
						if tt.cmdError != nil {
							return nil, nil, tt.cmdError
						}
						if tt.exitStatus != 0 {
							return []byte(tt.mockStdout), []byte(tt.mockStderr), &testingexec.FakeExitError{Status: tt.exitStatus}
						}
						return []byte(tt.mockStdout), []byte(tt.mockStderr), nil
					},
				},
			}

			fakeExec := &testingexec.FakeExec{
				CommandScript: []testingexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						expectedCmd := "kubectl"
						expectedArgs := []string{"get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}"}
						if cmd != expectedCmd {
							t.Errorf("expected command %q but got %q", expectedCmd, cmd)
						}
						if !reflect.DeepEqual(args, expectedArgs) {
							t.Errorf("expected args %v but got %v", expectedArgs, args)
						}
						return fcmd
					},
				},
			}

			client := NewClient()
			utilExec.SetExec(fakeExec)
			names, err := client.GetNamespaces()

			if tt.expectedErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectedErr && !reflect.DeepEqual(names, tt.expectedNames) {
				t.Errorf("expected %v, got %v", tt.expectedNames, names)
			}
		})
	}
}

func TestGetSecrets(t *testing.T) {
	tests := []struct {
		name          string
		namespace     string
		mockStdout    string
		mockStderr    string
		exitStatus    int
		cmdError      error
		expectedErr   bool
		expectedNames []string
	}{
		{
			name:          "success",
			namespace:     "default",
			mockStdout:    "secret1 secret2 secret3",
			exitStatus:    0,
			expectedNames: []string{"secret1", "secret2", "secret3"},
		},
		{
			name:          "empty",
			namespace:     "default",
			mockStdout:    "",
			exitStatus:    0,
			expectedNames: []string{},
		},
		{
			name:        "error",
			namespace:   "default",
			mockStderr:  "error getting secrets",
			exitStatus:  1,
			expectedErr: true,
		},
		{
			name:        "kubectl not found",
			namespace:   "default",
			cmdError:    errors.New("executable file not found in $PATH"),
			expectedErr: true,
		},
		{
			name:        "empty namespace",
			namespace:   "",
			mockStderr:  "error: required flag(s) \"namespace\" not set",
			exitStatus:  1,
			expectedErr: true,
		},
		{
			name:        "permission denied",
			namespace:   "default",
			mockStderr:  "error: forbidden: User \"system:anonymous\" cannot list resource \"secrets\" in API group \"\" in the namespace \"default\"",
			exitStatus:  1,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcmd := &testingexec.FakeCmd{
				RunScript: []testingexec.FakeAction{
					func() ([]byte, []byte, error) {
						if tt.cmdError != nil {
							return nil, nil, tt.cmdError
						}
						if tt.exitStatus != 0 {
							return []byte(tt.mockStdout), []byte(tt.mockStderr), &testingexec.FakeExitError{Status: tt.exitStatus}
						}
						return []byte(tt.mockStdout), []byte(tt.mockStderr), nil
					},
				},
			}

			fakeExec := &testingexec.FakeExec{
				CommandScript: []testingexec.FakeCommandAction{
					func(cmd string, args ...string) exec.Cmd {
						expectedCmd := "kubectl"
						expectedArgs := []string{"get", "secrets", "-n", tt.namespace, "--field-selector", "type=Opaque", "-o", "jsonpath={.items[*].metadata.name}"}
						if cmd != expectedCmd {
							t.Errorf("expected command %q but got %q", expectedCmd, cmd)
						}
						if !reflect.DeepEqual(args, expectedArgs) {
							t.Errorf("expected args %v but got %v", expectedArgs, args)
						}
						return fcmd
					},
				},
			}

			client := NewClient()
			utilExec.SetExec(fakeExec)
			names, err := client.GetSecrets(tt.namespace)

			if tt.expectedErr && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectedErr && !reflect.DeepEqual(names, tt.expectedNames) {
				t.Errorf("expected %v, got %v", tt.expectedNames, names)
			}
		})
	}
}
