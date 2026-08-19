# gdocs-markdown-sync

Synchronize Markdown documents with Google Docs. Intended to facilitate collaboration between stakeholders and developers working in generative AI environments. Get the collaborative benefits of Google Docs and then export to version-controlled Markdown files.

## Development container

The supported toolchain is the image in [`.devcontainer/Dockerfile`](.devcontainer/Dockerfile) (Go 1.24). Docker Engine or Docker Desktop is required. Coding agents must use this image and must not install Go on the host; see [`AGENTS.md`](AGENTS.md).

Rebuild the container after pulling toolchain changes:

- VS Code / Cursor: Command Palette → Dev Containers: Rebuild Container.
- Host Make (wraps Docker): `make test`, `make lint`, `make build`.
