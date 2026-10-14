---
schema: language-practice-step/v1
id: WI-0001
step: "07"
status: proven
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

## Agent handoff notes

### Agent execution

- Agent identity: API assistant (identity not available)
- Agent/session ID: not available
- Provider/model: not available
- Capability: repository implementation and validation
- Use case: execute selected WI-0001 step
- Repository root: /home/michael/src/github/language-practice
- Work-item ID: WI-0001
- Step file: .docs/work-items/WI-0001/steps/WI-0001-07-language-practice-testing.md
- Execution mode: Act
- Started at (UTC): 2026-10-14T01:09:40Z
- Finished at (UTC): 2026-10-14T01:10:00Z
- Effective system instructions summary: Implement only the selected step, preserve unrelated work, validate the final diff, stage explicit paths, commit completed in-scope work, and report evidence limitations truthfully.
- Effective task prompt summary: Process only this authoritative testing-skill step, preserve sibling skills and iterator artifacts, record truthful validation evidence, update this handoff, and commit only the selected step's files.

### Run reference

- Iterator run directory: .cline-iterations/WI-0001-07-language-practice-testing.md/20261014_010940/
- Run ID: 20261014_010940
- Plan session ID: not available
- Act session ID: not available
- Use the run directory's session.json, events.jsonl, prompt snapshots, and Act transcript to review the execution context and reasoning. The provider session was not available; this execution used the repository state directly.

### Work performed

- Rechecked the selected step, work-item authority, current worktree, target skill, and dependency commits 02, 03, and 05.
- Confirmed the step is atomic: one bounded testing guidance skill and its execution record.
- Reviewed the existing target skill, which covers public contracts, negative and boundary cases, deterministic isolation, cleanup, regression tests, learning checkpoints, and evidence limits.
- Preserved pre-existing untracked iterator artifacts, the skills README, and sibling skill files.

### Files changed

- .cline/skills/language-practice-testing/SKILL.md
- .docs/work-items/WI-0001/steps/WI-0001-07-language-practice-testing.md

### Validation and evidence

- `npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-testing` — passed (`Valid skill`).
- `git diff --check` — passed for the selected paths before staging.
- Evidence tier: static skill-schema and whitespace validation only; this does not prove runtime activation, prompt quality, learner outcomes, or language-toolchain behavior.

### Decisions and references

- Dependencies are satisfied by commits `2f14e66` (02), `c86081b` (03), and `f831922` (05).
- The exact required commit subject is `WI-0001: 07: feat(repo): add testing skill`.
- Stage only the target skill and this selected step record; do not stage `.cline-iterations/`, `.cline/skills/README.md`, or sibling skills.

### Blockers and risks

- No implementation or dependency blocker identified.
- The worktree contains unrelated pre-existing untracked iterator artifacts and sibling skill files; they must remain uncommitted.
- Runtime activation, provider behavior, learner outcomes, and real language-toolchain execution remain unproven.

### Exact next actions

1. Stage only the two explicit paths and inspect the staged name/status and complete diff.
2. Run `git diff --cached --check`, then commit with the exact declared subject and an evidence-focused body.
3. Inspect the commit summary and final worktree status, confirming unrelated paths remain uncommitted.
