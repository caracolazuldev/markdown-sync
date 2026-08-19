# Skill: run-in-devcontainer

Purpose
- Run any Go, Make, or staticcheck command in the image built from `.devcontainer/Dockerfile`.
- Prevent agents from installing a Go toolchain on the host.

Usage
- Inputs: `cmd` (string; the command to run inside `/workspace`).
- Output: the command's exit status, stdout, and stderr.

Detect environment
- In-container if `/.dockerenv` exists, `DEVCONTAINER=true`, or cwd is `/workspace` with Go 1.24.
- In-container: run `cmd` natively (e.g. `make test`). Do not start Docker from inside the container.
- On the host: Go Make targets error. Use `make docker-run CMD='make test'` (builds the image if needed). If Make is missing, `docker build -t gdocs-markdown-sync:dev -f .devcontainer/Dockerfile .devcontainer` then `docker run --rm -v <repo>:/workspace -w /workspace -e GOCACHE=/tmp/go-cache -e GOMODCACHE=/tmp/gomod -e DEVCONTAINER=true gdocs-markdown-sync:dev make test`.

Examples
- `cmd: make test`
- `cmd: make lint`
- `cmd: go test ./internal/markdown`

Forbidden
- Installing Go or staticcheck on the host (`brew`, `apt`, `choco`, `winget`, `go install` of a toolchain).
- Treating a missing host `go` binary, or a `make test` host error, as a reason to install Go locally.
