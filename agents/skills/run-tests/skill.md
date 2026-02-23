# Skill: run-tests

Purpose
- Run the project's test suite and report failures or flakiness in a machine- and human-readable form.

Usage
- Inputs: none or optional `packages` (list of packages or `./...`).
- Output: structured result with `status` (ok/fail), `failed_packages`, `summary`, and `logs`.

Examples
- Run all tests:
  - `packages: ["./..."]`

Expected behavior
- The skill runs `make test` (or `go test` directly) in the repo root.
- It captures exit status, standard output, and standard error.
- It returns a concise summary and identifies the first failing package and stack trace.
