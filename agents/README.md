# Agents — AI coding agent resources

This folder contains curated resources for AI coding agents used by the project.

- `skills/` — task-focused instructions the agent can execute (small, re-usable).
- `workflows/` — multi-step processes and rules for common activities (define, plan, implement, document features).
- `references/` — project-specific documentation and dependency notes that augment official docs.

See `MANIFEST.md` for discovery rules.



Automation policy and commit behavior
------------------------------------
Automation must not commit or push directly to `main`. See `agents/policies/automation-policy.md` for the project policy.

Assistant commit behavior
-------------------------
- Assistants must not run or request git commands that perform commits (for example `git commit`, `git push`, or wrappers that run them).
- Instead, when changes are ready, the assistant must always pause and ask a human to make the commit. The assistant should provide:
	- A proposed commit message (brief, imperative style),
	- A short summary of the changes staged or to be committed,
	- Any suggested branch name (e.g., `feature/<short-desc>`) if a feature branch is appropriate.
- If the assistant is operating on a protected branch (for example, `main`) and no feature branch exists, the assistant must not perform or request a commit; it must ask the human to create the feature branch and perform the commit.

This ensures humans make commits while still enabling the assistant to prepare and propose accurate commit messages.

