---
schema: language-practice-step/v1
id: WI-0001
step: "08"
status: planned
skill: language-practice-debugging
depends_on: ["02", "03", "05", "07"]
target: .cline/skills/language-practice-debugging/SKILL.md
---

# WI-0001 step 08 — Add debugging skill

## Objective

Commit hypothesis-driven debugging guidance using reproduction, minimized
failures, bounded experiments, regression tests, and explicit evidence tiers.

## Commit subject

```text
WI-0001: 08: feat(repo): add debugging skill
```

## Scope

- Add `.cline/skills/language-practice-debugging/SKILL.md`.
- Preserve learning-mode boundaries for diagnosis and correction.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-debugging
git diff --check
```

## Prerequisite

Testing and research policies must be available before debugging guidance is
committed.
