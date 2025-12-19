# Architecture & Design Principles 🔧

This document captures the high-level architectural decisions and design principles for `gosak` — a small, pragmatic DevOps CLI toolkit written in Go.

---

## Project Structure (what you'll find)
- `cmd/` — Cobra-based command implementations (one file / command). Keep commands small and focused.
- `gurl/` — a packaged HTTP client used by `gurl` and other commands.
- `utils/` — shared helpers: constants, version info, error handling.
- `build.mk`, `common.mk` — build and release helpers (Makefile-driven workflows).
- `test/` — static samples and fixtures.

---

## Core Design Principles ✅
- **Single responsibility:** Each command performs one well-defined task (e.g., `assume` only handles STS assume and outputs env vars).
- **Composability:** Commands are script-friendly; output should be easy to parse and re-use in shell pipelines.
- **Small surface area:** Prefer specialized commands rather than huge monolithic subcommands.
- **Explicit errors and exit codes:** Use non-zero exits for failure conditions and clear error messages (see `utils` wrappers).
- **Stable CLI contract:** Flags and positional args should remain stable across releases to avoid breaking automation.
- **Minimal, well-scoped dependencies:** Avoid heavy dependencies unless necessary; prefer stdlib and small, focused libraries.
- **Observability-friendly:** Keep outputs machine-parseable where possible and add human-friendly details only when needed.
- **Testability:** Logic should be separated from CLI plumbing so it can be unit-tested.

---

## Recommended Patterns & Conventions 🔁
- **Command layout:** Use `cobra-cli` to scaffold commands; keep business logic in package functions (not only `Run` closures).
- **I/O injection:** Commands should accept io.Reader/io.Writer or accept config structs referencing outputs (see `gurl.Config`). This makes testing easier.
- **Config & env:** Avoid silent global state; favor explicit config structs and environment fallbacks.
- **Secrets:** Never log secret values; areas that call AWS or other auth should avoid printing credentials.
- **Duration & timeouts:** Make sensible defaults (e.g., STS duration default 3600s) and validate bounds.

---

## Build, Release & CI ⚙️
- **Semver tags:** Use semantic version tags (`vX.Y.Z`) for releases.
- **Cross-platform builds:** Use `make` rules to build artifacts for target OS/ARCH; ensure the build matrix in CI covers Linux/Mac/Windows.
- **CD flow:** On tag, build artifacts, create release artifacts, and publish (release notes via changelog).
- **Static analysis:** Run `go vet`, `golangci-lint` or equivalent in CI.
- **Unit tests:** Add unit tests for pure logic; add compat tests in integration where feasible.

---

## Extending the Project (how to add a new tool)
1. Use `cobra-cli add <name>` to scaffold the new command in `cmd/`.
2. Keep CLI parsing isolated; put the primary logic into a package function so it can be unit-tested and reused.
3. Add docs in `COMMANDS.md` and a short example in `README.md`.
4. Add tests and CI verification before merging.

---

## Security & Licensing 🔐
- Keep third-party dependencies up-to-date and minimize attack surface.
- The repository is licensed under the project `LICENSE`; verify compatibility when adding dependencies.

---

## Quick Guideline: When in doubt
- Prefer readable, testable code over clever shortcuts.
- Keep CLIs script-friendly and deterministic.
- Make changes in small, reviewable commits and document assumptions in PR descriptions.
