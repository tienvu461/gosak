# Commands Reference ✅

This file documents the public CLI surface exposed by `gosak` and gives short usage examples.

## Summary
- `gosak` is a collection of small, single-purpose utilities. Each command focuses on one job and aims to be script-friendly and composable.

---

## Commands

### `gurl <URL>` 🔧
- **Description:** HTTP client (lightweight curl-like).
- **Key Flags:**
  - `-H, --headers`  custom headers (string slice)
  - `-u, --user-agent` set User-Agent (default: `gurl`)
  - `-d, --data` HTTP body
  - `-m, --method` HTTP method (default: `GET`)
  - `-k, --insecure` skip TLS verification
- **Example:**
  - `gosak gurl -m POST -d '{"x":1}' -H "Content-Type: application/json" https://example/api`

---

### `assume` 🔐
- **Description:** Assume AWS IAM role via STS and emit environment variables for shell evaluation.
- **Key Flags:**
  - `-r, --role-arn` ARN of role to assume (either this or `--profile` must be provided)
  - `-p, --profile` profile name in `~/.aws/config` to assume
  - `-d, --duration` session duration in seconds (default: `3600`)
  - `-o, --output` print output (reserved)
- **Behavior:** Prints exports such as `AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`/`AWS_SESSION_TOKEN` and `ASSUMED_ROLE`.
- **Example:**
  - `eval $(gosak assume -r arn:aws:iam::123456789012:role/example -d 3600)`

---

### `ec2`
- **Description:** Manage AWS EC2 instances.
- **Subcommands:**
  - `list` List all EC2 instances.
  - `start <instance-id>|<instance-name>` Start an EC2 instance.
  - `stop <instance-id>|<instance-name>` Stop an EC2 instance.
- **Example:**
  - `gosak ec2 list`
  - `gosak ec2 start i-1234567890`
  - `gosak ec2 stop instance-name-1`

---

### `version` ℹ️
- **Description:** Print build/version metadata (`Version`, `Hash`, `OS`, `Arch`, `GoVersion`).
- **Usage:** `gosak version`

---
### `b64` 📦
- **Description:** Base64 encode/decode small strings.
- **Key Flags:**
  - `-e, --encode` encode input (default: decode)
- **Usage:** `gosak b64 <text>`
- **Example:**
  - `gosak b64 -e "hello"` -> prints base64 of "hello"
  - `gosak b64 "aGVsbG8="` -> decodes and prints `"hello"`

---

### `format` 🎨
- **Description:** Pretty-print/format JSON or dictionary-like inputs.
- **Key Flags:**
  - `--json` (default true) treat input as JSON
- **Usage:** `gosak format '<json-or-dict>'`

---

### `sslcert <domain>:443` 🔒
- **Description:** Fetch SSL certificate info (issuer, expiry, common name, TLS connection).
- **Usage:** `gosak sslcert example.com:443`

---

### `ifconfig` 🌐
- **Description:** Fetch current public IP (uses internal `gurl` to call `https://ifconfig.me`).
- **Usage:** `gosak ifconfig`

---

### `get` (legacy) ⛔
- **Description:** Generic HTTP GET helper (simple proof-of-concept). Prefer `gurl` for advanced use.
- **Usage:** `gosak get [url]` (default: `https://ifconfig.me`)

---

### `version` ℹ️
- **Description:** Print build/version metadata (`Version`, `Hash`, `OS`, `Arch`, `GoVersion`).
- **Usage:** `gosak version`

---

## Notes & Best Practices
- Commands are intentionally small and script-friendly: they print machine-friendly outputs and exit with non-zero codes on errors.
- Prefer `gurl` over `get` for HTTP needs; `get` remains for backward compatibility.
- To add commands, use `cobra-cli add <name>` and keep each command in `cmd/` with a focused responsibility.
- All added commands need to be documented
