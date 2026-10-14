---
schema: language-practice-step/v1
id: WI-0001
step: "10"
status: planned
skill: language-practice-reliability-concurrency
depends_on: ["05", "07", "08", "09"]
target: .cline/skills/language-practice-reliability-concurrency/SKILL.md
---

# WI-0001 step 10 — Add reliability and concurrency skill

## Objective

Commit cross-language guidance for ownership, cancellation, timeouts,
backpressure, retries, idempotency, synchronization, shutdown, and recovery.

## Commit subject

```text
WI-0001: 10: feat(repo): add reliability skill
```

## Scope

- Add `.cline/skills/language-practice-reliability-concurrency/SKILL.md`.
- Preserve language-specific distinctions instead of flattening semantics.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-reliability-concurrency
git diff --check
```

## Prerequisite

Observability must be committed so concurrency exercises include diagnostic
evidence rather than only correctness claims.
