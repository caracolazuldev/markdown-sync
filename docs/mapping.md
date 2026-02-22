# Mapping Rules — markdown-sync

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
- Frontmatter → generated YAML containing `title`, `doc_id`, `last_modified`

Unsupported / Limitations
-------------------------
- Complex layout elements (textboxes, positioned images) are not preserved.
- Comments and suggestions are not mapped in v0.1 (future work).

Round-trip / determinism
------------------------
Mapping is designed to be idempotent: repeated export of the same document yields the same Markdown (modulo non-deterministic metadata like timestamps).
