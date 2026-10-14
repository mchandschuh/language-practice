---
schema: language-practice-step/v1
id: WI-0001
step: "11"
status: planned
skill: language-practice-infrastructure-design
depends_on: ["02", "03", "05", "07", "09", "10"]
target: .cline/skills/language-practice-infrastructure-design/SKILL.md
---

# WI-0001 step 11 — Add infrastructure design skill

## Objective

Commit infrastructure-oriented design guidance for APIs, queues, storage,
caching, scheduling, failure budgets, and operational evidence.

## Commit subject

```text
WI-0001: 11: feat(repo): add infrastructure design skill
```

## Scope

- Add `.cline/skills/language-practice-infrastructure-design/SKILL.md`.
- Keep exercises bounded, local, disposable, and learning-focused.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-infrastructure-design
git diff --check
```

## Prerequisite

Reliability and observability guidance must be committed before infrastructure
design guidance.
