# Design Document: Support Copilot CLI and Build-style Configuration

## User Story
* **Headline**: Enable Breakdown to use GitHub Copilot CLI and transition exclusively to build-style agent_adapter configuration.
* **Problem Statement**: Users want to be able to use the `copilot` CLI tool as an LLM provider in Breakdown. They also want to use the build-style `agent_adapter` format (e.g., `cliName:provider/model`) within the configuration file to cleanly define which CLI tool and model to run under the hood, eventually phasing out the legacy, verbose `llm` configuration section and its API key fields entirely, since authentication is handled at the CLI-tool level.
* **Objective**:
  1. Ensure robust support for the `copilot` CLI client.
  2. Implement a config parser option for `agent_adapter` following the exact parsing logic used in `build`.
  3. Seamlessly map the parsed `agent_adapter` components (CLIName, Provider, Model) to Breakdown's provider and model fields.
  4. Once `agent_adapter` works end-to-end, completely remove the legacy `llm` config key and its sub-options (`provider`, `model`, `api_key`), relying on CLI-level authentication.
* **Expected Outcome**: Breakdown can be configured with an `agent_adapter: "copilot:anthropic/claude-haiku-4.5"` or `agent_adapter: "opencode:google/gemini-3.5-flash"`. Legacy `llm` config settings and all LLM API key fields are fully deprecated and removed.

## Implementation Backlog

### Pending
- [ ] Task 1: Add `agent_adapter` field to `config.Config` and implement `ParseAdapter` logic in `pkg/config/`
- [ ] Task 2: Integrate `agent_adapter` parsing into `LoadConfig` in `pkg/config/config.go`
- [ ] Task 3: Map the parsed `agent_adapter` to internal Provider and Model representations
- [ ] Task 4: Verify and refine Copilot CLI integration and unit tests
- [ ] Task 5: Remove the legacy `llm` config key and sub-options, and clean up all files and tests referring to it

### Current
- (None)

### Completed
- (None)

## Architecture Overview
We will add `AgentAdapter` field to the root level of `config.Config`.
During `LoadConfig`, if `AgentAdapter` is provided, we parse it into `CLIName`, `Provider`, and `Model` using the parser rules modeled after `build`.

### Parser Rules
The format must match `cli:provider/model`.
1. Split by `:` into `cliName` and `remaining`.
2. Trim spaces.
3. Split `remaining` by `/` into `provider` and `model`.
4. Trim spaces.
5. If any component is empty, return an error.

### CLI Mapping
- If `CLIName` is `"opencode"`:
  - Internal Provider = `"opencode"`
  - Internal Model = Provider + "/" + Model (since opencode expects provider/model format)
- If `CLIName` is `"copilot"`:
  - Internal Provider = `"copilot"`
  - Internal Model = Model (since copilot CLI expects just model name)
- Otherwise:
  - Internal Provider = CLIName
  - Internal Model = Model

### Transitioning away from legacy `llm` config and API keys
We will remove:
```yaml
llm:
  provider: ...
  model: ...
  api_key: ...
```
And replace it entirely with:
```yaml
agent_adapter: "copilot:anthropic/claude-haiku-4.5"
```
Breakdown will no longer manage, parse, or require any LLM API keys in its configuration, as authentication is handled independently by the underlying CLI tools (`copilot` or `opencode`).

## Checklist & TDD Requirements
1. Unit tests for `ParseAdapter` logic with valid and invalid format inputs.
2. Unit tests for `LoadConfig` verifying that `agent_adapter` correctly populates internal LLM provider and model.
3. Verify that `copilot` client is correctly registered in factory and can execute tasks using mock runner.
4. Verify that removing `llm` configuration works perfectly and that default settings fall back to standard environment variables and sensible defaults.
