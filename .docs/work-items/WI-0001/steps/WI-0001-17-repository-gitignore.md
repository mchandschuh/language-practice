---
schema: language-practice-step/v1
id: WI-0001
step: "17"
status: proven
skill: repository-gitignore
depends_on: ["16"]
target: .gitignore
---

# WI-0001 step 17 — Ignore iterator artifacts

## Objective

Add the first repository-level ignore rule so local `cline-iterate` receipts,
transcripts, manifests, and monitor artifacts do not appear as unrelated
worktree changes.

## Commit subject

```text
WI-0001: 17: chore(repo): ignore iterator artifacts
```

## Scope

- Add a root `.gitignore` containing `/\.cline-iterations/`.
- Do not add language-toolchain patterns; those belong to WI-0002.
- Do not ignore tracked files or unrelated user work.

## Validation

```bash
git check-ignore -v .cline-iterations/example/session.json
git status --short --untracked-files=all
git diff --check
```

The ignore check must identify the root `.gitignore` rule. Existing tracked
files are unaffected by a new ignore pattern.

## Prerequisite

The WI-0001 root documentation update must be committed, and step 16 must
remain available in branch history.

## Agent handoff notes

### Agent execution

- Agent identity: API assistant (identity not available)
- Agent/session ID: `1791948563232_7w4fm`
- Provider/model: `openrouter / openai/gpt-5.6-luna`
- Capability: repository inspection, scoped implementation, validation, handoff, and commit
- Use case: Act execution for the selected repository gitignore step
- Repository root: `/home/michael/src/github/language-practice`
- Work-item ID: `WI-0001`
- Step file: `.docs/work-items/WI-0001/steps/WI-0001-17-repository-gitignore.md`
- Execution mode: `Act`
- Started at (UTC): `2026-10-14T03:29:23Z`
- Finished at (UTC): `2026-10-14T03:30:02Z`
- Effective system instructions summary: Execute only the selected atomic step, preserve unrelated work, validate the final snapshot, update this handoff truthfully, and commit only explicit in-scope paths.
- Effective task prompt summary: Add the root iterator-artifact ignore rule, run the declared Git checks, record evidence and limitations, and commit the selected step with its exact subject.

### Run reference

- Iterator run directory: `/tmp/language-practice-cline-runs/WI-0001-step17/WI-0001-17-repository-gitignore.md/20261014_032835/`
- Run ID: `20261014_032835`
- Plan session ID: `1791948515977_j8fkh`
- Act session ID: `1791948563232_7w4fm`
- Review `session.json`, `events.jsonl`, monitor summaries, prompt snapshots, and transcripts in the external run directory for execution context.

### Work performed

- Rechecked the clean isolated worktree and dependency gate before editing.
- Added the root `/.cline-iterations/` ignore rule.
- Updated this selected step record with the proven status and execution handoff.
- Preserved all sibling steps, skills, generated artifacts, and unrelated files.

### Files changed

- `.gitignore`
- `.docs/work-items/WI-0001/steps/WI-0001-17-repository-gitignore.md`

### Validation and evidence

- `git check-ignore -v .cline-iterations/example/session.json` — passed; the root `.gitignore` rule matched.
- `git status --short --untracked-files=all` — passed in the isolated execution worktree; only the selected `.gitignore` and step record were implementation changes before commit.
- `git diff --check` — passed.
- `git diff --cached --check` — passed before the worker commit.
- `cline-iterate` session manifest — completed; post-step workspace status was clean and no forced recovery was used.
- Evidence tier: static Git ignore behavior, worktree hygiene, and scoped delivery validation only; runtime iterator behavior, provider behavior, learner outcomes, and production safety are not proven.

### Decisions and references

- The step remained atomic and within its declared ownership boundary.
- Step 16 commit `265e215` and the WI-0001 root documentation commit `406da26` were verified before implementation.
- The exact required commit subject is `WI-0001: 17: chore(repo): ignore iterator artifacts`.
- Only the root `.gitignore` and this selected step record are staged for the primary checkpoint.

### Blockers and risks

- No dependency, validation, or commit-scope blocker identified.
- The primary worktree contains unrelated pre-existing untracked paths; the new ignore rule intentionally hides only `.cline-iterations/` from status and does not stage or delete it.
- Static validation does not prove runtime iterator behavior or learner outcomes.

### Exact next actions

1. Inspect the primary checkpoint commit and append the advisor review before integration.

## Advisor review

### Review round 1

- Advisor: API assistant (identity not available)
- Work-item: `WI-0001`
- Step: `17`
- Iterator run: `/tmp/language-practice-cline-runs/WI-0001-step17/WI-0001-17-repository-gitignore.md/20261014_032835/`
- Implementation commit: `bd1e28914d38a7871b46bebeb769a32679e562c2`
- Reviewed snapshot: primary branch `wi-0001-gitignore`, commit `bd1e289`
- Reviewed at (UTC): `2026-10-14T03:32:00Z`

### Findings

- Severity: none
- Finding: The worker changed only the root `.gitignore` and the selected step record. The exact root rule `/.cline-iterations/` matches iterator artifacts, and the isolated iterator session completed with a clean workspace without forced recovery.
- Evidence: `git show --stat bd1e289`, `git check-ignore -v .cline-iterations/example/session.json`, the completed iterator manifest, monitor phase summaries, and staged-diff checks.
- Affected paths: `.gitignore`, `.docs/work-items/WI-0001/steps/WI-0001-17-repository-gitignore.md`

### Disposition

- Decision: accepted
- Follow-up: WI-0001 is complete. Language-toolchain ignore patterns are deferred to WI-0002.