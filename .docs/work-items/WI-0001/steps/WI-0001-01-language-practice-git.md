---
schema: language-practice-step/v1
id: WI-0001
step: "01"
status: planned
skill: language-practice-git
depends_on: []
target: .cline/skills/language-practice-git/SKILL.md
---

# WI-0001 step 01 — Add Git hygiene skill

## Objective

Commit the repository-specific Git hygiene skill that governs work-item-linked
branches, validated checkpoint commits, staging review, subject/body formats,
and verified integration.

## Commit subject

```text
WI-0001: 01: feat(repo): add Git hygiene skill
```

## Scope

- Add `.cline/skills/language-practice-git/SKILL.md`.
- Do not stage `.docs`, other skills, or unrelated work.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-git
git diff --check
```

## Completion contract

The step is complete only when the target skill validates, the staged diff
contains only its declared target and any truthful execution update, and the
commit subject exactly matches the declared subject.
