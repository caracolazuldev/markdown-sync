# CLI Specification — gdocs-markdown-sync

Overview
--------
This document specifies the CLI surface for `gdocs-markdown-sync`.

Top-level commands
------------------
- `export` — export a Google Doc to a single Markdown file (fails if the Doc is tabbed)
- `import` — import a Markdown file into a Google Doc (refuses tracked trees)
- `preview` — render a Google Doc as Markdown to stdout (fails if the Doc is tabbed)
- `list` — list accessible Google Docs and basic metadata
- `track` — one-way pull of a tabbed Doc into a nested Markdown folder

Common flags
------------
- `--auth` (oauth|service) — choose auth flow (default: oauth)
- `--credentials` — path to service-account JSON when using service auth
- `--out` — output file for `export`, directory for `track`
- `--doc` — Google Doc ID
- `--file` — local Markdown file for import
- `--dry-run` — compute changes but do not apply
- `--force` — `track`: overwrite dirty local files or adopt a non-empty untracked directory
- `--overwrite` — overwrite remote doc on import (default: prompt)
- `--verbosity` — debug|info|warn|error (default: info)

Examples
--------
Export a single-tab (or legacy) doc to a file:

```
gdocs-markdown-sync export --doc 1a2B3cdE --out ./notes.md
```

A tabbed Doc cannot be exported as one file. The CLI exits non-zero and names `track`:

```
gdocs-markdown-sync track --doc 1a2B3cdE --out ./minutes/standing
```

Refresh a track root (doc id from `_track.toml`):

```
gdocs-markdown-sync track --out ./minutes/standing
```

Adopt an existing non-empty directory that is not yet a track root:

```
gdocs-markdown-sync track --doc 1a2B3cdE --out ./minutes/standing --force
```

Preview a doc as markdown:

```
gdocs-markdown-sync preview --doc 1a2B3cdE
```

Import a file into an existing doc (service account):

```
gdocs-markdown-sync import --file ./notes/article.md --doc 1a2B3cdE --auth service --credentials ./sa.json --overwrite
```

Behavior notes
--------------
- Commands are manual: no background sync.
- **Tabbed** means more than one tab at any depth, or any nested `childTabs`. `export` and `preview` fail closed and suggest `gdocs-markdown-sync track -doc <id> -out <dir>`.
- `track` is Docs → local only. Markdown files are chmod `0444` after write; the tool itself chmod-writes-chmod on later refreshes. See `Missives/2026-08-13-spec-track-tabbed-docs.md`.
- `import` refuses a file under a `_track.toml` root, a file with `track: true` frontmatter, or a `-doc` that matches that sidecar.
- Where destructive changes are possible, `--dry-run` and `--force` are explicit.
- Errors return non-zero exit codes.
