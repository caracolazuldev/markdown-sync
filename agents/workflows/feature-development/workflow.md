# Workflow: Feature development

Purpose
- A standard, discoverable process for implementing new features or fixes.

Exit criteria
- Feature implemented, tests passing, docs updated, PR opened with reflection notes.

Steps
1. Define scope and acceptance criteria (create or update issue).
2. Create feature branch: `feature/<short-desc>`.
3. Implement changes in small commits.
4. Run `run-tests` skill regularly and fix failures.
5. Update `go.mod`/`go.sum` and run `go mod tidy` as needed.
6. Add or update documentation in `docs/` or `agents/references/`.
7. Run static checks (`staticcheck`).
8. Prepare PR with summary, tests run, and post-feature reflection (see `agents/post_feature_checklist.md`).

Decision points
- If dependency upgrade required, include compatibility notes in `agents/references/dependencies.md` and prefer a separate PR for major upgrades.

Referenced skills
- `run-tests` — run tests and report failures.
