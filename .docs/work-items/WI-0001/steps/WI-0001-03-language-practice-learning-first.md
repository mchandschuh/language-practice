---
schema: language-practice-step/v1
id: WI-0001
step: "03"
status: planned
skill: language-practice-learning-first
depends_on: ["02"]
target: .cline/skills/language-practice-learning-first/SKILL.md
---

# WI-0001 step 03 — Add learning-first skill

## Objective

Commit the learning policy that distinguishes handwrite, collaborate, and
assist modes and protects understanding as the repository's primary outcome.

## Commit subject

```text
WI-0001: 03: feat(repo): add learning-first skill
```

## Scope

- Add `.cline/skills/language-practice-learning-first/SKILL.md`.
- Keep the hand-written-code boundary explicit.
- Do not stage language-specific skills.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-learning-first
git diff --check
```

## Prerequisite

Step 02 must be committed so learning work inherits the repository process.
