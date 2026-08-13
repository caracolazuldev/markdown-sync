# Approaches — implementation & algorithmic records

Purpose
- Record implementation-level approaches and algorithmic strategies that affect performance, cost, or correctness but may not be full architectural decisions.

When to create an Approach record
- When choosing between different implementation strategies whose tradeoffs matter (e.g., per-item API calls vs batching, synchronous vs asynchronous processing, use of caches, pagination strategies).
- When adopting or adapting a design pattern to solve a local problem (batching, circuit-breaker, retry/backoff strategy, worker pool).

Approach record format (recommended)
- id: APP-YYYY-MM-DD-short
- title: Short descriptive title
- date: YYYY-MM-DD
- author: Who proposed it
- type: `architectural` | `algorithmic`
- description: What the approach is
- rationale: Why chosen (tradeoffs: latency, cost, complexity)
- alternatives: Other approaches considered (pros/cons)
- design_patterns: list of relevant patterns (e.g., Batcher, Circuit-Breaker, Strategy, Bulkhead)
- performance_considerations: latency, throughput, memory, cost
- impact: components, tests, monitoring needed
- reassessment_triggers: conditions that should prompt re-evaluation (e.g., >X req/s, error rate >Y%)
- links: PR, issue, benchmarks, external docs

Example (APP-2026-02-23-batching)
- id: APP-2026-02-23-batching
- title: Use bulk API requests when available instead of per-item calls
- date: 2026-02-23
- author: repo-maintainer
- type: algorithmic
- description: Use the provider's batch endpoint to submit multiple items per request rather than issuing one request per item.
- rationale: Reduces network overhead and latency per item, improves throughput at scale; increases complexity for partial failure handling.
- alternatives: per-item calls (simpler, easier error handling); hybrid (small batches)
- design_patterns: Batcher, Retry with Backoff
- performance_considerations: reduces requests/sec by factor ~N, requires memory to accumulate batches
- impact: add batcher component, partial-failure handling tests, metrics for batch size/latency
- reassessment_triggers: high variance in batch sizes, increased error rates, provider rate-limit changes
- links: PR/ISSUE/benchmarks

Location and discovery
- Add approach records to `docs/decisions.md` or a dedicated `docs/approaches/` folder and link from `MANIFEST.md`.

- id: APP-2026-08-13-tab-paths
- title: Map the Docs tab tree to nested folders with a same-named parent body file
- date: 2026-08-13
- author: repo-maintainer
- type: algorithmic
- description: Walk `document.tabs` and `childTabs`. A leaf tab is `<slug>.md`. A tab with children is folder `<slug>/` plus `<slug>/<slug>.md` for the parent body (even if empty). Recurse. Identity is `tab_id`. On refresh, promote leaf→folder when children appear and flatten when they all disappear. If a child slug collides with the parent-body path, disambiguate the child (suffix a short tab id), never the parent file. `_track.toml` stays at the track root.
- rationale: Mirrors how people already nest meeting packs in git. Parent tabs can have their own body; skipping nested tabs would drop data. Flattening into one file or skipping children was rejected for the minutes workflow.
- alternatives: Skip nested tabs (lossy). Flatten to `parent-child.md` (loses hierarchy). Use `_index.md` instead of same-named parent file (less obvious). Path-only identity (breaks on rename).
- design_patterns: Composite (tab tree), Strategy (leaf vs parent path)
- performance_considerations: One `documents.get` with `includeTabsContent`; path mapping is in-memory over the tab tree.
- impact: `internal/sync` path mapping and tests; `track` writes; mapping.md
- reassessment_triggers: Docs tab depth or slug collisions in real minutes Docs; need for index files for empty parents
- links:
  - ADR-2026-08-13-track
  - Missives/2026-08-13-spec-track-tabbed-docs.md

