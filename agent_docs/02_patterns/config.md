# Configuration

The core configuration for the breakdown is managed via the `breakdown.Config` struct.

```go
type Config struct {
	OutputFolder   string `yaml:"output_folder"`
	MaxConcurrency int    `yaml:"max_concurrency"`
	MaxRetries     int    `yaml:"max_retries"`
	// ... (LLM and Atlassian configs)
}
```

## Configuration File

By default, `breakdown` looks for a YAML configuration file at `~/.breakdown/config.yml`. If this file does not exist, it falls back to a set of in-memory defaults.

An example `config.yml` looks like:

```yaml
output_folder: "~/.breakdown/output"
llm:
  provider: "gemini" # Supported: "gemini", "copilot", "opencode", "claude"
  model: "gemini-3.1-flash-lite-preview"
  api_key: "YOUR_API_KEY_HERE" # Optional: Can also be passed via GEMINI_API_KEY or ANTHROPIC_API_KEY env var
atlassian:
  base_url: "https://your-atlassian-instance.atlassian.net"
  user: "your-email@example.com"
  api_key: "YOUR_API_TOKEN_HERE" # Optional: Can also be passed via BREAKDOWN_ATLASSIAN_API_KEY env var
```

### LLM Providers

**Gemini (`provider: "gemini"`)**
- The default provider. Requires `api_key` to be set in the config file or via the `GEMINI_API_KEY` environment variable.
- Configure `model` to pick a specific model (e.g., `gemini-3.1-flash-lite-preview`).

**Anthropic Claude (`provider: "claude"`)**
- Requires `api_key` to be set in the config file or via the `ANTHROPIC_API_KEY` environment variable.
- Configure `model` to pick a specific model (e.g., `claude-3-5-sonnet-latest`).

**GitHub Copilot (`provider: "copilot"`)**
- Requires the `copilot` command line interface to be installed and authenticated (`copilot auth`).
- Does not require an `api_key` in the `breakdown` config since it relies on the CLI's existing session.
- The `model` configuration is optional.

**Opencode (`provider: "opencode"`)**
- Requires the `opencode` command line interface to be installed.
- Does not require an `api_key` in the `breakdown` config since it relies on the CLI's configuration.
- The `model` configuration is optional.

**Atlassian Integration**
- Optional configuration to automatically fetch content from Jira or Confluence links.
- Requires `base_url`, `user` (email), and `api_key` (Atlassian API token).
- These can be provided in `config.yml` or via environment variables: `BREAKDOWN_ATLASSIAN_BASE_URL`, `BREAKDOWN_ATLASSIAN_API_USER`, and `BREAKDOWN_ATLASSIAN_API_KEY`.

## CLI Behavior

1. **Config Path**
   - **Default:** `~/.breakdown/config.yml`

2. **Output Folder**
   - **Default:** `./breakdown-output/`
   - **Usage:** Specifies the directory where the generated plan hierarchy will be created.
