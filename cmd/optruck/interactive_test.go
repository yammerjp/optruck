package optruck

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/manifoldco/promptui"
	"k8s.io/utils/exec"
)

type MockInteractiveRunner struct {
	selectResponses []struct {
		index int
		value string
		err   error
	}
	inputResponses []struct {
		value string
		err   error
	}
	selectIndex int
	inputIndex  int
}

func (m *MockInteractiveRunner) Select(_ promptui.Select) (int, string, error) {
	if m.selectIndex >= len(m.selectResponses) {
		return 0, "", nil
	}
	resp := m.selectResponses[m.selectIndex]
	m.selectIndex++
	return resp.index, resp.value, resp.err
}

func (m *MockInteractiveRunner) Input(_ promptui.Prompt) (string, error) {
	if m.inputIndex >= len(m.inputResponses) {
		return "", nil
	}
	resp := m.inputResponses[m.inputIndex]
	m.inputIndex++
	return resp.value, resp.err
}

type MockExec struct {
	exec.Interface
	commands map[string][]string
	outputs  map[string]string
}

func NewMockExec() *MockExec {
	return &MockExec{
		commands: make(map[string][]string),
		outputs:  make(map[string]string),
	}
}

type MockCmd struct {
	cmd    string
	args   []string
	output string
	err    error
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
	env    []string
}

func (m *MockExec) Command(cmd string, args ...string) exec.Cmd {
	key := cmd + " " + strings.Join(args, " ")
	m.commands[key] = args
	var err error
	if m.outputs[key] == "" {
		err = fmt.Errorf("command failed: %s %s", cmd, strings.Join(args, " "))
	}
	return &MockCmd{
		cmd:    cmd,
		args:   args,
		output: m.outputs[key],
		err:    err,
	}
}

func (c *MockCmd) SetStdout(stdout io.Writer) {
	c.stdout = stdout
}

func (c *MockCmd) SetStderr(stderr io.Writer) {
	c.stderr = stderr
}

func (c *MockCmd) SetStdin(stdin io.Reader) {
	c.stdin = stdin
}

func (c *MockCmd) Run() error {
	if c.err != nil {
		return c.err
	}
	if c.stdout != nil {
		if w, ok := c.stdout.(*bytes.Buffer); ok {
			w.WriteString(c.output)
		}
	}
	return nil
}

func (c *MockCmd) CombinedOutput() ([]byte, error) {
	if c.err != nil {
		return nil, c.err
	}
	return []byte(c.output), nil
}

func (c *MockCmd) Output() ([]byte, error) {
	if c.err != nil {
		return nil, c.err
	}
	return []byte(c.output), nil
}

func (c *MockCmd) SetEnv(env []string) {
	c.env = env
}

func (c *MockCmd) SetDir(dir string) {
	// No-op for mock
}

func (c *MockCmd) StdinPipe() (io.WriteCloser, error) {
	return nil, nil
}

func (c *MockCmd) StdoutPipe() (io.ReadCloser, error) {
	return nil, nil
}

func (c *MockCmd) StderrPipe() (io.ReadCloser, error) {
	return nil, nil
}

func (c *MockCmd) Start() error {
	return nil
}

func (c *MockCmd) Wait() error {
	return nil
}

func (c *MockCmd) Stop() {
	// No-op for mock
}

func TestSetDataSourceInteractively(t *testing.T) {
	tests := []struct {
		name     string
		cli      *CLI
		mock     *MockInteractiveRunner
		mockExec *MockExec
		wantErr  bool
		wantFile string
		wantK8s  string
	}{
		{
			name: "select env file",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{
				selectResponses: []struct {
					index int
					value string
					err   error
				}{
					{0, "env file", nil},
				},
				inputResponses: []struct {
					value string
					err   error
				}{
					{".env", nil},
				},
			},
			mockExec: NewMockExec(),
			wantErr:  false,
			wantFile: ".env",
			wantK8s:  "",
		},
		{
			name: "select k8s secret",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{
				selectResponses: []struct {
					index int
					value string
					err   error
				}{
					{1, "k8s secret", nil},
					{0, "default", nil},
					{0, "mysecret", nil},
				},
			},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["kubectl get namespaces -o jsonpath={.items[*].metadata.name}"] = "default kube-system"
				m.outputs["kubectl get secrets -n default --field-selector type=Opaque -o jsonpath={.items[*].metadata.name}"] = "mysecret secret1"
				return m
			}(),
			wantErr:  false,
			wantFile: "",
			wantK8s:  "mysecret",
		},
		{
			name:     "data source already set with env file",
			cli:      &CLI{EnvFile: "existing.env"},
			mock:     &MockInteractiveRunner{},
			mockExec: NewMockExec(),
			wantErr:  false,
			wantFile: "existing.env",
			wantK8s:  "",
		},
		{
			name:     "data source already set with k8s secret",
			cli:      &CLI{K8sSecret: "existing-secret"},
			mock:     &MockInteractiveRunner{},
			mockExec: NewMockExec(),
			wantErr:  false,
			wantFile: "",
			wantK8s:  "existing-secret",
		},
		{
			name: "kubectl get namespaces fails",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{
				selectResponses: []struct {
					index int
					value string
					err   error
				}{
					{1, "k8s secret", nil},
				},
			},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["kubectl get namespaces -o jsonpath={.items[*].metadata.name}"] = ""
				return m
			}(),
			wantErr:  true,
			wantFile: "",
			wantK8s:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.runner = tt.mock
			tt.cli.exec = tt.mockExec
			err := tt.cli.setDataSourceInteractively()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.cli.EnvFile != tt.wantFile {
				t.Errorf("EnvFile = %v, want %v", tt.cli.EnvFile, tt.wantFile)
			}
			if tt.cli.K8sSecret != tt.wantK8s {
				t.Errorf("K8sSecret = %v, want %v", tt.cli.K8sSecret, tt.wantK8s)
			}
		})
	}
}

