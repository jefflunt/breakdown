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

func TestConfigLoadWithLegacyLLMFields(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test-legacy")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configContent := `
agent_adapter: "copilot:anthropic/claude-3-5-sonnet"
llm:
  provider: "gemini"
  model: "gemini-1.5-pro"
  api_key: "some-api-key"
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed with legacy fields: %v", err)
	}

	if cfg.AgentAdapter != "copilot:anthropic/claude-3-5-sonnet" {
		t.Errorf("Expected AgentAdapter 'copilot:anthropic/claude-3-5-sonnet', got %s", cfg.AgentAdapter)
	}
}

func TestConfigLoadWithEmptyAdapter(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config-test-empty-adapter")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configContent := `
agent_adapter: ""
`
	configPath := filepath.Join(tempDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Make sure env is clear to test pure default fallback
	t.Setenv("BREAKDOWN_AGENT_ADAPTER", "")

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.AgentAdapter != "gemini:google/gemini-3.1-flash-lite-preview" {
		t.Errorf("Expected fallback default AgentAdapter gemini:google/gemini-3.1-flash-lite-preview, got %s", cfg.AgentAdapter)
	}
}

func TestParseAdapter(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedCliName  string
		expectedProvider string
		expectedModel    string
		expectErr        bool
	}{
		{
			name:             "valid copilot",
			input:            "copilot:anthropic/claude-haiku-4.5",
			expectedCliName:  "copilot",
			expectedProvider: "anthropic",
			expectedModel:    "claude-haiku-4.5",
			expectErr:        false,
		},
		{
			name:             "valid opencode",
			input:            "opencode:google/gemini-3.5-flash",
			expectedCliName:  "opencode",
			expectedProvider: "google",
			expectedModel:    "gemini-3.5-flash",
			expectErr:        false,
		},
		{
			name:             "valid with spaces and trimming",
			input:            "  gemini : google / gemini-3.1-flash-lite-preview  ",
			expectedCliName:  "gemini",
			expectedProvider: "google",
			expectedModel:    "gemini-3.1-flash-lite-preview",
			expectErr:        false,
		},
		{
			name:      "empty input",
			input:     "",
			expectErr: true,
		},
		{
			name:      "missing colon",
			input:     "copilot",
			expectErr: true,
		},
		{
			name:      "missing slash",
			input:     "copilot:anthropic-claude",
			expectErr: true,
		},
		{
			name:      "empty components",
			input:     " : / ",
			expectErr: true,
		},
		{
			name:      "empty cli",
			input:     ":anthropic/claude-haiku",
			expectErr: true,
		},
		{
			name:      "empty provider",
			input:     "copilot:/claude-haiku",
			expectErr: true,
		},
		{
			name:      "empty model",
			input:     "copilot:anthropic/",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cliName, provider, model, err := ParseAdapter(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tc.input, err)
				}
				if cliName != tc.expectedCliName {
					t.Errorf("expected cliName %q, got %q", tc.expectedCliName, cliName)
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
