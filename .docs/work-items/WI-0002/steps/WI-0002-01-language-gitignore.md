---
schema: language-practice-step/v1
id: WI-0002
step: "01"
status: proven
skill: repository-gitignore
depends_on: []
target: .gitignore
---

# WI-0002 step 01 — Ignore language-toolchain artifacts

## Objective

Extend the root `.gitignore` for generated and local-only artifacts from the
Go, Python, Rust, and TypeScript/Node learning tracks while preserving source,
tests, documentation, committed lockfiles, `.env.example`, and the WI-0001
`.cline-iterations/` rule.

## Commit subject

```text
WI-0002: 01: chore(repo): ignore language artifacts
```

## Scope

- Add only conservative patterns justified by the selected four language tracks.
- Keep `/\.cline-iterations/` intact.
- Prefer directory-scoped generated output and cache patterns over broad
  extensions that could hide source or test fixtures.
- Do not add language manifests or dependencies.

## Validation

```bash
git check-ignore -v \
  go/example.test \
  go/coverage.out \
  python/__pycache__/module.cpython-314.pyc \
  python/.venv/bin/python \
  rust/target/debug/example \
  typescript/node_modules/package/index.js \
  typescript/dist/index.js \
  typescript/.tsbuildinfo \
  .cline-iterations/example/session.json
git check-ignore -v --no-index .env.example
git check-ignore -v --no-index go/main.go
git check-ignore -v --no-index rust/Cargo.lock
git diff --check
```

The representative generated paths must match. Source, `.env.example`, and
lockfile checks must remain visible; use `--no-index` because these fixtures are
not all present or tracked in the repository.

## Prerequisite

The WI-0002 root documentation commit must be an ancestor of `HEAD`. WI-0001's
existing `.cline-iterations/` rule must remain present and unchanged in meaning.

## Agent handoff notes

### Agent execution

- Agent identity: API assistant (identity not available)
- Agent/session ID: `1791948865069_gsdfh`
- Provider/model: `openrouter / openai/gpt-5.6-luna`
- Capability: scoped repository implementation, validation, handoff, and commit
- Use case: Act execution for the selected repository gitignore step
- Repository root: `/home/michael/src/github/language-practice`
- Work-item ID: `WI-0002`
- Step file: `.docs/work-items/WI-0002/steps/WI-0002-01-language-gitignore.md`
- Execution mode: `Act`
- Started at (UTC): `2026-10-14T03:34:25Z`
- Finished at (UTC): `2026-10-14T03:35:09Z`
- Effective system instructions summary: Execute only the selected atomic step, preserve unrelated work, validate the final snapshot, update this handoff truthfully, and commit only explicit in-scope paths.
- Effective task prompt summary: Extend the root ignore file for conservative Go, Python, Rust, and TypeScript/Node artifacts, preserve source/configuration visibility and the existing iterator rule, run the declared checks, and commit with the exact step subject.

### Run reference

- Iterator run directory: `/tmp/language-practice-cline-runs/WI-0002-step01/WI-0002-01-language-gitignore.md/20261014_033331/`
- Run ID: `20261014_033331`
- Plan session ID: `1791948811544_3hxv7`
- Act session ID: `1791948865069_gsdfh`
- Review `session.json`, `events.jsonl`, monitor summaries, prompt snapshots, and transcripts in the external run directory for execution context.

### Work performed

- Rechecked the clean isolated worktree and confirmed the WI-0002 root documentation commit is an ancestor of `HEAD`.
- Added directory-scoped ignore rules for Go test and coverage artifacts, Python cache and virtual-environment artifacts, Rust Cargo output, and TypeScript/Node dependencies and build output.
- Preserved the existing `/.cline-iterations/` rule.

### Files changed

- `.gitignore`
- `.docs/work-items/WI-0002/steps/WI-0002-01-language-gitignore.md`

### Validation and evidence

- Pre-edit `git diff --check` — passed.
- `git check-ignore -v` for all eight language-artifact paths and `.cline-iterations/example/session.json` — passed; each matched the intended scoped rule.
- `git check-ignore -v --no-index` for `.env.example`, `go/main.go`, and `rust/Cargo.lock` — passed; each remained visible.
- `git diff --check` — passed after implementation.
- `git diff --cached --check` — passed before commit; the staged diff contained only the two listed paths.
- `cline-iterate` session manifest — completed; post-step workspace status was clean and no forced recovery was used.
- Evidence tier: static Git ignore behavior, worktree hygiene, and scoped delivery validation only; language-toolchain behavior, builds, runtime behavior, learner outcomes, and production safety are not proven.

### Decisions and references

- The step remained atomic and within its declared ownership boundary.
- Patterns are rooted in the corresponding language directories to avoid broad extension-based ignores.
- The exact required commit subject is `WI-0002: 01: chore(repo): ignore language artifacts`.
- Only the root `.gitignore` and this selected step record are permitted implementation paths.

### Blockers and risks

- No dependency or scope blocker identified.
- Representative language paths are synthetic because no language manifests or generated directories exist in this repository.
- Static checks cannot prove actual toolchain output or runtime behavior.

### Exact next actions

1. Inspect the primary checkpoint commit and append the advisor review before integration.

## Advisor review

### Review round 1

- Advisor: API assistant (identity not available)
- Work-item: `WI-0002`
- Step: `01`
- Iterator run: `/tmp/language-practice-cline-runs/WI-0002-step01/WI-0002-01-language-gitignore.md/20261014_033331/`
- Implementation commit: `01b2e2bb8ed428c4dbdfcde0414730b65c4b7079`
- Reviewed snapshot: primary branch `wi-0002-language-gitignore`, commit `01b2e2b`
- Reviewed at (UTC): `2026-10-14T03:37:00Z`

### Findings

- Severity: none
- Finding: The worker added only conservative, language-directory-scoped patterns for the selected Go, Python, Rust, and TypeScript/Node tracks. The existing `/.cline-iterations/` rule remains intact. Source files, `.env.example`, a Rust lockfile, and a TypeScript source file remain visible.
- Evidence: `git show --stat 01b2e2b`, the generated-artifact and visibility `git check-ignore` matrix, staged-diff checks, the completed iterator manifest, and monitor phase summaries.
- Affected paths: `.gitignore`, `.docs/work-items/WI-0002/steps/WI-0002-01-language-gitignore.md`

### Disposition

- Decision: accepted
- Follow-up: WI-0002 is complete. Toolchain manifests and runtime behavior remain outside scope.