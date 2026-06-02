package llm

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestParseAgentAdapter(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedProvider string
		expectedModel    string
		expectErr        bool
	}{
		{
			name:             "valid copilot format",
			input:            "copilot:anthropic/claude-haiku-4.5",
			expectedProvider: "copilot",
			expectedModel:    "claude-haiku-4.5",
			expectErr:        false,
		},
		{
			name:             "valid opencode format",
			input:            "opencode:google/gemini-2.5-pro",
			expectedProvider: "opencode",
			expectedModel:    "google/gemini-2.5-pro",
			expectErr:        false,
		},
		{
			name:             "valid mock format",
			input:            "mock:mock-provider/anything",
			expectedProvider: "mock",
			expectedModel:    "anything",
			expectErr:        false,
		},
		{
			name:             "valid gemini format",
			input:            "gemini:google/model-name",
			expectedProvider: "gemini",
			expectedModel:    "model-name",
			expectErr:        false,
		},
		{
			name:             "valid with spaces",
			input:            "  copilot  :  anthropic/claude-haiku-4.5  ",
			expectedProvider: "copilot",
			expectedModel:    "claude-haiku-4.5",
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
			name:      "empty model name",
			input:     "copilot:anthropic/",
			expectErr: true,
		},
		{
			name:      "empty provider name",
			input:     "copilot:/claude",
			expectErr: true,
		},
		{
			name:      "empty provider and model",
			input:     "copilot:/",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			provider, model, err := parseAgentAdapter(tc.input)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error for input %q, got nil", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tc.input, err)
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

func TestNewClient(t *testing.T) {
	ctx := context.Background()

	// 1. Mock provider
	t.Run("mock client", func(t *testing.T) {
		client, err := NewClient(ctx, "mock:mock-provider/some-model")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := client.(*MockClient); !ok {
			t.Errorf("expected *MockClient, got %T", client)
		}
	})

	// 2. Unknown provider
	t.Run("unknown client", func(t *testing.T) {
		_, err := NewClient(ctx, "unknown:provider/model")
		if err == nil {
			t.Fatal("expected error for unknown provider, got nil")
		}
		if !strings.Contains(err.Error(), "unknown LLM provider") {
			t.Errorf("expected error to contain 'unknown LLM provider', got: %v", err)
		}
	})

	// 3. Claude client without API key
	t.Run("claude client without API key", func(t *testing.T) {
		// Clean env for the test
		origKey, exists := os.LookupEnv("ANTHROPIC_API_KEY")
		if exists {
			os.Unsetenv("ANTHROPIC_API_KEY")
			defer os.Setenv("ANTHROPIC_API_KEY", origKey)
		}

		_, err := NewClient(ctx, "claude:anthropic/claude-3-5-sonnet")
		if err == nil {
			t.Fatal("expected error without ANTHROPIC_API_KEY, got nil")
		}
		if !strings.Contains(err.Error(), "anthropic api key is required") {
			t.Errorf("expected key error, got: %v", err)
		}
	})

	// 4. Claude client with API key
	t.Run("claude client with API key", func(t *testing.T) {
		origKey, exists := os.LookupEnv("ANTHROPIC_API_KEY")
		os.Setenv("ANTHROPIC_API_KEY", "dummy-key")
		defer func() {
			if exists {
				os.Setenv("ANTHROPIC_API_KEY", origKey)
			} else {
				os.Unsetenv("ANTHROPIC_API_KEY")
			}
		}()

		client, err := NewClient(ctx, "claude:anthropic/claude-3-5-sonnet")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := client.(*ClaudeClient); !ok {
			t.Errorf("expected *ClaudeClient, got %T", client)
		}
	})

	// 5. Gemini client without API key
	t.Run("gemini client without API key", func(t *testing.T) {
		origKey, exists := os.LookupEnv("GEMINI_API_KEY")
		if exists {
			os.Unsetenv("GEMINI_API_KEY")
			defer os.Setenv("GEMINI_API_KEY", origKey)
		}

		_, err := NewClient(ctx, "gemini:google/gemini-1.5-pro")
		if err == nil {
			t.Fatal("expected error without GEMINI_API_KEY, got nil")
		}
		if !strings.Contains(err.Error(), "gemini api key is required") {
			t.Errorf("expected key error, got: %v", err)
		}
	})

	// 6. Gemini client with API key
	t.Run("gemini client with API key", func(t *testing.T) {
		origKey, exists := os.LookupEnv("GEMINI_API_KEY")
		os.Setenv("GEMINI_API_KEY", "dummy-key")
		defer func() {
			if exists {
				os.Setenv("GEMINI_API_KEY", origKey)
			} else {
				os.Unsetenv("GEMINI_API_KEY")
			}
		}()

		client, err := NewClient(ctx, "gemini:google/gemini-1.5-pro")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := client.(*GeminiClient); !ok {
			t.Errorf("expected *GeminiClient, got %T", client)
		}
	})

	// 7. Copilot client tests (checks executable presence)
	t.Run("copilot client", func(t *testing.T) {
		_, lookupErr := exec.LookPath("copilot")

		client, err := NewClient(ctx, "copilot:provider/some-model")
		if lookupErr != nil {
			if err == nil {
				t.Fatal("expected error because copilot CLI is not in PATH, got nil")
			}
			if !strings.Contains(err.Error(), "copilot command line interface not found in PATH") {
				t.Errorf("expected command not found error, got: %v", err)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error when copilot CLI is in PATH: %v", err)
			}
			if _, ok := client.(*CopilotClient); !ok {
				t.Errorf("expected *CopilotClient, got %T", client)
			}
		}
	})

	// 8. Opencode client tests (checks executable presence)
	t.Run("opencode client", func(t *testing.T) {
		_, lookupErr := exec.LookPath("opencode")

		client, err := NewClient(ctx, "opencode:provider/some-model")
		if lookupErr != nil {
			if err == nil {
				t.Fatal("expected error because opencode CLI is not in PATH, got nil")
			}
			if !strings.Contains(err.Error(), "opencode command line interface not found in PATH") {
				t.Errorf("expected command not found error, got: %v", err)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error when opencode CLI is in PATH: %v", err)
			}
			if _, ok := client.(*OpencodeClient); !ok {
				t.Errorf("expected *OpencodeClient, got %T", client)
			}
		}
	})

	// 9. Parsing error propagates
	t.Run("invalid adapter string format", func(t *testing.T) {
		_, err := NewClient(ctx, "invalid_adapter_format")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
