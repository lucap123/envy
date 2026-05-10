# envy

> Zero-config environment variable manager for developers.
> One tool. Every project. No cloud. No bullshit.

## Why envy?

Every developer knows the pain of managing `.env` files across multiple projects. Copy-pasting tokens, accidentally committing secrets, and forgetting which variables are needed for a project are common frustrations.

**envy** is different: **one command and it just works.**

- ✅ **Zero config:** Automatically detects your project type (Node, Go, Python, etc.)
- ✅ **Secure:** Stores your variables encrypted locally (~/.envy/)
- ✅ **Safe:** Built-in git hook to prevent accidental secret commits
- ✅ **No Cloud:** No account, no server, no SaaS dependency
- ✅ **Team Friendly:** Easily share encrypted bundles with your team

## Quick Start

### 1. Initialize
Detect your project type and create a template `.env.example`.
```bash
envy init
```

### 2. Set variables
Store variables securely in your local encrypted store.
```bash
envy set DATABASE_URL postgres://localhost:5432/mydb
```

### 3. Run your app
Inject your variables directly into any command.
```bash
envy run npm start
```

## Features

- `envy init` — Project auto-detection
- `envy set/get/list` — Manage encrypted variables
- `envy run` — Inject variables into processes
- `envy hook install` — Prevent secret leaks in git commits
- `envy profile add/use` — Manage multiple environments (dev, staging, prod)
- `envy share/import` — Offline encrypted team sharing

## Installation

```bash
# Coming soon!
go install github.com/lucap/envy@latest
```

## License

MIT
