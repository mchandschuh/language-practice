---
schema: language-practice-step/v1
id: WI-0001
step: "07"
status: planned
skill: language-practice-testing
depends_on: ["02", "03", "05"]
target: .cline/skills/language-practice-testing/SKILL.md
---

# WI-0001 step 07 — Add testing skill

## Objective

Commit testing guidance that treats tests as learning artifacts and requires
public contracts, negative paths, invariants, regressions, and honest evidence.

## Commit subject

```text
WI-0001: 07: feat(repo): add testing skill
```

## Scope

- Add `.cline/skills/language-practice-testing/SKILL.md`.
- Keep tests deterministic, bounded, isolated, and explicit about limitations.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-testing
git diff --check
```

## Prerequisite

Research must be committed so testing guidance can defer to current toolchain
documentation.
