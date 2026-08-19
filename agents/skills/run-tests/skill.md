# Skill: run-tests

Purpose
- Run the project's test suite and report failures or flakiness in a machine- and human-readable form.

Depends on
- `run-in-devcontainer` — tests run in the `.devcontainer` image (or natively if already inside it). Do not require `go` on the host PATH.

Usage
- Inputs: none or optional `packages` (list of packages or `./...`).
- Output: structured result with `status` (ok/fail), `failed_packages`, `summary`, and `logs`.

Examples
- Run all tests:
  - `packages: ["./..."]`

Expected behavior
- Default command is `make test` from the repo root (on the host this wraps Docker; in the container it runs `go test ./...`).
- Optional `packages` maps to `go test <packages>` via `run-in-devcontainer`.
- Capture exit status, standard output, and standard error.
- Return a concise summary and identify the first failing package and stack trace.
