---
schema: language-practice-step/v1
id: WI-0001
step: "09"
status: planned
skill: language-practice-observability-performance
depends_on: ["05", "07", "08"]
target: .cline/skills/language-practice-observability-performance/SKILL.md
---

# WI-0001 step 09 — Add observability and performance skill

## Objective

Commit measurement guidance for logs, metrics, traces, profiling, benchmarks,
baselines, noise, and evidence boundaries.

## Commit subject

```text
WI-0001: 09: feat(repo): add observability skill
```

## Scope

- Add `.cline/skills/language-practice-observability-performance/SKILL.md`.
- Keep measurement separate from unsupported production-readiness claims.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-observability-performance
git diff --check
```

## Prerequisite

Research, testing, and debugging skills must be committed first.
