# Delete the unused `prompts/generate_plan_name.md` and `prompts/execute_plan.md` files, and clean up all corresponding embed directives and variable references in `prompts/prompts.go`.

This task is responsible for cleaning up the prompt templates that are no longer needed after the removal of the interactive planning features. You must delete the files `prompts/generate_plan_name.md` and `prompts/execute_plan.md` from the repository.

After removing the files, you must update `prompts/prompts.go` to remove any `//go:embed` directives, variables, or constants that reference these deleted Markdown files. Ensure that the codebase compiles cleanly and only retains prompts required for the core non-interactive task decomposition logic (e.g., `analyze_task.md`).
