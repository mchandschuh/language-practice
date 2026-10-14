---
schema: language-practice-step/v1
id: WI-0001
step: "15"
status: planned
skill: language-practice-typescript
depends_on: ["05", "07", "08", "09", "10"]
target: .cline/skills/language-practice-typescript/SKILL.md
---

# WI-0001 step 15 — Add TypeScript conventions skill

## Objective

Commit TypeScript-specific guidance for strictness, type modeling, narrowing,
modules, declarations, JavaScript runtime behavior, tests, and async code.

## Commit subject

```text
WI-0001: 15: feat(typescript): add TypeScript conventions skill
```

## Scope

- Add `.cline/skills/language-practice-typescript/SKILL.md`.
- Preserve the distinction between static types and runtime validation.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-typescript
git diff --check
```

## Prerequisite

Shared research, testing, debugging, observability, and reliability policies
must be committed first.
