package llm

import (
	"context"
	"fmt"

	breakdown "github.com/jefflunt/breakdown/pkg/breakdown"
	"github.com/jefflunt/breakdown/pkg/config"
)

// NewClient returns a configured LLMClient based on the agentAdapter string.
func NewClient(ctx context.Context, agentAdapter string) (breakdown.LLMClient, error) {
	provider, model, err := parseAgentAdapter(agentAdapter)
	if err != nil {
		return nil, err
	}

	switch provider {
	case "gemini":
		return NewGeminiClient(ctx, model)
	case "copilot":
		return NewCopilotClient(ctx, model)
	case "opencode":
		return NewOpencodeClient(ctx, model)
	case "claude":
		return NewClaudeClient(ctx, model)
	case "mock":
		// Mock primarily used for tests, but can be forced via config
		return &MockClient{}, nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}
}

func parseAgentAdapter(adapter string) (string, string, error) {
	cliName, provider, model, err := config.ParseAdapter(adapter)
	if err != nil {
		return "", "", err
	}

	switch cliName {
	case "opencode":
		return "opencode", provider + "/" + model, nil
	case "copilot":
		return "copilot", model, nil
	default:
		return cliName, model, nil
	}
}

// MockClient is included for fallback and testing
type MockClient struct{}

func (m *MockClient) AnalyzeTask(ctx context.Context, req breakdown.LLMRequest) (breakdown.LLMResponse, error) {
	return breakdown.LLMResponse{
		Action: breakdown.ActionActionable,
	}, nil
}
