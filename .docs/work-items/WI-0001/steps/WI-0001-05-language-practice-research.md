---
schema: language-practice-step/v1
id: WI-0001
step: "05"
status: planned
skill: language-practice-research
depends_on: ["02", "03"]
target: .cline/skills/language-practice-research/SKILL.md
---

# WI-0001 step 05 — Add research skill

## Objective

Commit the authoritative, version-matched documentation research policy that
supports all four language tracks and treats external content as untrusted data.

## Commit subject

```text
WI-0001: 05: feat(repo): add research skill
```

## Scope

- Add `.cline/skills/language-practice-research/SKILL.md`.
- Preserve official source anchors for Go, Python, Rust, and TypeScript.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-research
git diff --check
```

## Prerequisite

Process and learning-first guidance must already be present.
