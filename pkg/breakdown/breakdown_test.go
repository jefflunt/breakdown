package breakdown

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// simpleMockClient replaces the complex MockClient from pkg/llm for predictable testing
type simpleMockClient struct {
	responses map[string]LLMResponse
}

func (m *simpleMockClient) AnalyzeTask(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	if resp, ok := m.responses[req.Task]; ok {
		return resp, nil
	}
	return LLMResponse{}, fmt.Errorf("unexpected task: %s", req.Task)
}





func TestPlannerStart(t *testing.T) {
	cfg := Config{}

	client := &simpleMockClient{
		responses: map[string]LLMResponse{
			"Do a simple task": {Action: ActionActionable},
		},
	}

	p := NewPlanner(cfg, client)
	ctx := context.Background()

	// Starting should initialize the root node and trigger Plan()
	err := p.Start(ctx, "Do a simple task")
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	if p.Root == nil {
		t.Fatalf("Root node was not created")
	}

	if p.Root.Task != "Do a simple task" {
		t.Errorf("Expected root task to be 'Do a simple task'")
	}

	if p.Root.Status != StatusActionable {
		t.Errorf("Expected root status to be Actionable")
	}
}

func TestPlannerPlanDecomposition(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "breakdown-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	
	cfg := Config{ }

	client := &simpleMockClient{
		responses: map[string]LLMResponse{
			"Complex Task": {
				Action:   ActionDecompose,
				Subtasks: []string{"Subtask 1", "Subtask 2"},
			},
			"Subtask 1": {Action: ActionActionable},
			"Subtask 2": {Action: ActionActionable},
		},
	}

	p := NewPlanner(cfg, client)
	p.Root = &Node{ID: "root", Task: "Complex Task", Status: StatusPending}

	err = p.Plan(context.Background(), p.Root)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	if len(p.Root.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(p.Root.Children))
	}

	if p.Root.Children[0].Status != StatusActionable || p.Root.Children[1].Status != StatusActionable {
		t.Errorf("Expected children to be actionable")
	}
}

func TestPlannerPlanAskUser(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "breakdown-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	
	cfg := Config{ }

	client := &simpleMockClient{
		responses: map[string]LLMResponse{
			"Ambiguous task": {
				Action:   ActionAskUser,
				Question: "What do you mean?",
			},
		},
	}

	p := NewPlanner(cfg, client)
	p.Root = &Node{ID: "root", Task: "Ambiguous task", Status: StatusPending}

	// Run planning in a separate goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- p.Plan(context.Background(), p.Root)
	}()

	// Wait for planning to finish
	select {
	case err := <-errChan:
		if err != nil {
			t.Fatalf("Plan failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out waiting for Plan to finish")
	}

	// Verify the node's Details string was updated
	if !strings.Contains(p.Root.Details, "[Need Input]") {
		t.Errorf("Expected details to contain [Need Input], got %q", p.Root.Details)
	}

	if p.Root.Status != StatusActionable {
		t.Errorf("Expected status to be actionable, got %s", p.Root.Status)
	}
}








