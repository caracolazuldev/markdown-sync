# Dependencies — project notes

This file captures project-specific notes about important dependencies and compatibility considerations.

- `golang.org/x/oauth2` — pinned to `v0.35.0` for enhanced OAuth features; requires Go 1.24+. See `go.mod`.
- `cloud.google.com/go/compute/metadata` — indirect dependency used by oauth2/google; keep updated via `go mod tidy`.

When upgrading a dependency:
1. Run `make test` and `make lint` in a feature branch (inside the container, or `make docker-run CMD='make test'` from the host).
2. Document compatibility notes here and in the PR description.
3. If upgrade requires a Go toolchain bump, update `.devcontainer/Dockerfile` (the supported toolchain) and `README.md` / `AGENTS.md` with migration notes. Do not document a host Go install as the agent path.
