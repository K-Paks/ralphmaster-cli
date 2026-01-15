# Ralphmaster CLI - Implementation Plan

## Overview
CLI tool that enforces structured GitHub issue communication for LLM agents. Wraps `gh` CLI to ensure consistent formatting.

## Global Behavior
- **Repo detection**: Auto-detect from current git directory (like `gh`), with optional `--repo owner/name` flag on all commands

---

## Commands

### 1. `ralphmaster new`
Create a new GitHub issue.

**Flags:**
- `--title` (required): Issue title
- `--description` (required): Issue body in markdown
- `--repo` (optional): Override repo detection

**Implementation:**
```bash
gh issue create --title "<title>" --body "<description>"
```

**Files to modify:** `cmd/new.go`

---

### 2. `ralphmaster metadata`
Add/replace metadata comment on an issue.

**Flags:**
- `--issue` (required): Issue number
- `--branch` (required): Branch name for the work
- `--overview` (required): Single-line task overview
- `--force`: Replace existing metadata comment
- `--repo` (optional): Override repo detection

**Logic:**
1. Fetch all comments on the issue via `gh api`
2. If no comments exist → add metadata comment
3. If exactly 1 comment exists with `[METADATA]` on first line:
   - With `--force` → delete it and add new one
   - Without `--force` → do nothing (exit 0)
4. If other comments exist → return error

**Comment format:**
```
[METADATA]
branch: <branch>
overview: <overview>
```

**Files to create:** `cmd/metadata.go`

---

### 3. `ralphmaster task add`
Add a task comment to an issue. Task number is auto-incremented based on existing tasks.

**Flags:**
- `--issue` (required): Issue number
- `--model` (required): `opus|sonnet|haiku`
- `--goal` (required): Task description
- `--comments` (optional): Extra context
- `--references` (optional): File references (e.g., `path/to/file.ts#10-17`)
- `--repo` (optional): Override repo detection

**Logic:**
1. Fetch all comments on the issue
2. Find highest task number from comments starting with `[UNDONE]` or `[DONE]`
3. Auto-increment to next number
4. Post new task comment

**Comment format:**
```
[UNDONE]
id: <number>
model: <model>
goal: <goal>
comments: <comments>
references: <references>
```

**Files to create:** `cmd/task.go`, `cmd/task_add.go`

---

### 4. `ralphmaster task list`
List undone tasks on an issue.

**Flags:**
- `--issue` (required): Issue number
- `--repo` (optional): Override repo detection

**Logic:**
1. Fetch all comments on the issue
2. Filter to those starting with `[UNDONE]`
3. Display task number, model, and goal

**Files to create:** `cmd/task_list.go`

---

### 5. `ralphmaster task done`
Mark a task as completed.

**Flags:**
- `--issue` (required): Issue number
- `--task` (required): Task number to mark done
- `--commit` (required): Commit hash/reference
- `--work-done` (required): Summary of completed work
- `--repo` (optional): Override repo detection

**Logic:**
1. Find comment with matching task number that starts with `[UNDONE]`
2. Replace `[UNDONE]` with `[DONE]`
3. Append separator and work summary

**Updated comment format:**
```
[DONE]
id: <number>
model: <model>
goal: <goal>
comments: <comments>
references: <references>

work_done: <work-done summary>
commit: <commit>
```

**Files to create:** `cmd/task_done.go`

---

## Project Structure
```
ralph-cli/
├── go.mod
├── main.go
└── cmd/
    ├── root.go        # Root command + global --repo flag
    ├── new.go         # ralphmaster new
    ├── metadata.go    # ralphmaster metadata
    ├── task.go        # ralphmaster task (parent)
    ├── task_add.go    # ralphmaster task add
    ├── task_list.go   # ralphmaster task list
    └── task_done.go   # ralphmaster task done
```

---

## Implementation Steps

1. **Add global --repo flag to root command** - Persistent flag inherited by all subcommands

2. **Update `new` command** - Add `--title` and `--description` flags, execute `gh issue create`

3. **Create `metadata` command** - Implement comment fetching, validation logic, and metadata insertion

4. **Create `task` parent command** - Simple parent for task subcommands

5. **Create `task add` command** - Validate model enum, auto-increment task number, format and post comment

6. **Create `task list` command** - Fetch comments, filter `[UNDONE]`, display

7. **Create `task done` command** - Find task comment, update content, edit via API

---

## gh CLI Commands Used

```bash
# Create issue
gh issue create --title "..." --body "..."

# List comments (JSON for parsing)
gh api repos/{owner}/{repo}/issues/{number}/comments

# Add comment
gh issue comment {number} --body "..."

# Delete comment
gh api repos/{owner}/{repo}/issues/comments/{id} -X DELETE

# Edit comment
gh api repos/{owner}/{repo}/issues/comments/{id} -X PATCH -f body="..."
```

---

## Verification
1. Install Go: `brew install go`
2. Build: `go mod tidy && go build -o ralphmaster`
3. Test each command:
   - `./ralphmaster new --title "Test" --description "Test desc"`
   - `./ralphmaster metadata --issue 1 --branch "feature/x" --overview "Do X"`
   - `./ralphmaster task add --issue 1 --model sonnet --goal "Implement X"`
   - `./ralphmaster task list --issue 1`
   - `./ralphmaster task done --issue 1 --task 1 --commit "abc123" --work-done "Did X"`
