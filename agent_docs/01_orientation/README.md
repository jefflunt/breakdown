# breakdown — Agent Documentation Index

**breakdown** is a non-interactive recursive task orchestrator written in Go. It accepts a high-level task description, recursively breaks it down iteratively using an LLM, and guarantees that the resulting leaf nodes represent *Logical Units of Work (LUoW)* by recursively polling the LLM for decomposition until every branch reaches an actionable state.

Unlike previous versions, `breakdown` is non-interactive: it generates the plan structure as a directory/file hierarchy on your local filesystem, where leaf nodes are actionable Markdown files.

---

## Quick Start

### Installation

```bash
# Build, test, and install the breakdown binary
./script/build-test-install
```

### Usage
Write your task description, requirements, or architecture ideas into a text/markdown file, then run:

```bash
breakdown my-prompt.md
# The plan is generated in ./breakdown-output/
```

---

## How to Use This Documentation

This folder follows **Progressive Disclosure** principles. Start here, then read only the detail files relevant to your task.

| File | What it covers | Read when… |
|------|---------------|------------|
| **This file** | Repo overview, file map, key facts | Always — start here |
| [`architecture.md`](../02_patterns/architecture.md) | The generic task tree, file generator logic | Changing core logic, adding new generator components |
| [`planning_workflow.md`](../03_deep_dives/planning_workflow.md) | Step-by-step walkthrough of how tasks are analyzed and decomposed | Changing the LLM interaction loop |
| [`building.md`](../02_patterns/building.md) | Build process and commands | Building the binaries |
| [`config.md`](../02_patterns/config.md) | Configuration options (output folder, LLM settings) | Changing CLI flags or configuration |

---

## Repo at a Glance

```
breakdown/
├── bin/                       ← Compiled output directory
├── cmd/
│   └── breakdown/             ← CLI entry point
├── pkg/
│   ├── atlassian/             ← Jira/Confluence integration client
│   ├── breakdown/             ← Core orchestrator & filesystem generator logic
│   ├── config/                ← YAML Configuration parsing
│   ├── llm/                   ← Gemini/Claude/etc LLM Clients
│   ├── logger/                ← Internal logger utilities
│   └── version/               ← Binary version definitions
├── script/                    ← Build, test, and automation scripts
└── agent_docs/                ← this documentation tree
```

**Module:** `github.com/jefflunt/breakdown`  
**Go version:** 1.26+  

---

## Key Facts

- **Non-Interactive:** Once invoked, `breakdown` proceeds to generate the full filesystem hierarchy without further user input.
- **Actionable Heuristic:** A leaf node is *only* actionable if it describes a cohesive "Logical Unit of Work" (LUoW) that can be implemented as a functional slice, spanning one or multiple files. The LLM must enforce this.
- **Test-Last Pipeline Synergy:** `breakdown` is designed to work with the `build` orchestrator's `Dev -> Tester` pipeline. It explicitly avoids generating testing tasks, allowing the downstream `Tester` agent to handle QA natively.
- **No Max Depth:** The breakdown does not rely on arbitrary depth limits. It continues to decompose infinitely until the LLM returns `Actionable` for all branches.
- **Filesystem Output:** The final plan is generated as a directory/file hierarchy in the directory specified by `output_folder` (default `./breakdown-output`).
- **Atlassian Integration:** If configured, the breakdown automatically fetches content from Jira or Confluence URLs found in tasks or details, providing the LLM with direct access to your issue tracking and documentation.
