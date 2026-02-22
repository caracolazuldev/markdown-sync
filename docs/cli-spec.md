# CLI Specification — markdown-sync

Overview
--------
This document specifies the initial CLI surface for `markdown-sync` (v0.1).

Top-level commands
------------------
- `export` — export a Google Doc to Markdown files
- `import` — import a Markdown file into a Google Doc
- `preview` — render a Google Doc as Markdown to stdout
- `list` — list accessible Google Docs and basic metadata

Common flags
------------
- `--auth` (oauth|service) — choose auth flow (default: oauth)
- `--credentials` — path to service-account JSON when using service auth
- `--out` — output directory or file for export
- `--doc` — Google Doc ID
- `--file` — local Markdown file for import
- `--dry-run` — compute changes but do not apply
- `--overwrite` — overwrite remote doc on import (default: prompt)
- `--verbosity` — debug|info|warn|error (default: info)

Examples
--------
Export a doc to local folder:

markdown-sync export --doc 1a2B3cdE --out ./docs

Preview a doc as markdown:

markdown-sync preview --doc 1a2B3cdE

Import a file into an existing doc (service account):

markdown-sync import --file ./notes/article.md --doc 1a2B3cdE --auth service --credentials ./sa.json --overwrite

Behavior notes
--------------
- Commands are manual: no background sync.
- Where destructive changes are possible, `--dry-run` and interactive confirmation are available.
- Errors are non-fatal where possible and return non-zero exit codes on failure.
