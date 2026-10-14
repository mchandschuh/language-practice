---
schema: language-practice-step/v1
id: WI-0001
step: "12"
status: planned
skill: language-practice-go
depends_on: ["05", "07", "08", "09", "10"]
target: .cline/skills/language-practice-go/SKILL.md
---

# WI-0001 step 12 — Add Go conventions skill

## Objective

Commit Go-specific conventions for packages, interfaces, errors, modules,
goroutines, channels, contexts, tests, and diagnostics.

## Commit subject

```text
WI-0001: 12: feat(go): add Go conventions skill
```

## Scope

- Add `.cline/skills/language-practice-go/SKILL.md`.
- Preserve official Go documentation anchors and learning boundaries.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-go
git diff --check
```

## Prerequisite

Research, testing, debugging, observability, and reliability guidance must be
committed first.
