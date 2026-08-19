---
name: run-in-devcontainer
description: Run build, test, lint, make, go, and staticcheck in the .devcontainer Docker image. Use when go is missing, when installing go or staticcheck on the host would seem necessary, or when running Makefile targets.
---

# Run in the devcontainer image

Never install Go or staticcheck on the host. The toolchain is `.devcontainer/Dockerfile`.

## Detect environment

Already in the container if any of these is true:

- `/.dockerenv` exists
- `DEVCONTAINER=true`
- cwd is `/workspace` and `go version` reports 1.24

If in the container, run Make natively (`make test`, `make lint`, `make build`). Do not start Docker from inside the container.

## On the host

Prefer Make (it builds the image and bind-mounts the repo at `/workspace`):

```sh
make test
make lint
make build
make docker-run CMD='go test ./internal/markdown'
```

If Make is unavailable:

```sh
docker build -t gdocs-markdown-sync:dev -f .devcontainer/Dockerfile .devcontainer
docker run --rm -v "$PWD:/workspace" -w /workspace -e GOCACHE=/tmp/go-cache -e GOMODCACHE=/tmp/gomod gdocs-markdown-sync:dev go test ./...
```

On Windows PowerShell use `${PWD}` instead of `$PWD`.

## Forbidden

- `brew` / `apt` / `choco` / `winget` install of Go
- `go install` or downloading a Go toolchain onto the host
- Treating a missing `go` binary as a problem to fix on the host
