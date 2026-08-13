---
title: "Spec: one-way track for tabbed Google Docs"
status: accepted
audience: maintainers
date: 2026-08-13
supersedes-in-part: 2026-08-12-feature-request-tabbed-doc-directory-sync.md
references:
  - 2026-08-13-response-tabbed-doc-directory-sync.md
---

# Spec: one-way `track` for tabbed Google Docs

## Decision context

Bidirectional tabbed sync (one Google Doc ↔ a folder of Markdown files with tab-scoped import) is **not feasible** and is not scheduled. That decision is recorded in [2026-08-13-response-tabbed-doc-directory-sync.md](2026-08-13-response-tabbed-doc-directory-sync.md). This spec is the one-way alternative: **Docs → local only**.

`track` needs Get + convert + write. It does not add `tabId` to `batchUpdate`. `import` / `ApplyDocument` stay 1:1 and must refuse tracked trees so they cannot clobber tab 0.

## Commands

```bash
gdocs-markdown-sync track -doc <id> -out ./minutes/standing
gdocs-markdown-sync track -out ./minutes/standing
gdocs-markdown-sync track -out ./minutes/standing --force
gdocs-markdown-sync track -doc <id> -out ./minutes/standing -dry-run
```

- First run requires `-doc` and `-out`.
- Later runs may omit `-doc`; `_track.toml` supplies it. If both are given and they disagree, fail (even with `--force`).
- `-dry-run` lists writes/renames/chmods and does not touch files.

`export` remains one Markdown file. After fetch, if the document is **tabbed**, `export` fails with a non-zero exit and tells the user to run `track` (suggested argv includes the same `-doc`). It must not write first-tab-only Markdown.

A document is **tabbed** when the tab tree has more than one tab at any depth, or any tab has `childTabs`. A legacy body or a single leaf tab remains exportable.

`preview` uses the same fail-closed rule. `import` stays 1:1 and refuses tracked paths.

## Local layout

The Docs tab tree is mirrored on disk. A **leaf** tab is a Markdown file. A tab **with child tabs** is a folder; the parent’s own body lives in a file **named the same as that folder**.

```text
minutes/standing/
  _track.toml
  welcome.md
  fy2024/
    fy2024.md
    jan.md
    q1/
      q1.md
      march.md
```

- Walk `document.tabs` and `childTabs`. Frontmatter: `title`, `tab_id`, `doc_id`, `track: true`.
- Parent with children: `<slug>/` plus `<slug>/<slug>.md` (write even if the body is empty).
- Identity is `tab_id`. Tab rename in Docs → rename file/folder. Leaf gains children → promote `slug.md` → `slug/slug.md`. Parent loses all children → flatten back. Moves follow dirty/`--force` rules.
- New remote tabs → new files. Removed remote tabs → keep the local file, warn, do not delete.
- If a child slug would occupy the parent-body path (`slug/slug.md`), disambiguate the **child**, never the parent-body file.
- `_track.toml` lives at the **track root only**: `mode = "track"`, `doc_id`, `title`, ordered `tabs` (`tab_id`, relative `path`, `title`, `sha256` of last written bytes). Directories `0755`. Markdown files `0444`.

## Pull policy

Google Docs is the source of truth.

1. After a successful write, chmod each `.md` to `0444`. The track root, nested folders, and `_track.toml` stay writable.
2. **`0444` is a human deterrent, not a lock on `track`.** A later `track` that is allowed to update a file must chmod owner-writable, write, then chmod `0444` again. The same cycle applies to `--force` and to promote/flatten (the tool may rename/remove a `0444` file; restore `0444` on the new path). Do not treat EACCES as “skip this file.”
3. Dirty = local bytes hash ≠ sidecar `sha256` from the last track (detect before chmod-write). Unchanged files can be left `0444` without rewriting.
4. Any dirty file → refuse the whole run unless `--force`.
5. Conversion drift is not dirty: dirty means the file bytes changed since this tool last wrote them.

## First-run guard

`-out` is a **track root** only if it contains `_track.toml` with `mode = "track"`.

| `-out` state | Without `--force` | With `--force` |
| --- | --- | --- |
| Missing | Create dir, init sidecar, write tree, chmod `.md` `0444` | Same |
| Exists, empty (or only empty dirs) | Same as missing | Same |
| Exists with `_track.toml` (valid track root) | Refresh (dirty rules) | Refresh, overwrite dirty |
| Exists, non-empty, **no** `_track.toml` | **Refuse** (`ErrNotTrackRoot`) | Adopt: write/overwrite, create sidecar, chmod |
| Exists, `_track.toml` but `doc_id` ≠ `-doc` | **Refuse** | Still refuse; pick another `-out` |

Refuse if `-out` is an existing file (must be a directory). Only the `-out` directory itself is the track root, not a parent’s sidecar.

## Block push

Refuse `import` when:

- `-file` is inside a directory that contains `_track.toml` with `mode = "track"`, or
- the file’s frontmatter has `track: true`, or
- `-doc` matches a `_track.toml` `doc_id` found by walking parents of `-file`.

Error text must say to use `track` for refresh, not `import`.

## Non-goals

- Tab-scoped `batchUpdate` / round-trip import (declined).
- Creating or deleting remote tabs.
- Drive folder recursion.
- Comments / suggestions / layout fidelity.
