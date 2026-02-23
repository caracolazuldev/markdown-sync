# Agents — AI coding agent resources

This folder contains curated resources for AI coding agents used by the project.

- `skills/` — task-focused instructions the agent can execute (small, re-usable).
- `workflows/` — multi-step processes and rules for common activities (define, plan, implement, document features).
- `references/` — project-specific documentation and dependency notes that augment official docs.

See `MANIFEST.md` for discovery rules.


Automation policy and guard
-------------------------
Automation must not commit or push directly to `main`. See `agents/policies/automation-policy.md` for the project policy and `agents/scripts/guard_commit.sh` for a local guard script.

Behavior expectation for assistants
----------------------------------
- Prefer creating and working on a feature branch for any changes.
- If the assistant is operating while `HEAD` is on a protected branch (for example, `main`) and a feature branch was not created, the assistant must pause and ask a human to review and perform the commit/push/merge unless explicitly instructed to create and use a feature branch.

This ensures automation does not make irreversible changes to protected branches without a human in the loop.

