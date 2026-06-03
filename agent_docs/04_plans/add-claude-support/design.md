# Design Document: Support Claude using Agent Adapter Config Format

## User Story
* **Headline**: Enable Claude support in Breakdown using the build-style `agent_adapter` format.
* **Problem Statement**: We recently transitioned to using the build-style `agent_adapter` format (e.g. `cliName:provider/model`) for providers like Copilot and OpenCode. Currently, Claude uses standard API client configuration, but does not explicitly handle the `claude:provider/model` adapter format in `parseAgentAdapter`'s switch statement. We want to align Claude support with the newly introduced `agent_adapter` format so users can seamlessly configure and use Claude.
* **Objective**:
  1. Add explicit support for `claude:provider/model` (e.g. `claude:anthropic/claude-3-5-sonnet`) in `parseAgentAdapter`.
  2. Implement comprehensive unit tests to verify parsing and initialization of the `ClaudeClient` under the `agent_adapter` configuration.
  3. Ensure compatibility with the existing Anthropic API client.
* **Expected Outcome**: Breakdown can be configured with `agent_adapter: "claude:anthropic/claude-3-5-sonnet"`, which resolves to initializing `ClaudeClient` with the model `claude-3-5-sonnet`.

## Implementation Backlog

### Pending
- [ ] `01-Explicitly_handle_the_'claude'_adapter_type_in_the.md`: Explicitly handle 'claude' adapter type in `parseAgentAdapter`.

### Current
- (None)

### Completed
- (None)

## Architecture Overview
In `pkg/llm/factory.go`, `parseAgentAdapter` processes the `agent_adapter` configuration string. We will explicitly handle the `claude` cliName, extracting the underlying `model` name and returning `"claude"` as the provider name.

```go
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
	case "claude":
		return "claude", model, nil
	default:
		return cliName, model, nil
	}
}
```

This maps cleanly to `NewClient`'s switch statement:
```go
	case "claude":
		return NewClaudeClient(ctx, model)
```

## Checklist & TDD Requirements
1. Unit test in `pkg/llm/factory_test.go` ensuring that `claude:anthropic/claude-3-5-sonnet` correctly parses into provider `"claude"` and model `"claude-3-5-sonnet"`.
2. Unit test in `pkg/llm/factory_test.go` ensuring that `NewClient` correctly returns a `ClaudeClient` when configured with `claude:anthropic/claude-3-5-sonnet` and the API key is present.
