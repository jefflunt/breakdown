# Remove Git Repository Requirement

## User Story
**Headline:** Remove the requirement to run `breakdown` strictly inside a Git repository.
**Problem Statement:** `breakdown` currently forces the user to run the tool from within a Git repository (checked via `git status` in `main.go`). This prevents users from using the tool in fresh, un-versioned directories or for general brainstorming tasks outside of a codebase.
**Objective:** Remove the `isGitRepo()` check and its associated failure branch from `cmd/breakdown/main.go`.
**Expected Outcome:** `breakdown` can be executed in any directory without exiting with an error about missing a Git repository.

## Implementation Backlog

### Pending

### Current

### Completed
- `remove-git-check.md`: Remove the `isGitRepo()` function and the `if !isGitRepo() { ... }` block from `cmd/breakdown/main.go`.

## Architecture Overview
The `main.go` entry point simply drops the guard clause that ensures `.git` presence. Codebase context relies on the file system tree builder which will just scan whatever directory it is executed in. If no `.gitignore` is present, it might scan more files, but that is acceptable for an un-versioned directory.

## Checklist & TDD Requirements
1. **Tests Pass:** Standard build and tests pass.
2. **Compile Check:** `go build ./...` succeeds.
3. **Behavior Verification:** Running the compiled binary in a non-git directory does not fail.