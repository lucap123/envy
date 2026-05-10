# envy

Zero-config environment variable manager for developers. One tool. Every project. No cloud. No bullshit.

## Why envy?

Managing `.env` files across multiple projects is a security risk and a productivity killer. Copying tokens, accidentally committing secrets, and losing track of required variables are problems every developer faces.

**envy** is a professional-grade, hardware-locked CLI tool that manages your environment variables locally and securely.

- **Zero config:** Automatically detects your project type (Node, Go, Python, Rust, PHP, Ruby).
- **Hardened Security:** Uses Argon2id key derivation and AES-256-GCM encryption.
- **Hardware-Locked:** Your encryption keys are tied to your specific machine using OS Keychains.
- **Memory Protected:** Secrets are locked in RAM and wiped after use to prevent memory inspection.
- **Git Safety:** Built-in hook to block commits if secrets are detected in your code.
- **No Cloud:** No account, no server, and no SaaS dependencies.

---

## Security Architecture

envy is built with a "security-first" architecture to ensure your secrets remain private even if your machine is compromised:

1. **Argon2id Passphrase:** Your master passphrase is used with the Argon2id algorithm to derive your encryption key on the fly. It is never stored on disk.
2. **OS Keychain Integration:** A unique device secret is stored in your system secure vault (Windows Credential Manager / macOS Keychain). This secret is required to derive your key, making your data machine-locked.
3. **RAM Isolation:** Sensitive data is managed with memory-locking, preventing it from being swapped to disk and ensuring it is zeroed out as soon as a command finishes.
4. **Zero-Knowledge Sharing:** Team bundles use independent encryption and per-bundle salts, keeping your main vault isolated and safe.

---

## Quick Start

### 1. Build and Initialize
Clone the repository and build the binary. On your first run, you will be prompted to set a Master Passphrase.
```bash
# Build the binary
make build

# Initialize a project
./envy init
```

### 2. Set variables
Store variables securely in your local encrypted store.
```bash
./envy set DATABASE_URL "postgres://localhost:5432/mydb"
```

### 3. Run your app
Inject your variables directly into any command without exposing them to your shell history.
```bash
./envy run npm start
```

---

## Command Reference

- `envy init` : Detect project type and create a `.env.example` template.
- `envy set <key> <value>` : Store a variable in the active profile.
- `envy get <key>` : Retrieve a variable (prints plain value).
- `envy list` : List all variables in the active profile (masked by default).
- `envy run <cmd>` : Run a command with all variables injected into the process.
- `envy export` : Generate a plaintext `.env` file from the active profile.
- `envy hook install` : Install a git pre-commit hook to prevent secret leaks.
- `envy profile add/use` : Manage multiple environments like development or production.
- `envy share/import` : Securely share encrypted bundles with your team.
- `envy logout` : Clear the memory-cached session and force passphrase re-entry.

---

## Installation

### Global Installation (Recommended)

If you have Go installed, you can install `envy` globally with a single command:

```bash
go install github.com/lucap/envy@latest
```

*Note: Ensure your `$GOPATH/bin` (or `%USERPROFILE%\go\bin` on Windows) is in your system PATH.*

### Build from Source

If you prefer to build locally:

```bash
git clone https://github.com/lucap/envy
cd envy
make build
```

## License

This project is licensed under the MIT License.
