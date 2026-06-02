package llm

import (
	"context"
	"os/exec"
	"testing"

	breakdown "github.com/jefflunt/breakdown/pkg/breakdown"
)

func TestOpencodeClient_AnalyzeTask(t *testing.T) {
	// Skip the test if the opencode CLI is not installed
	if _, err := exec.LookPath("opencode"); err != nil {
		t.Skip("opencode CLI not found in PATH; skipping test")
	}

	client, err := NewOpencodeClient(context.Background(), "google/gemini-2.5-pro")
	if err != nil {
		t.Fatalf("expected no error creating client, got: %v", err)
	}

	client.runner = func(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
		// Mock line-by-line JSON stream of opencode events
		response := `{"type": "text", "part": {"text": "{\"action\":"}}
{"type": "text", "part": {"text": " \"actionable\", \"reasoning\":"}}
{"type": "text", "part": {"text": " \"mock\"}"}}`
		return []byte(response), nil, nil
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
