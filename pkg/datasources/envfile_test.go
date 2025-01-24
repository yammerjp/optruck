package datasources

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestEnvFileSource_FetchSecrets(t *testing.T) {
	// Create a temporary .env file
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	err := os.WriteFile(envPath, []byte("KEY1=value1\nKEY2=value2"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	source := &EnvFileSource{Path: envPath}
	secrets, err := source.FetchSecrets()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	if !reflect.DeepEqual(secrets, expected) {
		t.Errorf("expected %v, got %v", expected, secrets)
	}
}

func TestEnvFileSource_FetchSecrets_FileNotFound(t *testing.T) {
	source := &EnvFileSource{Path: "nonexistent.env"}
	secrets, err := source.FetchSecrets()

	if err == nil {
		t.Error("expected error, got nil")
	}
	expected := map[string]string{}
	if !reflect.DeepEqual(secrets, expected) {
		t.Errorf("expected empty map, got %v", secrets)
	}
}

func TestNewSource(t *testing.T) {
	path := "/path/to/env"
	source := NewSource(path, EnvFile)

	envSource, ok := source.(*EnvFileSource)
	if !ok {
		t.Error("expected *EnvFileSource type")
	}
	if envSource.Path != path {
		t.Errorf("expected path %s, got %s", path, envSource.Path)
	}
}
