# Configuration

The core configuration for the breakdown is managed via the `config.Config` struct.

```go
type Config struct {
	AgentAdapter string `yaml:"agent_adapter"`
	Atlassian    struct {
		BaseURL string `yaml:"base_url"`
		User    string `yaml:"user"`
		APIKey  string `yaml:"api_key"`
	} `yaml:"atlassian"`
	MaxConcurrency int `yaml:"max_concurrency"`
	MaxRetries     int `yaml:"max_retries"`
}
```

## Configuration File

By default, `breakdown` looks for a YAML configuration file at `~/.breakdown/config.yml`. If this file does not exist, it falls back to a set of in-memory defaults.

An example `config.yml` looks like:

```yaml
agent_adapter: "gemini:google/gemini-3.1-flash-lite-preview"
max_concurrency: 4             # Optional: Maximum concurrent LLM requests (default 4)
max_retries: 3                 # Optional: Maximum retry attempts for failed LLM requests (default 3)
atlassian:
  base_url: "https://your-atlassian-instance.atlassian.net"
  user: "your-email@example.com"
  api_key: "YOUR_API_TOKEN_HERE" # Optional: Can also be passed via BREAKDOWN_ATLASSIAN_API_KEY env var
```

### Agent Adapter Format

The `agent_adapter` format must match `cli:provider/model`.
For example:
- `copilot:anthropic/claude-haiku-4.5`
- `opencode:google/gemini-3.5-flash`
- `gemini:google/gemini-3.1-flash-lite-preview`

No LLM API keys are handled in the Breakdown configuration block itself. Authentication is handled:
- Via standard environment variables (`GEMINI_API_KEY` for gemini, `ANTHROPIC_API_KEY` for claude).
- Independently by the underlying CLI tools (`copilot` or `opencode`).

**Atlassian Integration**
- Optional configuration to automatically fetch content from Jira or Confluence links.
- Requires `base_url`, `user` (email), and `api_key` (Atlassian API token).
- These can be provided in `config.yml` or via environment variables: `BREAKDOWN_ATLASSIAN_BASE_URL`, `BREAKDOWN_ATLASSIAN_API_USER`, and `BREAKDOWN_ATLASSIAN_API_KEY`.
