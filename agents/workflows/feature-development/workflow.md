# Workflow: Feature development

Purpose
- A standard, discoverable process for implementing new features or fixes.

Exit criteria
- Feature implemented, tests passing, docs updated, PR opened with reflection notes.

Steps
1. Define scope and acceptance criteria (create or update issue).
2. Create feature branch: `feature/<short-desc>`.
3. Implement changes in small commits.
4. Run `run-tests` skill regularly and fix failures (uses `run-in-devcontainer`; do not install Go on the host).
5. Update `go.mod`/`go.sum` and run `make tidy` as needed (Docker on the host).
6. Add or update documentation in `docs/` or `agents/references/`.
7. Run static checks via `make lint` (same container path; do not run `staticcheck` on the host).
8. Prepare PR with summary, tests run, and post-feature reflection (see `agents/post_feature_checklist.md`).
9. If the change touches sensitive files, include ADR/APP ids in the PR description and link to the decision/approach records (see `agents/policies/pr-guidelines.md`).

Decision points
- If dependency upgrade required, include compatibility notes in `agents/references/dependencies.md` and prefer a separate PR for major upgrades.

Referenced skills
- `run-in-devcontainer` — run Go/Make/staticcheck in `.devcontainer/Dockerfile`; never install Go on the host.
- `run-tests` — run tests and report failures.
