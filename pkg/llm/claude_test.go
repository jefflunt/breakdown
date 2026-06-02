package llm

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	breakdown "github.com/jefflunt/breakdown/pkg/breakdown"
)

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestClaudeClient_AnalyzeTask(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	client, err := NewClaudeClient(context.Background(), "claude-3-5-sonnet-latest")
	if err != nil {
		t.Fatalf("expected no error creating client, got: %v", err)
	}

	// Mock the HTTP client
	client.client = &http.Client{
		Transport: &mockTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				responseBody := `{
					"content": [
						{
							"type": "text",
							"text": "{\"action\": \"actionable\", \"reasoning\": \"mock\"}"
						}
					]
				}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
					Header:     make(http.Header),
				}, nil
			},
		},
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
