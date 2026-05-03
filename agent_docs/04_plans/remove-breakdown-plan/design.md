# Remove Breakdown Plan Feature

## User Story
**Headline:** Remove unused interactive "plan" features and legacy states from the codebase.
**Problem Statement:** The `breakdown` CLI is designed as a non-interactive task decomposer. However, the codebase contains remnants of an interactive planning feature (such as the unused `--plan` CLI flag, tree manipulation methods like `AddChild`/`EditNode`, plan state saving in `plans/`, and `execute_plan` / `generate_plan_name` LLM prompts). This dead code pollutes the project and creates confusion.
**Objective:** Strip out all remnants of the "breakdown plan" state/interactive tree manipulation capabilities from `main.go`, `pkg/breakdown`, `pkg/llm`, and `prompts/`.
**Expected Outcome:** A leaner codebase where only the core decomposition logic (`Plan`) remains, without `GeneratePlanName`, `GetExecCommand`, or node mutation methods. The test suite will be green and the `prompts/` directory will be cleaner.

## Implementation Backlog

### Pending
- `remove-cli-flag.md`: Remove the unused `--plan` flag and its associated variable from `cmd/breakdown/main.go`.
- `cleanup-prompts.md`: Delete `prompts/generate_plan_name.md` and `prompts/execute_plan.md`, and any leftover state files (e.g. `plans/my-plan.json` and the `plans` directory).
- `simplify-llmclient-interface.md`: Remove `GeneratePlanName` and `GetExecCommand` signatures from the `LLMClient` interface and `ExecRequest` struct in `pkg/breakdown/node.go`.
- `cleanup-llm-implementations.md`: Remove the implementations of `GeneratePlanName` and `GetExecCommand` from all LLM providers in `pkg/llm/` (`claude.go`, `copilot.go`, `gemini.go`, `opencode.go`, `factory.go`) along with their respective tests.
- `cleanup-breakdown-planner.md`: Remove dead tree mutation methods (`EditNode`, `ReplanNode`, `AddChild`, `AddSibling`, `InsertParent`, `DeleteNode`, `SerializePlan`, `GetExecCommand`) from `pkg/breakdown/breakdown.go`, `FormatPlanStructure` from `pkg/breakdown/node.go`, and their corresponding tests in `pkg/breakdown/breakdown_test.go`.

### Current

### Completed

## Architecture Overview
The core `Planner` struct inside `pkg/breakdown` will be simplified to its essential purpose: orchestrating `Start()` and `Plan()` recursively. By eliminating `ReplanNode`, `AddChild`, etc., the `Planner` becomes fundamentally "run-once" rather than acting as a stateful in-memory tree database for a frontend or REPL. 

This change also slims down the `LLMClient` interface, which will now exclusively focus on `AnalyzeTask`, solidifying the boundaries of what the LLMs are expected to do in this workflow.

## Checklist & TDD Requirements
1. **Tests Pass:** All unit tests across `pkg/breakdown` and `pkg/llm` must remain green after removing the code. No test must be left broken by the removal of the mock methods.
2. **Compile Check:** Running `go build ./...` must succeed without unresolved references.
3. **No Dead Code:** Grep for `GetExecCommand`, `GeneratePlanName`, `EditNode`, `AddChild`, etc. should yield 0 results in `.go` source files after the refactor.