# Design Patterns — practical notes and examples

Purpose
- Capture design patterns and concrete usage notes for this project. Include when to use, tradeoffs, and small code/architecture sketches where useful.

Suggested entries

- Batcher
  - When: network-bound APIs that support bulk requests.
  - Why: reduces per-item latency and overhead; improves throughput.
  - Tradeoffs: complexity for partial failures, memory for buffering, need for backpressure.

- Circuit Breaker
  - When: remote services that may intermittently fail.
  - Why: avoids cascading failures and provides graceful degradation.
  - Tradeoffs: increases state management; needs tuning of thresholds.

- Retry with Backoff
  - When: transient errors from upstream services.
  - Why: improves reliability without immediate failure.
  - Tradeoffs: increases latency; must avoid amplifying load.

- Strategy Pattern
  - When: multiple interchangeable algorithms or API variants (per-item vs bulk).
  - Why: cleanly encapsulates variants and simplifies testing.

How to use
- Reference a pattern from approach records (`design_patterns` field) and link to small examples here.
