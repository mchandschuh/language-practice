---
schema: language-practice-step/v1
id: WI-0001
step: "08"
status: proven
skill: language-practice-debugging
depends_on: ["02", "03", "05", "07"]
target: .cline/skills/language-practice-debugging/SKILL.md
---

# WI-0001 step 08 — Add debugging skill

## Objective

Commit hypothesis-driven debugging guidance using reproduction, minimized
failures, bounded experiments, regression tests, and explicit evidence tiers.

## Commit subject

```text
WI-0001: 08: feat(repo): add debugging skill
```

## Scope

- Add `.cline/skills/language-practice-debugging/SKILL.md`.
- Preserve learning-mode boundaries for diagnosis and correction.

## Validation

```bash
npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-debugging
git diff --check
```

## Prerequisite

Testing and research policies must be available before debugging guidance is
committed.

## Agent handoff notes

### Agent execution

- Agent identity: API assistant (identity not available)
- Agent/session ID: not available
- Provider/model: not available
- Capability: repository inspection, validation, and scoped commit
- Use case: Act execution for the selected debugging-skill step
- Repository root: `/home/michael/src/github/language-practice`
- Work-item ID: `WI-0001`
- Step file: `.docs/work-items/WI-0001/steps/WI-0001-08-language-practice-debugging.md`
- Execution mode: act
- Started at (UTC): `2026-10-14T02:32:06Z`
- Finished at (UTC): `2026-10-14T02:34:00Z`
- Effective system instructions summary: Execute only the selected atomic step, preserve unrelated worktree changes, validate the final snapshot, stage explicit paths, commit completed work, and report evidence limitations truthfully.
- Effective task prompt summary: Process only WI-0001 step 08, preserve sibling skills and iterator artifacts, record a truthful handoff, run the declared validation, and commit the target skill plus this selected step record with the exact subject.

### Run reference

- Iterator run directory: `.cline-iterations/WI-0001-08-language-practice-debugging.md/20261014_023206/`
- Run ID: `20261014_023206`
- Plan session ID: not available
- Act session ID: not available
- Use the run directory's `session.json`, `events.jsonl`, prompt snapshots, and Act transcript to review the execution context.

### Work performed

- Read the selected step completely and confirmed it is atomic: one bounded debugging guidance skill and its execution record.
- Rechecked the dependency gate: steps 02, 03, 05, and 07 are available in branch history and their process, learning-first, research, and testing policies were reviewed.
- Preserved pre-existing untracked iterator artifacts, the skills README, sibling skill directories, and the unrelated `python/` directory.
- Verified the existing debugging skill content against the selected objective and recorded this execution handoff.

### Files changed

- `.cline/skills/language-practice-debugging/SKILL.md`
- `.docs/work-items/WI-0001/steps/WI-0001-08-language-practice-debugging.md`

### Validation and evidence

- `npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-debugging` — passed (`Valid skill`).
- `git diff --check` — passed for the selected skill and step record before staging.
- Evidence tier: static skill-schema and whitespace validation only; runtime activation, prompt evaluation, learner outcomes, and language-toolchain behavior are not proven.

### Decisions and references

- The step is atomic and remains within its declared ownership boundary.
- Dependency commits include `2f14e66` (02), `c86081b` (03), `f831922` (05), and `249749a` (07).
- The exact required commit subject is `WI-0001: 08: feat(repo): add debugging skill`.
- Only the target skill and this selected step record may be staged.

### Blockers and risks

- No dependency or implementation blocker identified.
- The worktree contains unrelated pre-existing untracked iterator artifacts, the skills README, sibling skills, and `python/`; none may be staged.
- Static validation does not prove behavioral quality or runtime loading.

### Exact next actions

1. Re-run validation after this handoff edit.
2. Stage only the two listed files, inspect the staged names and complete diff, and run `git diff --cached --check`.
3. Commit with the exact declared subject, inspect the commit summary and final worktree status, and report the remaining unrelated paths.
