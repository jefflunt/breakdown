# Planning Workflow

The core mechanic of **breakdown** is its recursive interaction loop, defined in `pkg/breakdown/breakdown.go` within the `Plan(ctx context.Context, node *Node) error` method.

## The Goal

Unlike execution-focused agents that try to *do* the work immediately, this project strictly focuses on decomposition until a specific actionable threshold is met.

**The Actionable Heuristic:** A leaf node is actionable *if and only if* it describes the creation, deletion, or editing of a single file on disk. 
- Example: "Refactor the authentication module" -> *Not Actionable* (Multiple files, too vague)
- Example: "Rename `AuthUser` to `SessionUser` in `src/auth/models.go`" -> *Actionable* (Single file operation)

This ensures that the resulting plan can be executed with extreme predictability.

---

## The Loop (`Plan()` Method)

When `p.Plan(ctx, rootNode)` is called, the orchestrator begins a recursive loop over the tree.

1. **Ask LLM**: `AnalyzeTask(ctx, task)`
   The LLM is prompted with the current task description and the actionable heuristic. It must return a structured `LLMResponse`.

2. **Handle the LLM Response**:

   - **`ActionActionable`**: The LLM determines this is a single file operation. The node's status is set to `actionable`, and this branch terminates successfully.
   
   - **`ActionDecompose`**: The LLM determines the task is too broad. The node's status is set to `composite`, and it generates `N` subtasks. The breakdown creates child `Node`s for each subtask and recursively calls `Plan()` on each child.
   
   - **`ActionAskUser`**: The LLM determines the task is ambiguous (e.g., "Build a web scraper" -> "Which language?"). 

3. **Handling Ambiguity (`ActionAskUser`)**:
   In non-interactive mode, `breakdown` cannot halt to ask for user input. Instead:
   - It marks the node as `actionable`.
   - It appends the LLM's clarification question to the node's `Details` field prefixed with `[Need Input]`.
   - This ensures the generated Markdown file for that task prompts the *human* to clarify or define the requirement when they actually start working on that file.

This loop guarantees that no branch stops growing until it is either clearly `Actionable` or carries a clear instruction for the user to provide clarification.
