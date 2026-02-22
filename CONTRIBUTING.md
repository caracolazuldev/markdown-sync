## Contributing

Thank you for contributing to Markdown Sync. This guide covers repository layout, technology stack, development environment setup, and the expected developer workflow.

Repository layout
- `cmd/` — CLI entrypoint (packages for the `markdown-sync` binary).
- `internal/cli` — CLI helpers and argument parsing.
- `internal/google` — Google auth helpers and integrations.
- `internal/markdown` — Markdown conversion utilities.
- `internal/sync` — Sync logic between Google Docs and markdown files.
- `docs/` — Design and operational documentation.
- `.devcontainer/` — Development container configuration (Dockerfile + devcontainer.json).
- `Makefile` — Convenience tasks (setup, test, lint).

Technology stack
- Language: Go (module mode). Project targets Go `1.24` in the devcontainer.
- Primary libs: `golang.org/x/oauth2`, Google API client libraries, `cloud.google.com/go/*`.
- Tooling: `make`, `staticcheck`, `go test`, and VS Code Go extension.

Developer environment setup
1. Rebuild the devcontainer after pulling changes to use the pinned Go toolchain and tools.
   - In VS Code: Command Palette -> Remote-Containers: Rebuild Container.
2. Local setup (optional): install Go 1.24+ and run `make setup`.
3. Useful commands:
   - `make setup` — install tooling (if present in Makefile).
   - `make test` — run `go test ./...`.
   - `staticcheck ./...` — static analysis.
   - `go mod tidy` — tidy modules.

Developer workflow
- Branching: create a feature branch per change: `feature/`, `fix/`, `chore/`.
- Commits: keep commits small and focused; use descriptive messages.
- PRs: push a feature branch and open a Pull Request against `main`; include description and testing notes.
- Reviews: wait for CI success and at least one approving review before merging.

Git & automation policy (strict)
- Automated tools, bots, or assistants must:
	- Not commit or push directly to `main` (this repository assumes protected branch policies).
  - Create and operate on a feature branch.
  - Not push or merge into `main` without explicit human approval.
  - Include `go.mod` and `go.sum` changes in the feature branch when adding dependencies.

Adding dependencies
- On a feature branch, run `go get <pkg>@<version>` or `go install`, then `go mod tidy` and include `go.mod` and `go.sum` in the PR.

Code quality & testing
- Run `staticcheck` and add unit tests for new logic.
- `make test` runs all tests; CI will also run tests on PRs.

Contact
- If you're unsure about branching, tests, or releasing, open an issue or ask on the PR.
