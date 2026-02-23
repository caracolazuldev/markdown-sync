## Agent Practices

This document defines expectations for AI coding agents and automation operating in this repository.

Core principles
- Human-in-the-loop: Agents must consult a human for actions that change protected branches or modify sensitive configuration (credentials, CI, release tags).
- Least surprise: Changes should be minimal, well-explained, and reversible.
- Test-before-request: Agents must run local tests and static checks before proposing changes.
- Document intent: Every proposed change must include a short intent statement and expected outcomes.

Operational rules
- Work on feature branches by default. If working on `main`, pause and ask a human to commit/push.
- Include `go.mod` and `go.sum` changes in the same feature branch as code changes.
- For dependency upgrades, include upgrade rationale and compatibility notes in `agents/references/dependencies.md`.
- If a change affects developer workflow (Makefile, devcontainer, CI), include migration notes in `docs/` and notify maintainers.

Post-feature reflection
- After delivering a feature, the agent must run the `agents/post_feature_checklist.md` and produce a short reflection noting:
  - What worked well?
  - What failed or was fragile?
  - Suggested updates to skills or workflows.
  - Any documentation gaps.

These reflections should be included in the PR description or a follow-up issue for human review.
