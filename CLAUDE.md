# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ralphmaster is a Go CLI tool for structured LLM-GitHub communication via GitHub issues. It enforces consistent formatting for task management between AI agents and development teams, wrapping the `gh` CLI.

## Build Commands

```bash
go mod tidy        # Install dependencies
go build -o ralphmaster    # Build binary
```

Requires GitHub CLI (`gh`) to be installed and authenticated.

## Architecture

The codebase uses a **pluggable backend architecture** with three layers:

1. **CLI Layer** (`cmd/`): Cobra command handlers that interact only with the backend interface
2. **Backend Interface** (`backend/backend.go`): Platform-agnostic abstraction defining issue/comment operations
3. **Backend Implementations** (`backend/github/`): Concrete implementations (currently GitHub via `gh` CLI)

The factory pattern in `backend/factory.go` enables future support for other platforms (Linear, Jira, etc.). Backends self-register via `init()` functions.

## Key Patterns

- **Backend initialization**: Commands call `backend.New(cfg)` to get a backend instance, never interact with GitHub directly
- **Error handling**: Print to stderr via `fmt.Fprintf(os.Stderr, ...)`, exit with code 1
- **Model validation**: Valid models are `opus`, `sonnet`, `haiku`
- **Comment parsing**: Uses regex to identify `[METADATA]`, `[UNDONE]`, `[DONE]` prefixed comments
- **Task IDs**: Auto-incremented based on existing task comments

## GitHub Backend Comment Formats

These formats are specific to the GitHub backend (`backend/github/`). Other backends (Linear, Jira, etc.) can be implemented with their own storage mechanisms.

Metadata: `[METADATA]\nbranch: ...\noverview: ...`

Task (undone): `[UNDONE]\nid: N\nmodel: ...\ngoal: ...\ncomments: ...\nreferences: ...`

Task (done): Same as undone but with `[DONE]` prefix, plus `work_done:` and `commit:` fields appended after a blank line.

## Global Flags

All commands support:
- `--repo owner/name`: Override GitHub repo auto-detection
- `--backend github`: Backend type (default: github)
