# Enable explicit Claude support in Breakdown using the agent adapter config format ('claude:provider/model') by updating the adapter parser inside `pkg/llm/factory.go` to handle the 'claude' case explicitly.

This project integrates support for the standard agent adapter format with Claude within the Breakdown configuration system. Currently, other LLM clients like Copilot and OpenCode use the `agent_adapter` format (e.g., `copilot:provider/model`), and their handling is explicitly implemented inside `parseAgentAdapter`. Claude, however, falls back to the default parsing case.

To ensure robustness, consistency, and clean architecture, we will add an explicit case for `claude` in the `parseAgentAdapter` switch statement inside `pkg/llm/factory.go`. This maps the config value `claude:anthropic/claude-3-5-sonnet` explicitly to the provider name `claude` and the model name `claude-3-5-sonnet`, preventing accidental fallthrough and ensuring future-proof config handling.
