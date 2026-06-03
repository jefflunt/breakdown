# Architecture

The **breakdown** application is fundamentally a tree data structure orchestrator. It is built to recursively decompose a root task into an N-ary tree of subtasks.

## Core Concepts

### `breakdown.Planner`
The `Planner` struct is the central orchestrator. It manages the `Root` node, configuration state, and the LLM client, using a semaphore channel to enforce concurrency limits on LLM calls.

```go
type Planner struct {
	mu           sync.RWMutex
	Root         *Node         `json:"root"`
	Config       Config        `json:"config"`
	LLM          LLMClient     `json:"-"`
	llmSemaphore chan struct{} `json:"-"`
}
```

### `breakdown.Node`
A node represents a single task in the task tree.
- `ID`: UUID for tracking.
- `Task`: A string describing the work. This string grows as user clarifications are appended to it.
- `Type`: EITHER `TaskTypeAtomic` (a leaf node) or `TaskTypeComposite` (a node with children).
- `Status`: 
  - `pending` (waiting to be analyzed)
  - `composite` (analyzed, has children)
  - `actionable` (analyzed, single file operation)
  - `needs_input` (waiting for user clarification)

### Filesystem Generation

The core breakdown is entirely decoupled from user interactivity. The primary output mechanism is the `GenerateFilesystemStructure` method on the `Node`.

Once planning is complete, this method recursively traverses the task tree to generate a filesystem hierarchy:
- **Composite Nodes** become folders prefixed with a numeric ID (e.g., `01-task-name`).
- **Actionable Nodes** become Markdown files (e.g., `01-task-name.md`).

This ensures the plan is readily usable in any text editor or IDE.
