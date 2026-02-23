# Post-feature Checklist & Reflection

After implementing a feature or fix, agents (and maintainers) should run this checklist and record brief answers in the PR description or a follow-up issue.

1. Tests and checks
   - Did `make test` pass? (yes/no)
   - Did `staticcheck` pass? (yes/no)

2. Scope & acceptance
   - Does the change satisfy the acceptance criteria? (brief)

3. Risks and fragility
   - Any flaky tests or brittle code introduced? (brief)

4. Documentation
   - Which docs were updated? (list files)
   - Any outstanding documentation gaps? (brief)

5. Agent instructions and skills
   - Did implementing this feature expose missing or unclear agent skills/workflows? (yes/no)
   - Suggested improvements to agent instructions (short bullets).

6. Next steps
   - Follow-up tasks (issues, cleanup, monitoring).

Requirement for agents
- Agents must include answers to items 1, 3, and 5 when proposing to close a feature PR. If item 5 indicates agent guidance changes, the agent should propose an update to the relevant `agents/skills/` or `agents/workflows/` artifacts.
