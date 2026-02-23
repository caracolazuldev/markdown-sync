## Architecture & Technical Decision Policy

Purpose
- Ensure architecture and technical approach documentation remains current and that business-logic decisions are discoverable, justified, and revisited when conditions change.

Policy summary
- Every architectural or technical decision that affects system behavior, developer workflow, business logic, or implementation approach MUST be recorded in `docs/decisions.md` (for decisions) or `docs/approach.md` (for algorithmic/approach notes) and linked from `MANIFEST.md`.
- "Approach" records capture algorithmic or implementation strategies (for example: per-item API calls vs bulk requests, pagination strategies, batching, caching approaches) and should reference relevant design patterns.
- Records must include: decision title, date, author, rationale, alternative options considered, and re-evaluation conditions (triggers).
- If a change modifies or supersedes an existing decision, update `docs/decisions.md` with a new entry that references and supersedes the earlier decision.

- `id`: unique short id (e.g., `ADR-2026-02-23-auth`)
- `title`
- `date`
- `author`
- `decision` (short description)
- `rationale` (why chosen)
- `alternatives` (brief list)
- `impact` (components, tests, docs)
- `reassessment_triggers` (conditions that should prompt re-evaluation)
- `links` (PR, issue, specs)

Approach record format
- `id`: unique short id (e.g., `APP-2026-02-23-batching`)
- `title`
- `date`
- `author`
- `type`: `architectural` or `algorithmic`
- `description`: concise description of the approach
- `rationale`: why this approach was chosen (cost, latency, complexity)
- `alternatives`: other approaches considered (with brief pros/cons)
- `design_patterns`: list of design patterns applied or relevant (e.g., Bulkhead, Batcher, Circuit Breaker, Strategy)
- `impact`: components, tests, operational considerations
- `reassessment_triggers`: conditions that should prompt re-evaluation (scale, latency, error-rate)
- `links`: PR, issue, performance tests, external docs
- `id`: unique short id (e.g., `ADR-2026-02-23-auth`)
- `title`
- `date`
- `author`
- `decision` (short description)
- `rationale` (why chosen)
- `alternatives` (brief list)
- `impact` (components, tests, docs)
- `reassessment_triggers` (conditions that should prompt re-evaluation)
- `links` (PR, issue, specs)

Workflow expectation
- When proposing a change that includes business logic or architectural impact, the proposer must:
  1. Create or update a decision record in `docs/decisions.md` on the feature branch.
  2. Add a short summary and the decision `id` to the PR description.
  3. Add a link or reference to the MANIFEST under a `Decisions` section (PR reviewers will confirm).

Verification
- CI or maintainers should verify that PRs with architecture/logic changes include an updated or new decision record and a MANIFEST update linking to it.
