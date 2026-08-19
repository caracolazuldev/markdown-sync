---
title: "Bug: import applies inline styles at the wrong offsets after non-ASCII text"
status: review
audience: maintainers
date: 2026-08-18
from: co-op-strategic-planning (SPC workspace)
---

# Bug: import inline style offsets drift after non-ASCII characters

## Summary

`gdocs-markdown-sync import` still paints **bold, italic, code, and links on the wrong characters** whenever a paragraph contains non-ASCII text before the styled span. ASCII-only paragraphs look fine. A related import defect **collapses Markdown hard line breaks** (trailing two spaces) into a single paragraph, which then makes later style drift worse.

This was observed after deploying the recent inline-style work (`0b6355f` / `dc6d269`) into the SPC workspace and re-importing [Business Concept Framework](https://docs.google.com/document/d/1XpVaHEp1J9iFdlKm1M4o5ijRJMjAFBwAqlJoSb07D8o/edit).

## Environment

- Consumer: `/home/mzd/Lab/co-op-strategic-planning`
- CLI: `./bin/gdocs-markdown-sync` copied from this repo (`make sync-tool`)
- Command: `gdocs-markdown-sync import -auth service -file … -doc 1XpVaHEp1J9iFdlKm1M4o5ijRJMjAFBwAqlJoSb07D8o`
- Source: `GDrive/business-concepts/business-concept-framework.md`
- HEAD at report time: `0b6355f feat: enhance Markdown inline style handling and add tests`

## Reproduction

1. Import Markdown that has a non-ASCII character **before** an inline span, for example:

   ```markdown
   Committee → Business Concepts
   **Home (git):** `GDrive/example.md`

   They are the right tools **after** we know which idea deserves that investment.
   ```

   The first sample uses `→` (U+2192). The second uses curly quotes in the same paragraph in the real doc (`“fleshed-out business concept.”` then `**after**`).

2. Export or `preview` the Doc back to Markdown.

## Expected

- `**Home (git):**` is bold on those exact characters.
- `` `GDrive/example.md` `` is Courier New on those exact characters.
- `**after**` is bold on `after` only.
- Trailing-two-space hard line breaks remain separate visual lines in Docs (or at least a newline, not a join).

## Actual

Round-trip export of the SPC Doc after a successful import (`applied 145 body elements`):

- Header lines joined into one paragraph, then styles shifted:

  `Ho**me (git): G**D`rive/business-concepts/business-concept-framework.mdCo`**mpanion list: C**oncepts to flesh-out`

  Source was five separate hard-broken lines, including `**Home (git):** \`GDrive/…\``.

- Curly quotes earlier in the paragraph shifted a later bold by four characters:

  Source: `**after** we know`  
  Export: `afte**r we **know`

- ASCII-only spans in the same Doc were correct, e.g. `**guiding document**`, `**Destination (decade-scale).**`.

The shift size matches **UTF-8 extra bytes**, not Docs UTF-16:

| Preceding character | Extra UTF-8 bytes vs runes | Observed start shift |
|---|---|---|
| `→` (U+2192) | +2 | `Home` → bold starts at `me` |
| Two curly quotes (`“` `”`) | +2 each = +4 | `after` → bold starts at `r we ` |

## Likely cause

`parseInline` records span `Offset` / `Length` as **Go string byte lengths** (`out.Len()`, `len(val)`):

```467:473:internal/markdown/adapter.go
	addText := func(val string, st inlineStyleState) {
		if val == "" {
			return
		}
		off := out.Len()
		out.WriteString(val)
		length := len(val)
```

`appendInlineStyleReqs` then treats those numbers as **rune indices**:

```149:156:internal/markdown/adapter.go
	appendInlineStyleReqs := func(textStart int64, clean string, spans []inlineSpan) {
		for _, sp := range spans {
			runes := []rune(clean)
			if sp.Offset < 0 || sp.Offset+sp.Length > len(runes) {
				continue
			}
			s := textStart + utf16Len(string(runes[:sp.Offset]))
```

For ASCII, bytes == runes, so tests pass. For `→` / curly quotes, the byte offset is larger than the rune index, so `runes[:sp.Offset]` slices too far and `UpdateTextStyle` lands late. If `Offset+Length` exceeds `len(runes)`, the span is **silently skipped**.

Table cell fill is a second copy of the same mistake and does not even convert to UTF-16: `idx + int64(sp.Offset)`.

Existing coverage (`TestBuildDocsRequests_ParagraphInlineStyles`, `TestBuildDocsRequests_HeadingInlineStyles`) only uses ASCII (`Hello **bold**`), so it cannot catch this.

### Related: hard line breaks

`extractText` concatenates adjacent `ast.Text` nodes and does not emit a newline (or space) for `HardLineBreak` / `SoftLineBreak`. `parseInline` then turns remaining hard/soft breaks into a **single space**. Markdown lines ending in two spaces therefore become one Docs paragraph. That is a separate mapping bug; it also puts more styled spans into one paragraph so the offset drift is more visible.

## Acceptance criteria

1. Import of a paragraph with BMP non-ASCII before a span (`→`, `“”`, em dash) applies bold/italic/code/link to the same rune ranges as the Markdown.
2. Import of a paragraph with a non-BMP character (e.g. an emoji, 2 UTF-16 units) still applies styles on the correct Docs UTF-16 indexes.
3. Hard line breaks (trailing two spaces) do not join the following line onto the same visual line.
4. Unit tests cover (1)–(3) at `parseInline` / `buildDocsRequests` (assert `UpdateTextStyle` start/end indexes, not merely “a bold request exists”).
5. Re-import of the SPC Business Concept Framework Doc round-trips the header block and `**after**` without shifted markers.

## Workaround today

Keep committee-facing Markdown ASCII-only in styled paragraphs, or avoid inline markers after punctuation like `→` and curly quotes. Do not treat a successful `import` as visual fidelity until this is fixed.

## References

- Mapping contract: `docs/mapping.md` (inline styles on import)
- Apply path: `internal/markdown/adapter.go` (`parseInline`, `appendInlineStyleReqs`, `buildNativeTableFillRequests`)
- Consumer Doc: https://docs.google.com/document/d/1XpVaHEp1J9iFdlKm1M4o5ijRJMjAFBwAqlJoSb07D8o/edit
- Consumer source: `/home/mzd/Lab/co-op-strategic-planning/GDrive/business-concepts/business-concept-framework.md`
