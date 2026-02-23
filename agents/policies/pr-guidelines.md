## PR Guidelines: requiring ADR / APP references for sensitive changes

Purpose
- Provide a simple, discoverable, and human-readable procedure that maintainers and agents should follow when a PR modifies sensitive parts of the repo.

Sensitive files (example list)
- `internal/**`
- `docs/**`
- `.devcontainer/**`
- `go.mod`, `go.sum`
- `MANIFEST.md`
- `Makefile`

Procedure for contributors and agents
1. If you change any file matching the sensitive list, add (on the feature branch) a decision (`ADR-...`) or approach (`APP-...`) record before opening the PR. If you are unsure whether the change is architectural (business logic) or algorithmic (technical approach), add both and link them.
2. In the PR description include the ADR/APP id(s) and a one-line summary of the rationale.
3. Run `make test` and `staticcheck` locally and include a short test-results section in the PR description (pass/fail and any remaining issues).
4. Run the `agents/post_feature_checklist.md` and paste the checklist answers in the PR description.

Suggested quick enforcement (manual or CI)
- Maintainers can add a lightweight check that scans changed files and the PR body for ADR/APP ids. Example one-liner (local/CI):

  - List changed files relative to base branch:

    git fetch origin main && git diff --name-only origin/main..HEAD

  - If any file matches the sensitive patterns, verify PR body contains ADR/APP id (simple grep):

    grep -E '(ADR|APP)-[0-9]{4}-[0-9]{2}-[0-9]{2}' pr_body.txt || echo "Missing ADR/APP id"

PR template snippet suggestion
---------------------------
Add the following to any PR template used by the repo. It prompts contributors to include the needed references.

```
Summary:

Tests:
- make test: (pass/fail)
- staticcheck: (pass/fail)

Decisions/Approaches:
- ADR/APP ids (if applicable):

Post-feature checklist answers:
- (paste results from agents/post_feature_checklist.md)

```

Notes
- This is a human-friendly guideline. Automation should enforce it in CI where practical, but the text here is the canonical, repo-local spec for agents and contributors.
