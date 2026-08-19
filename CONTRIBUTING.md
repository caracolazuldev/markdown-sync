## Contributing

Thank you for contributing to Markdown Sync. This guide covers repository layout, technology stack, development environment setup, and the expected developer workflow.

Repository layout
- `cmd/` — CLI entrypoint (packages for the `gdocs-markdown-sync` binary).
- `internal/cli` — CLI helpers and argument parsing.
- `internal/google` — Google auth helpers and integrations.
- `internal/markdown` — Markdown conversion utilities.
- `internal/sync` — Sync logic between Google Docs and markdown files.
- `docs/` — Design and operational documentation.
- `.devcontainer/` — Development container configuration (Dockerfile + devcontainer.json). This is the supported toolchain.
- `Makefile` — Convenience tasks (setup, test, lint). Go targets require the container; from the host use `make docker-run CMD='make test'`.
- `AGENTS.md` — Instructions for coding agents (container-only toolchain).

Technology stack
- Language: Go (module mode). Project targets Go `1.24` in the devcontainer.
- Primary libs: `golang.org/x/oauth2`, Google API client libraries, `cloud.google.com/go/*`.
- Tooling: Docker, `make`, `staticcheck`, `go test`, and VS Code Go extension (inside the container).

Developer environment setup
1. Install Docker Engine or Docker Desktop. Rebuild the devcontainer after pulling changes to use the pinned Go toolchain and tools.
   - In VS Code / Cursor: Command Palette → Dev Containers: Rebuild Container.
2. Useful commands (inside the attached container they run natively; on the host Go targets error and tell you to use Docker):
   - `make test` — run `go test ./...` (container only).
   - `make lint` — run `staticcheck ./...` (container only).
   - `make build` — build `bin/gdocs-markdown-sync` (container only).
   - `make tidy` — `go mod tidy` (container only).
   - `make docker-run CMD='make test'` — run a Make target in the image from the host.
3. Host Go (human-only, optional): you may install Go 1.24+ yourself. **Agents must not** install Go on the host; see `AGENTS.md`.

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
  - Use `make docker-run CMD='make test'` (and `make lint`) from the host, or run Make inside the attached container; never install Go on the host.

Adding dependencies
- On a feature branch, run `go get <pkg>@<version>` then `make tidy` inside the container (from the host: `make docker-run CMD='go get ...'` then `make docker-run CMD='make tidy'`) and include `go.mod` and `go.sum` in the PR.

Code quality & testing
- Run `make lint` (or `make docker-run CMD='make lint'` from the host) and add unit tests for new logic.
- `make test` runs all unit tests inside the container; CI will also run tests on PRs.

Contact
- If you're unsure about branching, tests, or releasing, open an issue or ask on the PR.
