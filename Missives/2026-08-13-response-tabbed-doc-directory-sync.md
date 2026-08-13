---
title: "Response: tabbed Google Doc ↔ directory of Markdown files"
status: decided
audience: co-op-strategic-planning (SPC workspace); maintainers
date: 2026-08-13
in-reply-to: 2026-08-12-feature-request-tabbed-doc-directory-sync.md
---

# Response: tabbed document ↔ Markdown directory

## Decision

**This request is not feasible to implement** in `gdocs-markdown-sync` as it stands, and it will not be scheduled.

The tool’s contract is one Google Doc ID ↔ one Markdown file. Tab-aware round-trip is not a flag on that contract. It requires a second product surface (directory of files), a second identity model (`doc_id` + `tab_id`), and a rewrite of the Google adapter (`includeTabsContent`, walk `document.tabs`, `tabId` on every `Range`/`Location`). `FetchDocument` is still a stub; live `ApplyDocument` already rewrites the legacy body with no tab scope. Building tabbed directory sync on that base would be unsafe and would not land as a thin increment.

SPC should keep the existing workaround: **one Doc per meeting** (single tab) for anything that must sync, and **do not register multi-tab minutes Docs** for push/pull.

The analysis below is retained so the decision is reviewable and so a later major version, if any, does not restart from a blank page.

---

# Trade-off analysis (retained)

Feature request 2026-08-12. Current tool is one Google Doc ID ↔ one Markdown file. Multi-tab Docs would be lossy: reads default to the first tab, and `ApplyDocument` deletes/replaces `document.Body` without a tab ID.

| Today | Requested | Stable key | Nested tabs |
| --- | --- | --- | --- |
| 1:1 doc ↔ file | N tabs → directory of files | `tab_id` | skip in any v1 |

**Highest-risk existing behavior.** Import already calls Docs `batchUpdate` against the legacy body. On a multi-tab Doc that is equivalent to rewriting only the first tab, or failing in confusing ways. Any tabbed design would have to refuse unlabeled single-file apply when tabs are present.

**What would have to stay true** (if the request were ever revived): single-tab Docs keep the existing CLI; export of N top-level tabs yields N files with `tab_id`; import updates matching tabs only. Nested tabs, creating missing tabs, and Drive folder recursion stay out unless explicitly pulled in.

## Decision 1 — Local layout

| Option | Fit for minutes | Round-trip | Cost | Verdict |
| --- | --- | --- | --- | --- |
| Directory: one `.md` per tab + sidecar | Matches git meeting packs | Strong if `tab_id` is in frontmatter | New export/import path; CLI `-out` becomes a dir | Only layout that matches the request |
| One file, tab sections / delimiters | Poor: huge diffs, hard to review one meeting | Weak: order and titles collide with body | Smaller CLI change | Reject |
| Keep 1:1; one Doc per meeting | Works today; SPC already using as workaround | N/A — no tab mapping | Zero code; process cost on Docs side | **Chosen path** |

Source: `Missives/2026-08-12-feature-request-tabbed-doc-directory-sync.md` and current `cmd/gdocs-markdown-sync` + `ApplyDocument`.

## Decision 2 — Identity mapping

**Frontmatter as identity (if implemented).** Each file would carry `title`, `tab_id`, `doc_id`. Rename/move files freely; sync keys off `tab_id`, not path. Sidecar (`_doc.toml`) only for doc-level facts: `doc_id`, Doc title, tab order. If sidecar and frontmatter disagree on `tab_id`, frontmatter wins; sidecar is rebuilt on export.

**Path or registry only.** Mapping file of path ↔ `tab_id`, or slug-from-title as the key. Title-based slugs break when a tab is renamed in Docs. A central registry duplicates data git already has in each file and is easier to desync.

## Decision 3 — Nested (child) tabs

| Policy | Default in request | Risk | If ever implemented |
| --- | --- | --- | --- |
| Skip children, warn | Recommended default | Silent data omission if someone nests minutes | Warn; non-zero only with `-strict` |
| Flatten to `parent/child.md` | Configurable alternative | Import must recreate tab tree; API tab create is later-phase | Defer — needs create-tab |
| Error if any child exists | Safest / least convenient | Blocks real Docs that use grouping | Optional `-tabs=error` |

## Decision 4 — CLI surface

**Explicit mode (request sketch).** `export -doc ID -out DIR -tabs dir` and `import -doc ID -dir DIR -tabs dir`. Matches “manual-first” in MANIFEST. No surprise directory writes when `-out` is a file path today.

**Auto-detect.** If Get returns more than one top-level tab, export writes a directory. Convenient, but breaks callers that pass `-out file.md` and assume one blob. Preview also becomes ambiguous (first tab vs concatenated).

**Flags that would have been required.** Keep export/import. Add `-tabs none|dir`. Default `none`. If the Doc has multiple tabs and `-tabs none`, fail with a message naming tab count — do not silently use tab 0. Optional later: `-tab TAB_ID` for single-file round-trip of one tab without a directory.

## Decision 5 — Import mutations

| Mutation | Why it exists | Danger | If ever implemented |
| --- | --- | --- | --- |
| Replace body of tabs that have matching files | Core acceptance criterion | Same as today per tab; scope with `TabId` on `Range`/`Location` | Only safe mutation |
| Create tab when `tab_id` missing | New meetings added only in git | Tab order, permissions, and ID assignment need extra API work | Phase 2 (request already said optional) |
| Delete remote tab when local file gone | Mirror git as source of truth | Easy to destroy committee history from a bad checkout | Never |

## Adapter / model impact

`FetchDocument` is still a stub; `ApplyDocument` is live and uses `remote.Body` plus `Location`/`Range` without `TabId`. Tab support is not a thin CLI flag — it has to land in the Google adapter first. That is the main reason the request is not feasible on the current codebase.

| Layer | Today | Tabbed change |
| --- | --- | --- |
| Docs Get | No `includeTabsContent`; adapter stub ignores API | Get with `includeTabsContent`; walk `document.tabs` |
| Document model | Title + Body `[]Element` | Add `TabID` (and maybe `Tabs []Document`) or keep one Body per apply |
| `batchUpdate` | `DeleteContentRange` + inserts on default body | Set `tabId` on every `Range` and `Location`; per-tab `endIndex` |
| Frontmatter | `title` only | `title`, `doc_id`, `tab_id`; strip on import before `FromMarkdown` |
| `internal/sync` | Empty Export/Import stubs | Directory orchestration; CLI stays thin |

## Slice that was considered and declined

Fail closed on multi-tab; directory export/import; `tab_id` in frontmatter; skip nested tabs; no create/delete tabs. Order would have been: safety + real Get-with-tabs, then export N files, then import that only `batchUpdate`s listed tab IDs.

That slice is still a new mapping product on an unfinished 1:1 adapter. It is declined.

Consumer context: Co-op Strategic Planning minutes (one Doc, one tab per meeting). Non-goals from the request remain: Drive folder recursion, comments, layout fidelity.
