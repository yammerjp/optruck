package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yammerjp/optruck/pkg/op"
)

func TestEnvTemplateDestWrite(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "optruck-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name            string
		dest            *EnvTemplateDest
		secretReference *op.SecretReference
		expected        string
	}{
		{
			name: "basic case",
			dest: &EnvTemplateDest{
				Path: filepath.Join(tmpDir, "test1.env"),
			},
			secretReference: &op.SecretReference{
				Account:     "test.1password.com",
				VaultName:   "TestVault",
				VaultID:     "vault-id",
				ItemID:      "item-id",
				ItemName:    "TestItem",
				FieldLabels: []string{"DB_USER", "DB_PASS"},
			},
			expected: `# This file was generated by optruck.
#   - 1password account: test.1password.com
#   - 1password vault: TestVault
# To restore, run the following command:
#   $ op inject -i test1.env --account test.1password.com -o .env
DB_USER={{op://vault-id/item-id/DB_USER}}
DB_PASS={{op://vault-id/item-id/DB_PASS}}
`,
		},
		{
			name: "without account name",
			dest: &EnvTemplateDest{
				Path: filepath.Join(tmpDir, "test2.env"),
			},
			secretReference: &op.SecretReference{
				VaultName:   "TestVault",
				VaultID:     "vault-id",
				ItemName:    "TestItem",
				ItemID:      "item-id",
				FieldLabels: []string{"API_KEY"},
			},
			expected: `# This file was generated by optruck.
#   - 1password vault: TestVault
# To restore, run the following command:
#   $ op inject -i test2.env -o .env
API_KEY={{op://vault-id/item-id/API_KEY}}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write the template
			err := tc.dest.Write(tc.secretReference)
			if err != nil {
				t.Fatalf("Write() error = %v", err)
			}

			// Read the generated file
			content, err := os.ReadFile(tc.dest.Path)
			if err != nil {
				t.Fatalf("failed to read generated file: %v", err)
			}

			// Compare the content
			if got := string(content); got != tc.expected {
				t.Errorf("Write() generated content mismatch\nwant:\n%s\ngot:\n%s", tc.expected, got)
			}
		})
	}
}
