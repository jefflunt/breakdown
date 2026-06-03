# Explicitly handle the 'claude' adapter type in the parseAgentAdapter function's switch block within pkg/llm/factory.go to process 'claude:provider/model' configurations cleanly and consistently with other adapter types.

This task involves adding an explicit case for 'claude' within the 'parseAgentAdapter' switch statement inside 'pkg/llm/factory.go'. Currently, any 'claude:provider/model' format is handled by the default case. Making the handling of 'claude' explicit ensures proper model extraction and config mapping matching the style used for 'copilot' and 'opencode' adapters, and provides a clear extension point for any provider-specific parameters in the future.

Once the explicit case is added, we can ensure that calling NewClient with 'claude:provider/model' cleanly parses the adapter configuration and instantiates the Claude client correctly, keeping our adapter configuration fully integrated and aligned with the project's design.
