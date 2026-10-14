---
schema: language-practice-step/v1
id: WI-0001
step: "02"
status: planned
skill: language-practice-process
depends_on: ["01"]
target: .cline/skills/language-practice-process/SKILL.md
---

# WI-0001 step 02 — Add development process skill

## Objective

Commit the repository-wide process skill after Git hygiene exists. It defines
scope capture, bounded implementation, evidence states, side-effect controls,
and completion reporting.

## Commit subject

```text
WI-0001: 02: feat(repo): add development process skill
```

## Scope

- Add `.cline/skills/language-practice-process/SKILL.md`.
- Preserve its delegation to `language-practice-git`.
- Do not stage other skills.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-process
git diff --check
```

## Prerequisite

Step 01 must be committed and remain available in the branch history.
