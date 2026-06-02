package llm

import (
	"context"
	"os/exec"
	"testing"

	breakdown "github.com/jefflunt/breakdown/pkg/breakdown"
	"github.com/jefflunt/breakdown/pkg/config"
)

func TestCopilotClient_AnalyzeTask(t *testing.T) {
	// Skip the test if the copilot CLI is not installed
	if _, err := exec.LookPath("copilot"); err != nil {
		t.Skip("copilot CLI not found in PATH; skipping test")
	}

	cfg := &config.Config{}
	cfg.LLM.Provider = "copilot"

	client, err := NewCopilotClient(context.Background(), cfg)
	if err != nil {
		t.Fatalf("expected no error creating client, got: %v", err)
	}

	client.runner = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		return []byte(`{"action": "actionable", "reasoning": "mock"}`), nil, nil
	}

	req := breakdown.LLMRequest{
		Task: "Create a simple python hello world script in hello.py",
	}

	resp, err := client.AnalyzeTask(context.Background(), req)
	if err != nil {
		t.Fatalf("AnalyzeTask failed: %v", err)
	}

	if resp.Action != "actionable" {
		t.Fatalf("expected actionable, got %v", resp.Action)
	}
}
