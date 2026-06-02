package llm

import (
	"context"
	"fmt"
	"strings"

	breakdown "github.com/jefflunt/breakdown/pkg/breakdown"
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
	if adapter == "" {
		return "", "", fmt.Errorf("empty adapter string")
	}

	parts := strings.SplitN(adapter, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid adapter format: %q (expected 'provider:model')", adapter)
	}

	provider := strings.TrimSpace(parts[0])
	model := strings.TrimSpace(parts[1])

	if provider == "" || model == "" {
		return "", "", fmt.Errorf("invalid adapter: components cannot be empty")
	}

	return provider, model, nil
}

// MockClient is included for fallback and testing
type MockClient struct{}

func (m *MockClient) AnalyzeTask(ctx context.Context, req breakdown.LLMRequest) (breakdown.LLMResponse, error) {
	return breakdown.LLMResponse{
		Action: breakdown.ActionActionable,
	}, nil
}
