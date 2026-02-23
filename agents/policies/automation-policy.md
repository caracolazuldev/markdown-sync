## Automation Policy


Effective rules for automated tools and AI assistants operating on this repository:

- Do not commit, push, or merge directly to `main` or other protected branches.
- Automation SHOULD create and operate on a feature branch (naming convention: `feature/*`, `fix/*`, `chore/*`) when feasible.
- If the assistant or automation is operating in the context of a protected branch (for example, `main`) and a feature branch was not created, the assistant MUST pause and ask a human to review and perform the commit/push/merge on the protected branch unless explicitly instructed to create and use a feature branch.
- All changes to `go.mod`/`go.sum` should be prepared on the feature branch and included in the PR.
- Automation must run tests (`make test`) and static analysis (`staticcheck`) locally before opening a PR or requesting a human to commit on a protected branch.
- Automation must not push or merge without an explicit human approval recorded in the PR.


Enforcement aid:

- CI should enforce these policies by refusing merges that bypass PR checks.

Agent commit behavior (mandatory)

- Agents must not run or request git commands that perform commits (for example `git commit`, `git push`, or wrappers that run them).
- When changes are ready to be committed, the agent must always ask a human to make the commit and provide:
	- A proposed commit message (imperative style),
	- A short summary of staged or proposed changes,
	- A suggested branch name (if creating a feature branch is appropriate).
- If the agent is working with `HEAD` on a protected branch (e.g., `main`) and no feature branch exists, the agent must not attempt to create or commit; instead it must ask the human to create a feature branch and perform the commit.

If you are implementing automation, follow these rules and include a short note in the PR describing the automation's purpose and the human approver.

Post-feature requirement
- Agents must run the `agents/post_feature_checklist.md` after implementing a feature and include the checklist answers and any suggested updates to agent instructions in the PR description.

PR checklist requirement
- When a change touches sensitive files (code under `internal/`, `docs/`, `.devcontainer/`, `go.mod`, `go.sum`, `MANIFEST.md`, `Makefile`), the PR must include one or more of the following in the PR description:
	- A reference to an `ADR-YYYY-MM-DD-...` id documenting the business/architectural decision, or
	- A reference to an `APP-YYYY-MM-DD-...` id documenting the implementation/approach choice.

Suggested PR content:
- Short summary of the change.
- Tests run and result (brief).
- Links to affected decision/approach records (ADR/APP ids) or a note if none exist and one will be added.

Guidance for maintainers and automation:
- If maintainers set up CI checks to enforce this policy, they can require that PR descriptions include `ADR-` or `APP-` ids when the sensitive file patterns are modified. See `agents/policies/pr-guidelines.md` for a suggested, human-readable enforcement procedure and example grep commands.


