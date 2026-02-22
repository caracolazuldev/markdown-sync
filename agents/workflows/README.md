# Workflows

Workflows are rule-driven, step-by-step processes for common activities. Each workflow should:

- Define goals and exit criteria.
- Provide ordered steps and decision points.
- Reference skills to perform atomic actions.

Suggested structure for a workflow:

- `workflows/<workflow-name>/workflow.md` — human-readable steps and decision rules.
- `workflows/<workflow-name>/workflow.yaml` — machine-readable process graph (steps, transitions, required skills).

Example workflows: `feature-development`, `security-audit`, `release-cut`, `context-update`.
