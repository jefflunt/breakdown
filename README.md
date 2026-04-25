# breakdown

`breakdown` is a non-interactive CLI tool that takes a complex task and breaks it down into an arbitrarily nested structure of sub-tasks. It outputs the result as a folder/file hierarchy on your local filesystem, where leaf nodes are individual files detailing exactly what to do.

## Usage

```bash
# Provide the task as an argument
breakdown "Implement a new feature for the user dashboard"

# Or pipe the task via STDIN
echo "Implement a new feature for the user dashboard" | breakdown
```

The tool will generate a directory structure in `./breakdown-output` representing the plan.

## Configuration

`breakdown` can be configured using a YAML file located at `~/.planner/config.yml`.

Example `~/.planner/config.yml` for Gemini:

```yaml
plans_dir: "~/.planner/plans"
llm:
  provider: "gemini"
  model: "gemini-3.1-flash-lite-preview"
  api_key: "YOUR_API_KEY_HERE"
```

## Atlassian Integration

If configured, `breakdown` automatically fetches content from Jira or Confluence URLs found in tasks or details.

```yaml
atlassian:
  base_url: "https://your-atlassian-instance.atlassian.net"
  user: "your-email@example.com"
  api_key: "YOUR_API_TOKEN_HERE"
```

See `agent_docs/` for more detailed documentation.
