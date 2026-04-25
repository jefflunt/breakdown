# breakdown

`breakdown` is a non-interactive CLI tool designed for complex software engineering tasks. It takes a high-level task description, recursively breaks it down into an arbitrarily nested structure of sub-tasks, and outputs the result as a folder/file hierarchy on your local filesystem.

Leaf nodes in this hierarchy are individual Markdown files, each describing a specific, actionable detail required to complete the task.

## Getting Started

### Installation

`breakdown` relies on a standard Go workspace. To build and install:

```bash
# Clone the repository and navigate to the root
cd breakdown

# Build, test, and install
./script/build-test-install
```

### Configuration

`breakdown` requires a configuration file located at `~/.breakdown/config.yml`.

Example `~/.breakdown/config.yml`:

```yaml
output_folder: "~/.breakdown/output" # Directory where the plan will be generated
llm:
  provider: "gemini" # Supported: "gemini", "copilot", "opencode", "claude"
  model: "gemini-3.1-flash-lite-preview"
  api_key: "YOUR_API_KEY_HERE"
```

## Usage

```bash
# Provide the task as an argument
breakdown "Implement a TUI TODO application"

# Or pipe the task via STDIN
echo "Implement a TUI TODO application" | breakdown

# Use -v to see the decomposition progress in real-time
breakdown -v "Implement a TUI TODO application"
```

The tool will generate a directory structure in the directory specified by `output_folder` (defaulting to `./breakdown-output`) representing the plan.

## Mindset for Users & AI Agents

`breakdown` is designed to be an **agentic task orchestrator**. 

1.  **Iterative Decomposition:** The tool treats tasks as trees. If a task is too complex, it decomposes it into subtasks. It repeats this process until every branch of the task tree reaches an `Actionable` state (typically a single-file edit).
2.  **Context-Awareness:** Before analyzing a task, `breakdown` inspects your current directory's file system tree. This provides necessary context to the LLM about the project's language, framework, and architecture.
3.  **Actionable Output:** The goal is not just a plan, but a roadmap for execution. Each leaf node represents a concrete instruction.
4.  **Agent-First Docs:** If you are an AI agent, please check the `agent_docs/` directory. It contains architectural overviews, design patterns, and deep dives into the `breakdown` planning loop.

## Atlassian Integration

If configured, `breakdown` can automatically fetch context from Jira or Confluence URLs found in tasks or details.

```yaml
atlassian:
  base_url: "https://your-atlassian-instance.atlassian.net"
  user: "your-email@example.com"
  api_key: "YOUR_API_TOKEN_HERE"
```
