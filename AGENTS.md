# Agent instructions — gdocs-markdown-sync

Supported toolchain is the image built from [`.devcontainer/Dockerfile`](.devcontainer/Dockerfile). Docker Engine or Docker Desktop is required on the host. A host Go install is not.

## Do not install Go on the host

If `go` is missing, that is expected. Do **not** install Go or `staticcheck` on the host: no `brew`, `apt`, `choco`, `winget`, `go install`, or similar.

## How to build, test, and lint

**Already in the devcontainer** (`/.dockerenv` exists, or `DEVCONTAINER=true`, or cwd is `/workspace` with Go 1.24): run Make natively.

```sh
make test
make lint
make build
```

**On the host** (any OS): `make test` / `make lint` / `make build` fail on purpose. Use Docker:

```sh
make docker-run CMD='make test'
make docker-run CMD='make lint'
make docker-run CMD='make build'
```

If `make` is unavailable (typical Windows PowerShell):

```powershell
docker build -t gdocs-markdown-sync:dev -f .devcontainer/Dockerfile .devcontainer
docker run --rm -v "${PWD}:/workspace" -w /workspace -e GOCACHE=/tmp/go-cache -e GOMODCACHE=/tmp/gomod -e DEVCONTAINER=true gdocs-markdown-sync:dev make test
```

Do not attempt docker-in-docker from inside the attached container.

## Catalog

Detailed skills, workflows, and policies: [`agents/index.yaml`](agents/index.yaml). Cursor-native skill: [`.cursor/skills/run-in-devcontainer/SKILL.md`](.cursor/skills/run-in-devcontainer/SKILL.md).
