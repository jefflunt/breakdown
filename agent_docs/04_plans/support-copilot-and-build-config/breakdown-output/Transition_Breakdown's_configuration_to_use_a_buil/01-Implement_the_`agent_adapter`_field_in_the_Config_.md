# Implement the `agent_adapter` field in the Config struct, create a `ParseAdapter` parser to decompose the `cliName:provider/model` string format, update `LoadConfig` to utilize this new format, and remove legacy `llm` configuration block fields within `pkg/config/config.go`.

This task involves updating the configuration layer in `pkg/config/config.go` to support the new build-style `agent_adapter` format. You will add the `agent_adapter` field to the main Config struct (and remove the legacy `llm` block configuration fields). 

Next, you will implement a `ParseAdapter` function (or similar logic) to handle the string format `cliName:provider/model` (e.g., `copilot:anthropic/claude-haiku-4.5`). This function should extract the CLI tool name, the LLM provider, and the specific model. Finally, update `LoadConfig` to correctly read this new field from the environment or configuration files, parse it, and populate the internal configuration state for use by the LLM factory downstream.