func TestSetTargetAccountInteractively(t *testing.T) {
	tests := []struct {
		name      string
		cli       *CLI
		mock      *MockInteractiveRunner
		mockExec  *MockExec
		wantErr   bool
		wantValue string
	}{
		{
			name: "select account",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{
				selectResponses: []struct {
					index int
					value string
					err   error
				}{
					{0, "my.1password.com", nil},
				},
			},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op account list --format json"] = `[{"url": "my.1password.com", "email": "test@example.com"}]`
				return m
			}(),
			wantErr:   false,
			wantValue: "my.1password.com",
		},
		{
			name:      "account already set",
			cli:       &CLI{Account: "existing.1password.com"},
			mock:      &MockInteractiveRunner{},
			mockExec:  NewMockExec(),
			wantErr:   false,
			wantValue: "existing.1password.com",
		},
		{
			name: "op account list fails",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op account list --format json"] = ""
				return m
			}(),
			wantErr:   true,
			wantValue: "",
		},
		{
			name: "no accounts available",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op account list --format json"] = "[]"
				return m
			}(),
			wantErr:   true,
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.runner = tt.mock
			tt.cli.exec = tt.mockExec
			err := tt.cli.setTargetAccountInteractively()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.cli.Account != tt.wantValue {
				t.Errorf("Account = %v, want %v", tt.cli.Account, tt.wantValue)
			}
		})
	}
}

func TestSetTargetVaultInteractively(t *testing.T) {
	tests := []struct {
		name      string
		cli       *CLI
		mock      *MockInteractiveRunner
		mockExec  *MockExec
		wantErr   bool
		wantValue string
	}{
		{
			name: "select vault",
			cli:  &CLI{Account: "my.1password.com"},
			mock: &MockInteractiveRunner{
				selectResponses: []struct {
					index int
					value string
					err   error
				}{
					{0, "vault1", nil},
				},
			},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op vault list --account my.1password.com --format json"] = `[{"id": "vault1", "name": "Vault 1"}]`
				return m
			}(),
			wantErr:   false,
			wantValue: "vault1",
		},
		{
			name:      "vault already set",
			cli:       &CLI{Account: "my.1password.com", Vault: "existing-vault"},
			mock:      &MockInteractiveRunner{},
			mockExec:  NewMockExec(),
			wantErr:   false,
			wantValue: "existing-vault",
		},
		{
			name:      "account not set",
			cli:       &CLI{},
			mock:      &MockInteractiveRunner{},
			mockExec:  NewMockExec(),
			wantErr:   true,
			wantValue: "",
		},
		{
			name: "op vault list fails",
			cli:  &CLI{Account: "my.1password.com"},
			mock: &MockInteractiveRunner{},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op vault list --account my.1password.com --format json"] = ""
				return m
			}(),
			wantErr:   true,
			wantValue: "",
		},
		{
			name: "no vaults available",
			cli:  &CLI{Account: "my.1password.com"},
			mock: &MockInteractiveRunner{},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op vault list --account my.1password.com --format json"] = "[]"
				return m
			}(),
			wantErr:   true,
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.runner = tt.mock
			tt.cli.exec = tt.mockExec
			err := tt.cli.setTargetVaultInteractively()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.cli.Vault != tt.wantValue {
				t.Errorf("Vault = %v, want %v", tt.cli.Vault, tt.wantValue)
			}
		})
	}
}

