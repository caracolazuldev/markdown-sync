# Skills

Skills are focused, testable instructions or capabilities that an AI coding agent can perform. Each skill should be:

- Atomic: accomplishes a single, well-defined task.
- Parameterized: accepts inputs (e.g., file path, function name).
- Documented: includes examples and expected outputs.

Suggested structure for a skill:

- `skills/<skill-name>/skill.md` — human-readable description and examples.
- `skills/<skill-name>/spec.yaml` — machine-readable spec (inputs, outputs, preconditions).

Example skill names: `run-tests`, `add-dependency`, `generate-unit-test`, `refactor-function`.
