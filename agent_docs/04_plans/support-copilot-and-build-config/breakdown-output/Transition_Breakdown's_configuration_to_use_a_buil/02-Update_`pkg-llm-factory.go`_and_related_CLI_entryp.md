# Update `pkg/llm/factory.go` and related CLI entrypoints to parse and utilize the new `agent_adapter` configuration format, correctly routing to and initializing the appropriate LLM providers like `copilot` or `opencode`.

Update `pkg/llm/factory.go` to accept and utilize the new `AgentAdapter` configuration format (e.g., `copilot:anthropic/claude-haiku-4.5`). The factory must parse this string to separate the provider prefix (e.g., `copilot`, `opencode`) from the target model identifier. Based on the parsed provider, the factory should instantiate the corresponding internal client (Copilot or OpenCode) and pass the model identifier to it.

After updating the factory logic, modify the CLI entrypoints, specifically `cmd/breakdown/main.go` and `pkg/breakdown/breakdown.go`, to ensure they are passing the new `AgentAdapter` configuration property to the factory during setup, effectively replacing any transitional usage of the legacy configuration mapping.
