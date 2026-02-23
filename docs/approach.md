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
