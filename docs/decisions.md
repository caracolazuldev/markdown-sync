# Decisions (Architecture & Technical)

This file records architecture decisions and business-logic rationales. Each entry should follow the Decision record format described in `docs/architecture-policy.md`.


Example decision entry

- id: ADR-2026-02-23-oauth2
- title: Upgrade `golang.org/x/oauth2` to v0.35.0 and require Go 1.24 in devcontainer
- date: 2026-02-23
- author: repo-maintainer
- decision: Upgrade oauth2 to obtain improved OAuth features; require Go 1.24 in the devcontainer and update `go.mod`.
- rationale: v0.35.0 contains bug fixes and features needed for service flows; requires Go 1.24 which is acceptable for our devcontainer.
- alternatives: pin older oauth2 and avoid toolchain bump (rejected due to missing features).
- impact: update `.devcontainer/Dockerfile`, `go.mod`, and tests; update README and MANIFEST.
- reassessment_triggers:
  - security advisory on oauth2
  - breaking change in oauth2 or Go toolchain
  - CI/test flakiness introduced by the upgrade
- links:
  - PR: TODO

Example approach entry

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
- links:
  - PR: TODO
