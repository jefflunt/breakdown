# User Story
- **Headline**: Integrate the Breakdown Design and Build Workflow into `breakdown`
- **Problem Statement**: The `breakdown` tool generates excellent hierarchical task trees, but utilizing them effectively requires a specific, disciplined workflow (design drafting, auditing, strict TDD, and a user-approval loop). Users and AI agents need a codified definition of this workflow to maximize the value of the tool.
- **Objective**: Formalize the highly successful "pair programming" rhythm developed during the `money` CLI project into an official agent definition file (`breakdown-design-and-build.md`). This agent profile will guide AI assistants on exactly how to wield the `breakdown` tool and execute the resulting plans.
- **Expected Outcome**: A polished `breakdown-design-and-build.md` file residing within the `breakdown` repository, ready to be used as a system prompt or documentation for agents interacting with `breakdown`.

---

# Implementation Backlog

## Pending
- [DOCS] Write the final version of `breakdown-design-and-build.md` into the `breakdown/agent_docs/` directory, incorporating the improvement to use `breakdown`'s direct output folder argument.

## Current
*(Empty)*

## Completed
*(Empty)*

---

# Architecture Overview

## The Breakdown Design and Build Workflow
The `breakdown-design-and-build` agent definition establishes a four-phase lifecycle for building software using `breakdown`:
1. **Design Collaboration**: Architecting the feature and drafting the authoritative `design.md`.
2. **Breakdown & Audit**: Running `breakdown` (outputting directly to the plan folder) and manually pruning redundant tasks or fixing dependency ordering to match the design.
3. **The Execution Loop**: A rigid `Read-Analyze-Explain-Propose-HALT!` cycle that forces Test-Driven Development (RED/GREEN) and prevents runaway code generation.
4. **Finalization**: Updating global project documentation (patterns, deep dives) based on the completed feature.

---

# Checklist & TDD Requirements

- **Legend**:
  - `[DOCS]` - Documentation generation

- **TDD Enforcement**:
  - N/A for this specific task, as it is purely documentation generation.
