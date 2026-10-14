---
schema: language-practice-step/v1
id: WI-0001
step: "13"
status: planned
skill: language-practice-python
depends_on: ["05", "07", "08", "09", "10"]
target: .cline/skills/language-practice-python/SKILL.md
---

# WI-0001 step 13 — Add Python conventions skill

## Objective

Commit Python-specific guidance for readability, modules, data structures,
exceptions, typing, iterators, generators, environments, tests, and async code.

## Commit subject

```text
WI-0001: 13: feat(python): add Python conventions skill
```

## Scope

- Add `.cline/skills/language-practice-python/SKILL.md`.
- Preserve PEP 8 and official Python documentation anchors.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-python
git diff --check
```

## Prerequisite

Shared research, testing, debugging, observability, and reliability policies
must be committed first.
