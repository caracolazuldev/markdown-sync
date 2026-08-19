---
title: "Enhancement: reconstruct lists on export"
status: implemented
audience: maintainers
date: 2026-08-18
from: gdocs-markdown-sync maintainers
---

# Enhancement: reconstruct lists on export

## Summary

Export should emit Markdown list items (`- `, `1. `, nested indent) from Google Docs bullets instead of flattening every non-heading paragraph to a plain `Paragraph`. Import already creates bullets; this is **block-level export** only.

Tracked as roadmap **R3**. Implemented: export reconstructs Markdown list items from Docs bullets.

## Motivation

Compose-in-Markdown lists import correctly via `CreateParagraphBullets`. Export / `track` refresh then drops list structure: bullets come back as ordinary paragraphs with no `- ` / `1. `. Inline styles inside those lines can still round-trip; the list is lost. Committee Docs that mix prose and bullets cannot stay in git as lists.

## Current behavior

- Import: `ListItem` → tab-prefixed `InsertText` + `CreateParagraphBullets` (`BULLET_DISC_CIRCLE_SQUARE` or `NUMBERED_DECIMAL_ALPHA_ROMAN`) in `internal/markdown/adapter.go`.
- Export: `bodyToDocument` in `internal/markdown/fetch_api.go` never reads `Paragraph.Bullet`. Every non-heading paragraph becomes `Paragraph{Text: paragraphText(...)}`.
- `DocumentToMarkdown` already knows how to serialize `ListItem` (`- `, `1. `, two-space indent per `Level`).

## Proposed behavior

### Docs API shape

- `paragraph.Bullet != nil` means a list item.
- Nesting: `Bullet.NestingLevel` (0-based). Import already uses `ListItem.Level` as a tab prefix.
- Ordered vs unordered is **not** a boolean on the paragraph. `Bullet.ListId` indexes `document.Lists[listId].ListProperties.NestingLevels[i].GlyphType` (`DECIMAL`, `ALPHA`, `ROMAN`, unspecified / bullet glyphs).
- `bodyToDocument` currently takes only `*docs.Body`. Glyph lookup needs the parent `*docs.Document` (or a `map[listId]properties`). `fetchedFromAPI` / `convertAPITabs` must pass that through. Tabbed docs: list IDs live on the tab’s document, not the legacy body.

### Mapping

- If `Bullet != nil` → append `ListItem{Ordered, Level: int(NestingLevel), Text: paragraphText(...)}` instead of `Paragraph`.
- Contiguous items stay separate `ListItem` elements (existing model).
- Do **not** put Docs tab characters into `Text`. Import added tabs only as insert prefixes. If export `TextRun`s include leading `\t`, strip them. Indent is usually paragraph style, not content.
- If `NamedStyleType` is a heading **and** `Bullet` is set, **heading wins**.
- Nested lists: `Level` from `NestingLevel`. Mixed ordered/unordered at different levels: per-item `Ordered` from that level’s `GlyphType`.

## Tests

- Fake `docs.Document` with `Bullet{ListId, NestingLevel}` plus a `Lists` map: unordered level 0, nested level 1, ordered decimal. `bodyToDocument` returns matching `ListItem`s.
- Inline bold inside a bullet: `Text` contains `**bold**`.
- `DocumentToMarkdown` emits `- `, indented `- `, and `1. `.
- Regression: non-bullet `NORMAL_TEXT` remains a `Paragraph`.

## Files

- `internal/markdown/fetch_api.go` (`bodyToDocument` signature, paragraph branch)
- `internal/markdown/fetched_test.go`
- `docs/mapping.md`

Import path unchanged except if tab-prefix leakage is discovered on export.

## Prerequisites

Inline style export (markers in `paragraphText`) and UTF-16 import offsets are done. List **item text** already round-trips as a plain paragraph; this work restores structure.

## Out of scope

- Strikethrough / underline / sub / super (see `Missives/2026-08-18-enhancement-strikethrough-sub-super.md`).
- Fenced code blocks, images, blockquotes, and thematic breaks on export (same class of block mapping, not requested here).
