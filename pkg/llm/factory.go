package llm

import (
	"context"
	"fmt"
	"os/exec"

	 "github.com/jefflunt/breakdown/pkg/config"
	breakdown  "github.com/jefflunt/breakdown/pkg/breakdown"
)

// NewClient returns a configured LLMClient based on the config provider.
func NewClient(ctx context.Context, cfg *config.Config) (breakdown.LLMClient, error) {
	switch cfg.LLM.Provider {
	case "gemini":
		return NewGeminiClient(ctx, cfg)
	case "copilot":
		return NewCopilotClient(ctx, cfg)
	case "opencode":
		return NewOpencodeClient(ctx, cfg)
	case "claude":
		return NewClaudeClient(ctx, cfg)
	case "mock":
		// Mock primarily used for tests, but can be forced via config
		return &MockClient{}, nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", cfg.LLM.Provider)
	}
}

// MockClient is included for fallback and testing
type MockClient struct{}

func (m *MockClient) AnalyzeTask(ctx context.Context, req breakdown.LLMRequest) (breakdown.LLMResponse, error) {
	return breakdown.LLMResponse{
		Action: breakdown.ActionActionable,
	}, nil
}

func (m *MockClient) GeneratePlanName(ctx context.Context, task string) (string, error) {
	return "mock-plan-name", nil
}

func (m *MockClient) GetExecCommand(ctx context.Context, req breakdown.ExecRequest) (*exec.Cmd, error) {
	return exec.Command("echo", "mock execution"), nil
}
