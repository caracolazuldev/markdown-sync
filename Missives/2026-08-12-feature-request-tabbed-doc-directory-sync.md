---
title: "Feature request: tabbed Google Doc ↔ directory of Markdown files"
status: review
audience: maintainers
date: 2026-08-12
from: co-op-strategic-planning (SPC workspace)
---

# Feature request: tabbed document ↔ Markdown directory

## Summary

Support syncing a **multi-tab Google Doc** as a **folder of Markdown files** (one file per document tab), with stable round-trip mapping via tab IDs.

## Motivation

In the Co-op Strategic Planning workspace, official committee minutes live as **one Google Doc with one tab per meeting**, plus separate Docs for special meetings. Stakeholders collaborate in Docs; Markdown in git is the system of record elsewhere.

Today `gdocs-markdown-sync` maps **one document ID ↔ one Markdown file** and does not use `includeTabsContent` or `document.tabs`. For multi-tab Docs:

- Reads effectively see the **first tab only** (API default when tabs content is not requested).
- Writes/`batchUpdate` default to the **first tab** unless a tab ID is specified.
- Syncing a multi-tab minutes Doc would be **lossy** and unsafe.

A folder-of-files model matches how teams already version meeting packs in git and how MKA Missives are organized.

## Proposed behavior

### Export (Docs → Markdown)

1. `documents.get` with `includeTabsContent=true`.
2. Walk `document.tabs` (define policy for **child tabs**: flatten with path, skip, or error — recommend **configurable**; default skip nested child tabs or flatten as `parent/child.md`).
3. Emit a directory, e.g.:

```text
out/<doc-slug>/
  _doc.toml          # doc_id, title, tab index
  <tab-slug>.md      # frontmatter: title, tab_id, doc_id
```

4. CLI sketch:

```bash
gdocs-markdown-sync export -doc <id> -out ./minutes/standing -tabs dir
```

### Import (Markdown → Docs)

1. Read directory + frontmatter `tab_id` (create tab if missing — optional later phase).
2. Apply each file’s body to the corresponding tab via tab-scoped `batchUpdate` requests.
3. CLI sketch:

```bash
gdocs-markdown-sync import -doc <id> -dir ./minutes/standing -tabs dir
```

### Registry / mapping

Extend mapping beyond single `path` ↔ `doc_id` to support:

- `doc_id` + `tab_id` ↔ file path, or
- `doc_id` ↔ directory path with per-file `tab_id` in frontmatter.

## Non-goals (v1)

- Drive **folder** recursion (separate feature; consumers may list Drive themselves).
- Comments / suggestions sync.
- Perfect fidelity for complex layout.

## Acceptance criteria

1. Export of a Doc with N top-level tabs yields N Markdown files with correct `tab_id` metadata.
2. Import updates the matching tabs without clobbering other tabs’ content.
3. Single-tab Docs keep working with the existing one-file CLI.
4. Documented limitation: Docs without tabs / legacy body field behavior.
5. Unit tests for tab tree walk; gated integration test against a multi-tab fixture Doc.

## Workarounds today

- Use **one Doc per meeting** (single tab) for anything that must sync.
- Do **not** register multi-tab minutes Docs for push/pull with the current CLI.
- Optional one-off pull scripts outside this tool (not a substitute for round-trip).

## References

- [Google Docs API — Work with tabs](https://developers.google.com/docs/api/how-tos/tabs)
- Consumer context: `/home/mzd/Lab/co-op-strategic-planning` (SPC); minutes sync deferred until this lands.
