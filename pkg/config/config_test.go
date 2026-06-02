package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoadDefault(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.LLM.Provider != "gemini" {
		t.Errorf("Expected default provider gemini, got %s", cfg.LLM.Provider)
	}
	if cfg.LLM.Model != "gemini-3.1-flash-lite-preview" {
		t.Errorf("Expected default model gemini-3.1-flash-lite-preview, got %s", cfg.LLM.Model)
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

	if cfg.LLM.Provider != "gemini" {
		t.Errorf("Expected fallback default config")
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

	if cfg.LLM.Provider != "mock" {
		t.Errorf("Expected provider mock, got %s", cfg.LLM.Provider)
	}
}

func TestParseAdapter(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedCli      string
		expectedProvider string
		expectedModel    string
		expectErr        bool
	}{
		{
			name:             "valid format",
			input:            "copilot:anthropic/claude-haiku-4.5",
			expectedCli:      "copilot",
			expectedProvider: "anthropic",
			expectedModel:    "claude-haiku-4.5",
			expectErr:        false,
		},
		{
			name:             "valid format with spaces",
			input:            "  copilot  :  anthropic / claude-haiku-4.5  ",
			expectedCli:      "copilot",
			expectedProvider: "anthropic",
			expectedModel:    "claude-haiku-4.5",
			expectErr:        false,
		},
		{
			name:             "valid format with extra slashes in model",
			input:            "copilot:anthropic/claude/haiku/4.5",
			expectedCli:      "copilot",
			expectedProvider: "anthropic",
			expectedModel:    "claude/haiku/4.5",
			expectErr:        false,
		},
		{
			name:      "empty input",
			input:     "",
			expectErr: true,
		},
		{
			name:      "missing colon",
			input:     "copilot-anthropic/claude",
			expectErr: true,
		},
		{
			name:      "missing slash",
			input:     "copilot:anthropic-claude",
			expectErr: true,
		},
		{
			name:      "empty cli name",
			input:     ":anthropic/claude",
			expectErr: true,
		},
		{
			name:      "empty provider name",
			input:     "copilot:/claude",
			expectErr: true,
		},
		{
			name:      "empty model name",
			input:     "copilot:anthropic/",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cli, provider, model, err := ParseAdapter(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tc.input, err)
				}
				if cli != tc.expectedCli {
					t.Errorf("expected cli %q, got %q", tc.expectedCli, cli)
				}
				if provider != tc.expectedProvider {
					t.Errorf("expected provider %q, got %q", tc.expectedProvider, provider)
				}
				if model != tc.expectedModel {
					t.Errorf("expected model %q, got %q", tc.expectedModel, model)
				}
			}
		})
	}
}

func TestConfigLoadInvalidAdapter(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test-invalid")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configContent := `
agent_adapter: "invalid-format"
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	_, err = LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error when loading config with invalid adapter format, got nil")
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
	if cfg.LLM.Provider != "custom-cli" {
		t.Errorf("Expected LLM Provider to be custom-cli, got %q", cfg.LLM.Provider)
	}
	if cfg.LLM.Model != "custom-model" {
		t.Errorf("Expected LLM Model to be custom-model, got %q", cfg.LLM.Model)
	}
}
