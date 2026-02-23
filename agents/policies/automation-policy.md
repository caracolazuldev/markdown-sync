## Automation Policy


Effective rules for automated tools and AI assistants operating on this repository:

- Do not commit, push, or merge directly to `main` or other protected branches.
- Automation SHOULD create and operate on a feature branch (naming convention: `feature/*`, `fix/*`, `chore/*`) when feasible.
- If the assistant or automation is operating in the context of a protected branch (for example, `main`) and a feature branch was not created, the assistant MUST pause and ask a human to review and perform the commit/push/merge on the protected branch unless explicitly instructed to create and use a feature branch.
- All changes to `go.mod`/`go.sum` should be prepared on the feature branch and included in the PR.
- Automation must run tests (`make test`) and static analysis (`staticcheck`) locally before opening a PR or requesting a human to commit on a protected branch.
- Automation must not push or merge without an explicit human approval recorded in the PR.

Enforcement aid:

- Use the guard script at `agents/scripts/guard_commit.sh` to prevent accidental commits on `main`.
- CI should enforce these policies by refusing merges that bypass PR checks.

If you are implementing automation, follow these rules and include a short note in the PR describing the automation's purpose and the human approver.

Post-feature requirement
- Agents must run the `agents/post_feature_checklist.md` after implementing a feature and include the checklist answers and any suggested updates to agent instructions in the PR description.

