package breakdown

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jefflunt/breakdown/pkg/atlassian"
	"github.com/jefflunt/breakdown/pkg/logger"
	"golang.org/x/sync/errgroup"
)

// Config represents the planner configuration
type Config struct {
	Workspace      string // Directory to hold workspaces
	MaxConcurrency int
	MaxRetries     int
	Verbose        bool
	AgentAdapter   string
	Atlassian      struct {
		BaseURL string
		User    string
		APIKey  string
	}
}

// Planner is the central orchestrator for the task tree
type Planner struct {
	mu           sync.RWMutex
	Root         *Node         `json:"root"`
	Config       Config        `json:"config"`
	LLM          LLMClient     `json:"-"`
	llmSemaphore chan struct{} `json:"-"`
}

// NewPlanner creates a new planner instance
func NewPlanner(cfg Config, llm LLMClient) *Planner {
	if cfg.MaxConcurrency == 0 {
		cfg.MaxConcurrency = 4
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}

	return &Planner{
		Config:       cfg,
		LLM:          llm,
		llmSemaphore: make(chan struct{}, cfg.MaxConcurrency),
	}
}

func (p *Planner) RLock() {
	p.mu.RLock()
}

func (p *Planner) RUnlock() {
	p.mu.RUnlock()
}

// SerializePlan returns the plan as a string representation.

// Start initiates the planning process for a root task
func (p *Planner) Start(ctx context.Context, task string) error {
	p.mu.Lock()
	p.Root = &Node{
		ID:     uuid.New().String(),
		Task:   task,
		Status: StatusPending,
		Depth:  0,
	}
	p.mu.Unlock()

	err := p.Plan(ctx, p.Root)
	return err
}

// analyzeTaskWithRetry wraps the LLM call with a semaphore and retry logic.
func (p *Planner) analyzeTaskWithRetry(ctx context.Context, req LLMRequest) (LLMResponse, error) {
	p.llmSemaphore <- struct{}{}
	defer func() { <-p.llmSemaphore }()

	var lastErr error
	for i := 0; i <= p.Config.MaxRetries; i++ {
		resp, err := p.LLM.AnalyzeTask(ctx, req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond) // Simple backoff
	}
	return LLMResponse{}, fmt.Errorf("failed after %d retries: %w", p.Config.MaxRetries, lastErr)
}

// ReplanNode clears a node's children, resets its status, and saves.

// AddChild adds a new child node to the specified parent.

// AddSibling adds a new node immediately before or after the specified sibling.

// InsertParent inserts a new node directly above the target node.
// The target node becomes the only child of the new node.

// DeleteNode removes a node and all its children from the tree.

// Find finds a node by its ID
func (p *Planner) Find(id string) *Node {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.Root == nil {
		return nil
	}
	return p.Root.Find(id)
}

// GetAncestry returns a list of parent tasks, from Root down to the immediate parent of the given node.
func (p *Planner) GetAncestry(node *Node) []string {
	var ancestry []string
	current := node

	for current != nil && current.ParentID != "" {
		parent := p.Find(current.ParentID)
		if parent == nil {
			break
		}
		// Prepend
		ancestry = append([]string{parent.Task}, ancestry...)
		current = parent
	}
	return ancestry
}

