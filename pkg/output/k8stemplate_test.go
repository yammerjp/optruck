package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yammerjp/optruck/pkg/op"
)

func TestK8sSecretTemplateDestWrite(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "optruck-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name         string
		dest         *K8sSecretTemplateDest
		resp         *op.SecretReference
		expected     string
		overwrite    bool
		existingFile bool // file exists
	}{
		{
			name: "basic case",
			dest: &K8sSecretTemplateDest{
				Path:       filepath.Join(tmpDir, "test1.yaml"),
				Namespace:  "default",
				SecretName: "TestItem",
			},
			resp: &op.SecretReference{
				VaultName:   "TestVault",
				VaultID:     "vault-id",
				Account:     "test.1password.com",
				ItemName:    "TestItem",
				ItemID:      "item-id",
				FieldLabels: []string{"DB_USER", "DB_PASS"},
			},
			expected: `# This file was generated by optruck.
#   - 1password account: test.1password.com
#   - 1password vault: TestVault
# To restore, run the following command:
#   $ op inject -i test1.yaml | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: TestItem
  namespace: default
type: Opaque
data:
  DB_USER: {{op://vault-id/item-id/DB_USER}}
  DB_PASS: {{op://vault-id/item-id/DB_PASS}}
`,
			overwrite:    false,
			existingFile: false,
		},
		{
			name: "without account name",
			dest: &K8sSecretTemplateDest{
				Path:       filepath.Join(tmpDir, "test2.yaml"),
				Namespace:  "production",
				SecretName: "APISecret",
			},
			resp: &op.SecretReference{
				VaultName:   "TestVault",
				VaultID:     "vault-id",
				ItemName:    "APISecret",
				ItemID:      "item-id",
				FieldLabels: []string{"API_KEY"},
			},
			expected: `# This file was generated by optruck.
#   - 1password vault: TestVault
# To restore, run the following command:
#   $ op inject -i test2.yaml | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: APISecret
  namespace: production
type: Opaque
data:
  API_KEY: {{op://vault-id/item-id/API_KEY}}
`,
			overwrite:    false,
			existingFile: false,
		},
		{
			name: "overwrite existing file",
			dest: &K8sSecretTemplateDest{
				Path:       filepath.Join(tmpDir, "test3.yaml"),
				Namespace:  "default",
				SecretName: "TestItem",
			},
			resp: &op.SecretReference{
				VaultName:   "TestVault",
				VaultID:     "vault-id",
				Account:     "test.1password.com",
				ItemName:    "TestItem",
				ItemID:      "item-id",
				FieldLabels: []string{"API_KEY"},
			},
			expected: `# This file was generated by optruck.
#   - 1password account: test.1password.com
#   - 1password vault: TestVault
# To restore, run the following command:
#   $ op inject -i test3.yaml | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: TestItem
  namespace: default
type: Opaque
data:
  API_KEY: {{op://vault-id/item-id/API_KEY}}
`,
			overwrite:    true,
			existingFile: true,
		},
		{
			name: "error when file exists and overwrite is false",
			dest: &K8sSecretTemplateDest{
				Path:       filepath.Join(tmpDir, "test4.yaml"),
				Namespace:  "default",
				SecretName: "TestItem",
			},
			resp: &op.SecretReference{
				VaultName:   "TestVault",
				VaultID:     "vault-id",
				Account:     "test.1password.com",
				ItemName:    "TestItem",
				ItemID:      "item-id",
				FieldLabels: []string{"API_KEY"},
			},
			expected:     "", // Content doesn't matter as we expect an error
			overwrite:    false,
			existingFile: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.existingFile {
				// First write with some initial content
				initialContent := []byte("initial content")
				err := os.WriteFile(tc.dest.Path, initialContent, 0644)
				if err != nil {
					t.Fatalf("failed to write initial content: %v", err)
				}
			}

			// Write the template
			err := tc.dest.Write(tc.resp, tc.overwrite)

			if tc.existingFile && !tc.overwrite {
				if err == nil {
					t.Error("Write() should return error when file exists and overwrite is false")
				}
				return
			}

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
