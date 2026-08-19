# Testing — gdocs-markdown-sync

Unit tests
----------
- Keep unit tests fast and deterministic. Mock Google adapters (`markdown.Fetcher`) and only test mapping logic, track layout, and local storage behavior. Do not call the network.
- Tab walk: legacy body; one leaf; two top-level; nested `childTabs`.
- Path mapping: leaf vs parent folder; nested; promote; flatten; child slug collision; empty parent body still writes `slug/slug.md`.
- Export: single-tab still writes a file; tabbed returns an error containing `track` and does not write `-out`. Preview uses the same fail-closed rule.
- Track write: nested tree, `_track.toml` hashes, `.md` mode `0444`, dirs writable.
- Track refresh: chmod-write-chmod on `0444` files; unchanged files skipped; new tab adds file; remote rename; dirty without `--force` fails; `--force` overwrites and re-chmods.
- First-run guard: empty dir ok; non-empty untracked dir refuses; `--force` adopts; `doc_id` mismatch refuses even with `--force`; `-out` file path refuses.
- Import: file in tracked tree refused; `track: true` frontmatter refused; untracked file still allowed (diff-only in unit tests).
- Dry-run: no writes, no chmod.

Integration tests
-----------------
- Integration tests exercise Google Docs API and require credentials. They must be gated behind CI secrets and marked `//go:build integration` or run via `go test -tags=integration`.
- Provide instructions for creating test documents and service accounts in `docs/INTEGRATION.md`.
- Optional (not required for `track` unit coverage): a multi-tab fixture Doc via `INTEGRATION_TABBED_DOC_ID`.

CI
--
- CI runs unit tests for every PR. Integration tests run on a separate pipeline with secrets provided.

Local dev
---------
- Unit tests: `make test` inside the container. From the host: `make docker-run CMD='make test'`. Agents must not install Go on the host; see `AGENTS.md`.
- For local integration testing, set `GOOGLE_APPLICATION_CREDENTIALS` to a service account JSON and run `go test -tags=integration ./...` **inside the container** (or `make docker-run CMD='go test -tags=integration ./...'` on the host).
