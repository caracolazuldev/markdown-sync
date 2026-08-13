# Mapping Rules — gdocs-markdown-sync

Scope
-----
This document defines the initial deterministic mapping rules between Google Docs document structure and Markdown.

Supported elements
------------------
- Headings → `#` / `##` / `###` based on heading level
- Paragraphs → plain paragraphs
- Code blocks → fenced code blocks with language where available
- Images → saved to an `assets/` folder; referenced by relative paths in Markdown
- Tables → converted to GitHub-flavored Markdown tables where possible; complex tables flagged
- Links & footnotes → inline links and footnote-style where detected
- Frontmatter (1:1 export) → YAML containing `title` (and `doc_id` / `last_modified` when available)
- Frontmatter (`track`) → YAML containing `title`, `tab_id`, `doc_id`, `track: true`

Tab tree → local paths (`track`)
--------------------------------
See APP-2026-08-13-tab-paths and `Missives/2026-08-13-spec-track-tabbed-docs.md`.

- Leaf tab → `<slug>.md` in the parent folder.
- Tab with children → folder `<slug>/` plus parent body `<slug>/<slug>.md` (written even if empty).
- Recurse on `childTabs`. Identity is `tab_id`, not path.
- `_track.toml` at the track root records `mode`, `doc_id`, `title`, and per-tab `path` + `sha256`.

Unsupported / Limitations
-------------------------
- Complex layout elements (textboxes, positioned images) are not preserved.
- Comments and suggestions are not mapped in v0.1 (future work).
- `export` does not flatten tabbed documents; use `track`.
- `track` does not push Markdown back to Docs (see ADR-2026-08-13-track).

Round-trip / determinism
------------------------
1:1 mapping is designed to be idempotent: repeated export of the same single-tab document yields the same Markdown (modulo non-deterministic metadata like timestamps). `track` refreshes from Docs; local edits are dirty unless `--force`.
