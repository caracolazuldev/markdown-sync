# Decisions (Architecture & Technical)

This file records architecture decisions and business-logic rationales. Each entry should follow the Decision record format described in `docs/architecture-policy.md`.

- id: ADR-2026-08-18-devcontainer-toolchain
- title: Devcontainer image is the only supported agent toolchain; do not install Go on the host
- date: 2026-08-18
- author: repo-maintainer
- decision: Coding agents (any host OS) must build, test, and lint using the image from `.devcontainer/Dockerfile`. They must not install Go or staticcheck on the host. Makefile Go targets (`test`, `lint`, `build`, …) require the container and error on the host with `make docker-run CMD='make <target>'`. Cross-tool instructions live in `AGENTS.md`; Cursor also has `.cursor/rules/dev-environment.mdc` and `.cursor/skills/run-in-devcontainer`.
- rationale: Agents were installing host Go because `agents/` is not auto-loaded by Cursor/Copilot/Claude Code and `run-tests` required `go` on PATH. Implicit Make-to-Docker wrapping was rejected as hard to read; a hard error pointing at `docker-run` is explicit. The Dockerfile already pins Go 1.24 and tools.
- alternatives: Optional host Go for agents (rejected; diverges from the pinned image). Dev Containers CLI as the only runner (rejected; not installed on most agent hosts). Makefile silently wrapping Docker for `make test` (rejected; implicit and surprising).
- impact: Makefile parse-time container guard plus `docker-image`/`docker-run`; `AGENTS.md`; Cursor rule and skill; `agents/skills/run-in-devcontainer`; rewritten `run-tests`; CONTRIBUTING, README, docs/testing.md. CI continues to use `actions/setup-go` and does not call `make test`.
- reassessment_triggers:
  - Dev Containers CLI is universally available to agents
  - CI switched to the same Dockerfile
  - Docker is unavailable on a required agent host
- links:
  - Dockerfile: `.devcontainer/Dockerfile`
  - Agent instructions: `AGENTS.md`

- id: ADR-2026-08-13-track
- title: One-way `track` for tabbed Docs; `export` stays 1:1 and fails closed
- date: 2026-08-13
- author: repo-maintainer
- decision: Do not implement bidirectional tabbed sync. Add `track` (Docs → nested local folder, pull-only). `export` and `preview` remain one file and fail if the document is tabbed (more than one tab at any depth, or any `childTabs`), directing the user to `track`. `import` refuses tracked trees. Round-trip tab `batchUpdate` stays out of scope.
- rationale: Tab-scoped apply is a second product on an unfinished 1:1 adapter (`FetchDocument` was a stub; `ApplyDocument` rewrites the legacy body). One-way Get + convert + write is feasible and covers the minutes-in-git use case. Failing `export` on tabbed Docs avoids silent first-tab-only loss.
- alternatives: Bidirectional directory sync (declined; see Missives/2026-08-13-response-tabbed-doc-directory-sync.md). Concatenate tabs into one file (poor diffs). Auto-detect tabs and change `export -out` into a directory (breaks 1:1 callers). Skip nested tabs (rejected; nested folders are required).
- impact: CLI `track`; injectable `Fetcher`; `internal/sync` orchestration; `_track.toml`; chmod `0444` with tool-owned chmod-write-chmod; first-run guard for untracked dirs; unit tests offline; docs/cli-spec, mapping, architecture, testing, INTEGRATION, MANIFEST.
- reassessment_triggers:
  - Need to create remote tabs from new local files
  - Docs API tab-create/delete becomes a supported, tested path
  - Consumers require comments or Drive-folder recursion
- links:
  - Spec: Missives/2026-08-13-spec-track-tabbed-docs.md
  - Decline (do not amend): Missives/2026-08-13-response-tabbed-doc-directory-sync.md
  - Approach: APP-2026-08-13-tab-paths

Example decision entry (format)

- id: ADR-2026-02-23-oauth2
- title: Upgrade `golang.org/x/oauth2` to v0.35.0 and require Go 1.24 in devcontainer
- date: 2026-02-23
- author: repo-maintainer
- decision: Upgrade oauth2 to obtain improved OAuth features; require Go 1.24 in the devcontainer and update `go.mod`.
- rationale: v0.35.0 contains bug fixes and features needed for service flows; requires Go 1.24 which is acceptable for our devcontainer.
- alternatives: pin older oauth2 and avoid toolchain bump (rejected due to missing features).
- impact: update `.devcontainer/Dockerfile`, `go.mod`, and tests; update README and MANIFEST.
- reassessment_triggers:
  - security advisory on oauth2
  - breaking change in oauth2 or Go toolchain
  - CI/test flakiness introduced by the upgrade
- links:
  - PR: TODO
