---
schema: language-practice-step/v1
id: WI-0001
step: "05"
status: proven
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

## Agent handoff notes

### Agent execution

- Agent identity: not available
- Agent/session ID: `1791939744174_egifq`
- Provider/model: `openrouter/openai/gpt-5.6-luna`
- Capability: tools enabled; spawn and teams disabled
- Use case: Act execution for the selected research-skill step
- Repository root: `/home/michael/src/github/language-practice`
- Work-item ID: `WI-0001`
- Step file: `.docs/work-items/WI-0001/steps/WI-0001-05-language-practice-research.md`
- Execution mode: act-only
- Started at (UTC): `2026-10-14T01:02:24.163Z`
- Finished at (UTC): `2026-10-14T01:05:00Z`
- Effective system instructions summary: Execute only the selected atomic step, preserve unrelated worktree changes, validate the final snapshot, update this handoff, and commit only explicit in-scope paths.
- Effective task prompt summary: Add the authoritative research skill with version-matching, official anchors for four languages, untrusted-content handling, required validation, evidence limits, and the exact step commit subject.

### Run reference

- Iterator run directory: `.cline-iterations/WI-0001-05-language-practice-research.md/20261014_010224/`
- Run ID: `20261014_010224`
- Plan session ID: not available
- Act session ID: `1791939744174_egifq`
- Use the run directory's `session.json`, `events.jsonl`, prompt snapshots, and Act transcript to review the execution context.

### Work performed

- Inspected the selected step, dependency skills, current worktree, and existing research skill.
- Preserved the already-present research skill content, normalized its final newline, and recorded this execution handoff.

### Files changed

- `.cline/skills/language-practice-research/SKILL.md`
- `.docs/work-items/WI-0001/steps/WI-0001-05-language-practice-research.md`

### Validation and evidence

- `npx --yes skills-ref@0.1.5 validate .cline/skills/language-practice-research` — passed (`Valid skill`).
- `git diff --check` — passed for the selected skill and step record before the final handoff edit; final staged diff check is required before commit.
- Evidence tier: static skill-schema and whitespace validation only; no runtime activation, prompt evaluation, learner outcome, or language-toolchain evidence is claimed.

### Decisions and references

- The step is atomic: one skill file plus its execution record.
- Dependencies 02 and 03 are satisfied by commits `2f14e66` and `c86081b`; bootstrap commit `ac528d5` is also present but is outside this step's commit scope.
- The exact required commit subject is `WI-0001: 05: feat(repo): add research skill`.

### Blockers and risks

- Pre-existing untracked files under `.cline-iterations/` and other skill directories remain outside this step and must not be staged.
- No blocker to this step's static validation or isolated commit.

### Exact next actions

1. Re-run the required validation after the handoff edit.
2. Stage only the research skill and this selected step record; inspect the staged diff and run `git diff --cached --check`.
3. Commit with the exact subject above, then inspect the commit summary and worktree status.