func TestSetTargetItemInteractively(t *testing.T) {
	tests := []struct {
		name      string
		cli       *CLI
		mock      *MockInteractiveRunner
		mockExec  *MockExec
		wantErr   bool
		wantValue string
	}{
		{
			name: "create new item",
			cli: &CLI{
				Account: "my.1password.com",
				Vault:   "vault1",
			},
			mock: &MockInteractiveRunner{
				selectResponses: []struct {
					index int
					value string
					err   error
				}{
					{1, "create new", nil},
				},
				inputResponses: []struct {
					value string
					err   error
				}{
					{"new-item", nil},
				},
			},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op item list --account my.1password.com --vault vault1 --format json"] = `[{"id": "item1", "title": "Item 1"}]`
				return m
			}(),
			wantErr:   false,
			wantValue: "new-item",
		},
		{
			name: "select existing item",
			cli: &CLI{
				Account:   "my.1password.com",
				Vault:     "vault1",
				Overwrite: true,
			},
			mock: &MockInteractiveRunner{
				selectResponses: []struct {
					index int
					value string
					err   error
				}{
					{0, "item1", nil},
				},
			},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op item list --account my.1password.com --vault vault1 --format json"] = `[{"id": "item1", "title": "Item 1"}]`
				return m
			}(),
			wantErr:   false,
			wantValue: "item1",
		},
		{
			name: "item already set",
			cli: &CLI{
				Account: "my.1password.com",
				Vault:   "vault1",
				Item:    "existing-item",
			},
			mock:      &MockInteractiveRunner{},
			mockExec:  NewMockExec(),
			wantErr:   false,
			wantValue: "existing-item",
		},
		{
			name:      "account not set",
			cli:       &CLI{},
			mock:      &MockInteractiveRunner{},
			mockExec:  NewMockExec(),
			wantErr:   true,
			wantValue: "",
		},
		{
			name:      "vault not set",
			cli:       &CLI{Account: "my.1password.com"},
			mock:      &MockInteractiveRunner{},
			mockExec:  NewMockExec(),
			wantErr:   true,
			wantValue: "",
		},
		{
			name: "op item list fails",
			cli: &CLI{
				Account: "my.1password.com",
				Vault:   "vault1",
			},
			mock: &MockInteractiveRunner{},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op item list --account my.1password.com --vault vault1 --format json"] = ""
				return m
			}(),
			wantErr:   true,
			wantValue: "",
		},
		{
			name: "no items available with overwrite",
			cli: &CLI{
				Account:   "my.1password.com",
				Vault:     "vault1",
				Overwrite: true,
			},
			mock: &MockInteractiveRunner{},
			mockExec: func() *MockExec {
				m := NewMockExec()
				m.outputs["op item list --account my.1password.com --vault vault1 --format json"] = "[]"
				return m
			}(),
			wantErr:   true,
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.runner = tt.mock
			tt.cli.exec = tt.mockExec
			err := tt.cli.setTargetItemInteractively()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.cli.Item != tt.wantValue {
				t.Errorf("Item = %v, want %v", tt.cli.Item, tt.wantValue)
			}
		})
	}
}

func TestSetDestInteractively(t *testing.T) {
	tests := []struct {
		name      string
		cli       *CLI
		mock      *MockInteractiveRunner
		wantErr   bool
		wantValue string
	}{
		{
			name: "set output path",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{
				inputResponses: []struct {
					value string
					err   error
				}{
					{"output.env", nil},
				},
			},
			wantErr:   false,
			wantValue: "output.env",
		},
		{
			name: "set output path with confirmation",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{
				inputResponses: []struct {
					value string
					err   error
				}{
					{"existing.env", nil},
				},
				selectResponses: []struct {
					index int
					value string
					err   error
				}{
					{0, "overwrite", nil},
				},
			},
			wantErr:   false,
			wantValue: "existing.env",
		},
		{
			name:      "output already set",
			cli:       &CLI{Output: "existing.env"},
			mock:      &MockInteractiveRunner{},
			wantErr:   false,
			wantValue: "existing.env",
		},
		{
			name: "invalid output path",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{
				inputResponses: []struct {
					value string
					err   error
				}{
					{"/invalid/path/output.env", fmt.Errorf("directory /invalid/path does not exist")},
				},
			},
			wantErr:   true,
			wantValue: "",
		},
		{
			name: "empty output path",
			cli:  &CLI{},
			mock: &MockInteractiveRunner{
				inputResponses: []struct {
					value string
					err   error
				}{
					{"", fmt.Errorf("output path is required")},
				},
			},
			wantErr:   true,
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.runner = tt.mock
			err := tt.cli.setDestInteractively()
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.cli.Output != tt.wantValue {
				t.Errorf("Output = %v, want %v", tt.cli.Output, tt.wantValue)
			}
		})
	}
}
