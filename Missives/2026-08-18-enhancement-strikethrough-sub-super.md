---
title: "Enhancement: strikethrough, underline, subscript, and superscript"
status: implemented (strikethrough); underline/sub/super deferred
audience: maintainers
date: 2026-08-18
from: gdocs-markdown-sync maintainers
---

# Enhancement: strikethrough, underline, subscript, superscript

## Summary

Extend inline mapping beyond bold, italic, code, and links so `~~strike~~` (and, if we choose a dialect, underline / sub / super) round-trips through Google Docs `TextStyle`. UTF-16 span offsets are already correct; this is a **mapping-vocabulary** change.

Tracked as roadmap **R2**. **Strikethrough** (`~~text~~`) is implemented in both directions. Underline, subscript, and superscript remain deferred until a Markdown dialect is chosen.

## Motivation

Docs authors use strikethrough, underline, and baseline offset (sub/super) routinely. Export currently ignores those fields; import has no Goldmark strikethrough extension, so `~~text~~` is stripped or left as literal tildes. SPC and other Markdown-first workflows cannot preserve redlines or scientific notation.

## Current behavior

- `parseInline` / `extractText` handle emphasis, code spans, and links only.
- `parseInline` uses `goldmark.DefaultParser()` (no GFM strikethrough).
- `FromMarkdown` enables `extension.Table` only.
- Export `runStyle` in `internal/markdown/fetch_api.go` records bold, italic, Courier New, and link URL. `Strikethrough`, `Underline`, and `BaselineOffset` are unused.
- `docs/mapping.md` lists strikethrough, underline, subscript, and superscript as unsupported.

## Proposed behavior

### Import (Markdown → Docs)

1. Enable Goldmark strikethrough in both `FromMarkdown` and `parseInline`: `goldmark.WithExtensions(extension.Strikethrough, extension.Table)` (parseInline today has no Table either; keep parsers aligned).
2. `extractText`: serialize strikethrough as `~~text~~` so markers survive the `Document` model (same markdown-in-strings convention as bold).
3. `parseInline`: walk `*extast.Strikethrough`, set `inlineStyleState.strike`, emit span kind `"strike"`.
4. `appendInlineStyleRequests` / table fill: `UpdateTextStyle` with `TextStyle.Strikethrough = true`, `Fields: "strikethrough"`.
5. Underline, sub, super are **not** CommonMark. Decide in the implementation plan:
   - Skip until we adopt a dialect (`<u>` / `<sub>` / `<sup>`, or `~sub~` / `^super^`).
   - Or export-only (lossy Docs → Markdown) even if import cannot round-trip.

### Export (Docs → Markdown)

1. Extend `runStyle` with `strike`, `underline`, `baseline` from `TextStyle.Strikethrough`, `Underline`, `BaselineOffset` (`SUPERSCRIPT` / `SUBSCRIPT`).
2. Merge adjacent runs that share the full style set (already done for the existing four styles).
3. Wrap after code, before link: e.g. `~~**bold strike**~~`, then `[...](url)`.
4. Sub/super: if no Markdown dialect is chosen, export as HTML `<sub>`/`<sup>` (requires Goldmark HTML) or omit. Do not invent syntax in an ad-hoc way.

## Tests

- ASCII and non-ASCII before `~~strike~~` (UTF-16 + new span kind).
- Export fake `TextRun` with `Strikethrough: true`.
- Document underline/sub/super as unsupported until a dialect is chosen.

## Files

- `internal/markdown/adapter.go` (`parseInline`, style apply)
- `internal/markdown/convert.go` (`extractText`, Goldmark extensions)
- `internal/markdown/fetch_api.go` (`styledRunToMarkdown` / `runStyle`)
- `docs/mapping.md`

## Prerequisites

Import inline styles already use UTF-16 indexes. Do not regress those tests.

## Out of scope

- List reconstruction on export (see `Missives/2026-08-18-enhancement-export-lists.md`).
- Fenced code blocks, images, blockquotes, thematic breaks on export.
