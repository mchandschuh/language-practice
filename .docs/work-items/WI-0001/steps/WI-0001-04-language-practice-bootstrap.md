---
schema: language-practice-step/v1
id: WI-0001
step: "04"
status: planned
skill: language-practice-bootstrap
depends_on: ["02", "03"]
target: .cline/skills/language-practice-bootstrap/SKILL.md
---

# WI-0001 step 04 — Add bootstrap skill

## Objective

Commit the repository-initialization guidance for the four independent language
tracks without inventing versions, dependencies, or exercises.

## Commit subject

```text
WI-0001: 04: feat(repo): add bootstrap skill
```

## Scope

- Add `.cline/skills/language-practice-bootstrap/SKILL.md`.
- Keep setup separate from exercise implementation.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-bootstrap
git diff --check
```

## Prerequisite

Process and learning boundaries must already be committed.
