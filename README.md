# codegen

CLI tool that connects you to OpenRouter LLMs with file system access and shell execution.

## Quick start

```bash
cp .env.example .env   # add your OPENROUTER_KEY
go run .
```

Type messages normally. The LLM can read/write files and run commands when needed.

## Configuration

| Env var | Required | Default | Purpose |
|---------|----------|---------|---------|
| `OPENROUTER_KEY` | yes | — | API key |
| `MODEL` | no | OpenRouter default | Model slug (e.g. `openai/gpt-4o-mini`) |
| `COLUMNS` | no | 60 | Terminal width for box rendering |

## Commands

| Input | Action |
|-------|--------|
| `exit` | Quit |
| `list` | Show all available models |
| `list free` | Show free models only |

## LLM tools

The LLM can invoke these tools by outputting a JSON tool call:

- **readFile** — read a file from disk
- **writeFile** — write content to a file
- **bash** — execute a shell command (dangerous commands blocked)
- **lsDir** — list directory contents
- **glob** — find files by glob pattern (e.g. `**/*.go`)

Tool calls use raw JSON:
```json
{"tool": "readFile", "arguments": {"filename": "main.go"}}
```

## UI

```
┌─ You ─────────────────────────────────────
│ check if main.go exists
└───────────────────────────────────────────
┌─ Assistant (openai/gpt-4o-mini) ──────────
Let me check that file...
└───────────────────────────────────────────
┌─ Tool: readFile (main.go) ────────────────
│ package main
│ ...
└───────────────────────────────────────────
```

## Test

```bash
go test ./...
```