// Plan recursively decomposes a node, polling the LLM until Actionable
func (p *Planner) Plan(ctx context.Context, node *Node) error {
	p.mu.RLock()
	status := node.Status
	children := node.Children
	isRoot := p.Root != nil && node.ID == p.Root.ID
	p.mu.RUnlock()

	// 1. Fetch Atlassian Content if URLs present
	if p.Config.Atlassian.BaseURL != "" && p.Config.Atlassian.APIKey != "" {
		client := atlassian.NewClient(p.Config.Atlassian.BaseURL, p.Config.Atlassian.User, p.Config.Atlassian.APIKey)

		textToScan := node.Task + "\n" + node.Details

		// Simple URL detection
		lines := strings.Split(textToScan, " ")
		for _, token := range lines {
			if strings.Contains(token, p.Config.Atlassian.BaseURL) {
				// Clean URL (remove trailing punctuation)
				url := strings.TrimRight(token, ".,!?)")
				content, err := client.Fetch(url)
				if err == nil {
					p.mu.Lock()
					node.Details = fmt.Sprintf("%s\n\n[Atlassian Context from %s]:\n%s", node.Details, url, content)
					p.mu.Unlock()
				}
			}
		}
	}

	// If already fully analyzed, just skip or recurse
	if status == StatusActionable {
		return nil
	}

	if p.Config.Verbose && node.Depth > 0 {
		fmt.Fprintf(os.Stderr, "%s%s\n", strings.Repeat("  ", node.Depth), node.Task)
	}

	if status == StatusComposite {
		g, gCtx := errgroup.WithContext(ctx)
		for _, child := range children {
			c := child
			g.Go(func() error {
				if err := p.Plan(gCtx, c); err != nil {
					p.mu.Lock()
					c.Status = StatusError
					c.Task = "!" + c.Task
					p.mu.Unlock()
					_ = logger.Log(err)
				}
				return nil
			})
		}
		return g.Wait()
	}

	pwd, _ := os.Getwd()
	fsTree := GetFileSystemTree(pwd)

	for {
		p.mu.Lock()
		node.Status = StatusPending
		p.mu.Unlock()

		ancestry := p.GetAncestry(node)

		req := LLMRequest{
			Task:           node.Task,
			Ancestry:       ancestry,
			IsVision:       isRoot,
			FileSystemTree: fsTree,
		}

		// Ask LLM what to do
		resp, err := p.analyzeTaskWithRetry(ctx, req)
		if err != nil {
			p.mu.Lock()
			node.Status = StatusError
			node.Task = "!" + node.Task
			p.mu.Unlock()
			_ = logger.Log(err)
			return fmt.Errorf("failed to analyze task %q: %w", node.Task, err)
		}

		p.mu.Lock()
		if resp.RewrittenTask != "" && resp.RewrittenTask != node.Task {
			node.Task = resp.RewrittenTask
		}
		if resp.Title != "" {
			node.Title = resp.Title
		}
		if resp.Details != "" {
			node.Details = resp.Details
		}
		if resp.AsciiDiagram != "" {
			node.AsciiDiagram = resp.AsciiDiagram
		}
		p.mu.Unlock()

		switch resp.Action {
		case ActionActionable:
			p.mu.Lock()
			node.Type = TaskTypeAtomic
			node.Status = StatusActionable
			p.mu.Unlock()
			return nil // Branch terminates successfully
		case ActionDecompose:
			p.mu.Lock()
			node.Type = TaskTypeComposite
			node.Status = StatusComposite
			for _, st := range resp.Subtasks {
				child := &Node{
					ID:       uuid.New().String(),
					ParentID: node.ID,
					Task:     st,
					Status:   StatusPending,
					Depth:    node.Depth + 1,
				}
				node.Children = append(node.Children, child)
			}
			p.mu.Unlock()

			// Recursively plan children
			g, gCtx := errgroup.WithContext(ctx)
			p.mu.RLock()
			currentChildren := node.Children
			p.mu.RUnlock()
			for _, child := range currentChildren {
				c := child
				g.Go(func() error {
					if err := p.Plan(gCtx, c); err != nil {
						p.mu.Lock()
						c.Status = StatusError
						c.Task = "!" + c.Task
						p.mu.Unlock()
						_ = logger.Log(err)
					}
					return nil
				})
			}
			return g.Wait()

		case ActionAskUser:
			// In non-interactive mode, we can't ask user.
			// Just treat as Actionable for now, or log an error.
			p.mu.Lock()
			node.Status = StatusActionable
			node.Type = TaskTypeAtomic
			if node.Details == "" {
				node.Details = fmt.Sprintf("[Need Input]: %s", resp.Question)
			} else {
				node.Details = fmt.Sprintf("%s\n\n[Need Input]: %s", node.Details, resp.Question)
			}
			p.mu.Unlock()
			return nil
		}
		// If we reach here (and it wasn't Actionable), the loop continues
		// to re-analyze if the node was not Actionable.
		// For ActionAskUser, we returned.
	}
}

// GetExecCommand gathers the node's context and the overall plan structure,
// then returns an un-started exec.Cmd that will execute the plan natively.
