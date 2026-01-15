# Ralphmaster

CLI tool for structured LLM-GitHub communication via issues. Forces consistent formatting for task management between AI agents.

## Installation

```bash
brew install go
go mod tidy
go build -o ralphmaster
```

Requires [GitHub CLI](https://cli.github.com/) (`gh`) to be installed and authenticated.

## Commands

### Create Issue
```bash
ralphmaster new --title "Feature X" --description "Implement feature X"
```

### Add Metadata
```bash
ralphmaster metadata --issue 1 --branch "feature/x" --overview "Implement feature X"
ralphmaster metadata --issue 1 --branch "feature/x" --overview "Updated overview" --force
```

### Add Task
```bash
ralphmaster task add --issue 1 --model sonnet --goal "Create API endpoint"
ralphmaster task add --issue 1 --model opus --goal "Design architecture" --comments "Consider scalability" --references "src/api/main.go#10-20"
```

Models: `opus` (complex), `sonnet` (simple), `haiku` (trivial)

### List Undone Tasks
```bash
ralphmaster task list --issue 1
```

### Mark Task Done
```bash
ralphmaster task done --issue 1 --task 1 --commit "abc123" --work-done "Created REST endpoint with validation"
```

### Start Issue
```bash
ralphmaster start --issue 1
```
Adds `[IN PROGRESS]` to the first line of the issue description.

### Complete Issue
```bash
ralphmaster done --issue 1
```
Adds `[DONE]` to the first line of the issue description and closes the issue.

## Global Flags

All commands support `--repo owner/name` to override auto-detection from git remote.

## Comment Formats

**Metadata:**
```
[METADATA]
branch: feature/x
overview: Single line description
```

**Task (undone):**
```
[UNDONE]
id: 1
model: sonnet
goal: Create API endpoint
comments: Additional context here
references: src/api/main.go#10-20
```

**Task (done):**
```
[DONE]
id: 1
model: sonnet
goal: Create API endpoint
comments: Additional context here
references: src/api/main.go#10-20

work_done: Summary of completed work
commit: abc123
```
