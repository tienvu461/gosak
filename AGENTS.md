# Agents & Roles (who does what) 🤖👩‍💻

This document defines **human** and **automated** agents that interact with the `gosak` project, and the expectations for each agent.

---

## Human Roles

### Maintainer / Owner
- Approves PRs and manages releases.
- Responsible for semantic version tagging and release notes.
- Ensures CI is green before merging and publishing.

### Contributor / Developer
- Opens PRs with focused changes and tests.
- Adds or updates `COMMANDS.md` and `ARCHITECT.md` when the design or CLI surface changes.
- Fixes bugs and implements new tools following established patterns.

### Reviewer
- Reviews for correctness, test coverage, and potential security impacts.
- Ensures changes follow the architecture & style guidelines.

---

## Automated Agents

### CI Agent (Build / Test)
- Runs on PRs and the default branch.
- Jobs: `build`, `unit-tests`, `lint`, `vet`, `cross-build` (optional), `compat-checks` (OS matrix).
- Fails PRs with actionable logs. Should be required for merging.

### Release Agent
- Triggered by semver tag push (e.g., `v1.2.3`).
- Tasks: build artifacts for target platforms, attach binaries to GitHub Release, publish checksums, and optionally upload digest/signatures.
- Produces release notes based on changelog or PR titles.

### Linter & Static Analyzer (Auto-review)
- Runs linters (`golangci-lint`) and comments on PRs if configured.
- Suggests fixes, but does not block unless the maintainer enforces it.

### Security Scanner
- Runs dependency scanning (e.g., `go list -m -u all` and vulnerability checks) and reports known CVEs on PRs.

### Packaging / Distribution Agent
- Builds distribution archives (tar/zip) and optionally publishes to a release or artifact store.

---

## Agent Interaction Rules (best practices)
- **Fail fast:** CI should surface failures clearly to aid triage.
- **Immutable artifacts:** Release artifacts should be immutable and traceable to a tag and build hash.
- **Least privilege:** Release agents and any tokens must have minimal permissions required for their job.
- **Human approval gates:** Critical operations like publishing a release may require human approval (either via protected branches or manual workflow step).
- **Automation transparency:** CI and release logs must be retained and discoverable on the Release/PR.

---

## Example GitHub Actions Jobs (suggested)
- `lint`: run linters
- `test`: run `go test ./...`
- `build`: `go build` for the current platform
- `cross-build`: build artifacts for `linux, darwin, windows` (matrix)
- `release`: runs on tag, produces artifacts and publishes release

---

## Summary
Automate routine checks, keep humans in the loop for high-trust actions, and ensure reproducible builds and clear traceability from source to artifacts.
