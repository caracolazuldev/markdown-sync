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

`make test` (and other Go targets) fail unless you are in the container. Run them via Docker:

```sh
make docker-run CMD='make test'
make docker-run CMD='make lint'
make docker-run CMD='make build'
make docker-run CMD='go test ./internal/markdown'
```

If Make is unavailable:

```sh
docker build -t gdocs-markdown-sync:dev -f .devcontainer/Dockerfile .devcontainer
docker run --rm -v "$PWD:/workspace" -w /workspace -e GOCACHE=/tmp/go-cache -e GOMODCACHE=/tmp/gomod -e DEVCONTAINER=true gdocs-markdown-sync:dev make test
```

On Windows PowerShell use `${PWD}` instead of `$PWD`.

## Forbidden

- `brew` / `apt` / `choco` / `winget` install of Go
- `go install` or downloading a Go toolchain onto the host
- Treating a missing `go` binary as a problem to fix on the host
- Running `make test` on the host and then installing Go to make it pass
