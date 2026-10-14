---
schema: language-practice-step/v1
id: WI-0001
step: "14"
status: planned
skill: language-practice-rust
depends_on: ["05", "07", "08", "09", "10"]
target: .cline/skills/language-practice-rust/SKILL.md
---

# WI-0001 step 14 — Add Rust conventions skill

## Objective

Commit Rust-specific guidance for ownership, borrowing, lifetimes, enums,
traits, errors, Cargo, tests, concurrency, and async runtimes.

## Commit subject

```text
WI-0001: 14: feat(rust): add Rust conventions skill
```

## Scope

- Add `.cline/skills/language-practice-rust/SKILL.md`.
- Preserve compiler diagnostics as a learning resource.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-rust
git diff --check
```

## Prerequisite

Shared research, testing, debugging, observability, and reliability policies
must be committed first.
