---
schema: language-practice-step/v1
id: WI-0001
step: "16"
status: planned
skill: language-practice-cross-language-comparison
depends_on: ["05", "12", "13", "14", "15"]
target: .cline/skills/language-practice-cross-language-comparison/SKILL.md
---

# WI-0001 step 16 — Add cross-language comparison skill

## Objective

Commit the comparison method that relates shared behavioral contracts across
Go, Python, Rust, and TypeScript without pretending their semantics are equal.

## Commit subject

```text
WI-0001: 16: feat(repo): add cross language skill
```

## Scope

- Add `.cline/skills/language-practice-cross-language-comparison/SKILL.md`.
- Require comparison of types, ownership, errors, concurrency, tooling, and
  runtime behavior.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-cross-language-comparison
git diff --check
```

## Prerequisite

All four language skills and the shared research policy must be committed before
this comparison skill is processed.
