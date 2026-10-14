---
schema: language-practice-step/v1
id: WI-0001
step: "06"
status: planned
skill: language-practice-cline-iterate
depends_on: ["01", "02", "03"]
target: .cline/skills/language-practice-cline-iterate/SKILL.md
---

# WI-0001 step 06 — Add cline-iterate skill

## Objective

Commit the bounded `cline-iterate` protocol for ordered, isolated, resumable
exercise processing with plan/Act handoffs and explicit failure boundaries.

## Commit subject

```text
WI-0001: 06: feat(repo): add cline iterate skill
```

## Scope

- Add `.cline/skills/language-practice-cline-iterate/SKILL.md`.
- Define worker scope, retry, transcript, and target-isolation rules.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-cline-iterate
git diff --check
```

## Prerequisite

Git, process, and learning policies must already be committed.
