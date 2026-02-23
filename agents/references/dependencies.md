# Dependencies — project notes

This file captures project-specific notes about important dependencies and compatibility considerations.

- `golang.org/x/oauth2` — pinned to `v0.35.0` for enhanced OAuth features; requires Go 1.24+. See `go.mod`.
- `cloud.google.com/go/compute/metadata` — indirect dependency used by oauth2/google; keep updated via `go mod tidy`.

When upgrading a dependency:
1. Run tests and static analysis in a feature branch.
2. Document compatibility notes here and in the PR description.
3. If upgrade requires a Go toolchain bump, update `.devcontainer/Dockerfile` and `README.md` with migration notes.
