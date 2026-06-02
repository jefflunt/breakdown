package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadDefault(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.AgentAdapter != "gemini:google/gemini-3.1-flash-lite-preview" {
		t.Errorf("Expected default AgentAdapter gemini:google/gemini-3.1-flash-lite-preview, got %s", cfg.AgentAdapter)
	}
}

func TestConfigExpandTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot resolve home directory")
	}

	expanded := expandTilde("~/foo/bar")
	expected := filepath.Join(home, "foo/bar")

	if expanded != expected {
		t.Errorf("Expected %s, got %s", expected, expanded)
	}

	notTilde := expandTilde("/foo/bar")
	if notTilde != "/foo/bar" {
		t.Errorf("Expected /foo/bar, got %s", notTilde)
	}
}

func TestConfigLoadMissingFile(t *testing.T) {
	cfg, err := LoadConfig("/path/to/nonexistent/file.yml")
	if err != nil {
		t.Errorf("Expected no error when loading missing config file, got %v", err)
	}

	if cfg.AgentAdapter != "gemini:google/gemini-3.1-flash-lite-preview" {
		t.Errorf("Expected fallback default config, got %s", cfg.AgentAdapter)
	}
}

func TestConfigLoadValidFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configContent := `
plans_dir: "~/custom/plans"
agent_adapter: "mock:mock-provider/some-model"
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AgentAdapter != "mock:mock-provider/some-model" {
		t.Errorf("Expected AgentAdapter mock:mock-provider/some-model, got %s", cfg.AgentAdapter)
	}
}

func TestConfigLoadFromEnv(t *testing.T) {
	t.Setenv("BREAKDOWN_AGENT_ADAPTER", "custom-cli:custom-provider/custom-model")
	
	tempDir, err := os.MkdirTemp("", "config-test-env")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// An empty config file to force default fallback logic with env
	configContent := `{}`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AgentAdapter != "custom-cli:custom-provider/custom-model" {
		t.Errorf("Expected AgentAdapter to be loaded from environment, got %q", cfg.AgentAdapter)
	}
}
