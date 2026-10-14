---
schema: language-practice-step/v1
id: WI-0004
step: "05"
status: proven
skill: language-practice-typescript
depends_on: ["01"]
target: typescript/bootstrap.sh
---

# WI-0004 step 05 — Document TypeScript/Node bootstrap setup

## Objective

Write the executable `./typescript/bootstrap.sh` for Linux Node/npm version
selection, local TypeScript installation, package metadata and lockfiles,
`tsconfig`, verification, and runtime/module boundaries.

## Commit subject

```text
WI-0004: 05: doc(typescript): document bootstrap setup
```

## Required boundaries

- Write `./typescript/bootstrap.sh` as the sole implementation target for this
  step; the user will inspect and run it explicitly from a terminal.
- Distinguish Node installation/version management from project-local npm
  dependencies and TypeScript installation.
- Prefer a project-local TypeScript dependency and `npx tsc` over a global
  compiler for reproducibility; explain why a global install is only an optional
  one-off tool.
- Verify `node --version`, `npm --version`, local `npx tsc --version`, package
  lockfile state, `tsconfig`, type checking, and a fresh-process runtime check.
- Record local Node 24.21.0/npm 11.19.0 and the absence of local `tsc` as
  observations, not repository support decisions.
- Do not execute the script, install Node/npm/TypeScript, modify shell profiles,
  or download npm packages while processing this step.

## Validation

```bash
test -x typescript/bootstrap.sh
bash -n typescript/bootstrap.sh
grep -q 'node --version' typescript/bootstrap.sh
grep -q 'npm --version' typescript/bootstrap.sh
grep -q 'npx tsc\|typescript' typescript/bootstrap.sh
grep -q 'package-lock.json\|lockfile' typescript/bootstrap.sh
grep -q 'tsconfig\|tsc' typescript/bootstrap.sh
grep -q 'nvm\|Node 24.21.0\|apt\|package' typescript/bootstrap.sh
git diff --check
```

The script must link to or cite official Node/npm and TypeScript setup guidance
in its comments or the shared `BOOTSTRAP.md`. Static checks do not execute the
script or prove Node, npm, registry, or TypeScript behavior.

## Prerequisite

Step 01 must be proven and reviewed. No Node installation, npm registry access,
or dependency download is authorized by this step.

## Agent handoff notes

### Agent execution

- Agent identity: API assistant
- Agent/session ID: not available
- Provider/model: not available
- Capability: repository implementation and static validation
- Use case: execute WI-0004 step 05 Act plan
- Repository root: `/home/michael/src/github/language-practice`
- Work-item ID: `WI-0004`
- Step file: `.docs/work-items/WI-0004/steps/WI-0004-05-typescript-bootstrap.md`
- Execution mode: Act
- Started at (UTC): 2026-10-14T21:30:00Z
- Finished at (UTC): interrupted before commit; corrective review completed at `2026-10-14T21:42:48Z`
- Effective system instructions summary: Implement only the selected step, preserve unrelated work, validate honestly, update this handoff, and commit only explicit in-scope paths.
- Effective task prompt summary: Create the TypeScript/Node bootstrap script, perform static-only validation without running installers or package downloads, record evidence limitations, and use the declared commit subject.

### Run reference

- Iterator run directory: `/tmp/cline-iterate-WI-0004-rebased/run/WI-0004-05-typescript-bootstrap.md/20261014_212733/`
- Run ID: `20261014_212733`
- Plan session ID: not available
- Act session ID: not available
- The iterator artifacts contain the session and event history for the interrupted implementation attempt.

### Work performed

- Rechecked the clean topic branch and proven step 01 dependency.
- Added only `typescript/bootstrap.sh` as the implementation target.
- The script now runs from the `typescript/` track directory, separates Node installation/version management from project-local npm state, prefers local TypeScript through `npx --no-install tsc`, documents package metadata and lockfiles, creates a strict NodeNext `tsconfig` only when absent, handles an empty source tree explicitly, and includes type-check and fresh-process runtime verification for a user-run bootstrap.
- No bootstrap script, installer, npm registry access, package download, global install, shell-profile mutation, or destructive cleanup was executed.

### Files changed

- `typescript/bootstrap.sh`
- `.docs/work-items/WI-0004/steps/WI-0004-05-typescript-bootstrap.md`

### Validation and evidence

- Command/result: `bash -n typescript/bootstrap.sh` — PASS; executable mode — PASS; required step content checks — PASS; `git diff --check` — PASS.
- Command/result: local compiler fixture with valid and invalid TypeScript — PASS; empty-project `TS18003` classification — PASS; mixed `TS18003` plus real diagnostic classification — PASS; isolated temporary runtime-check compilation — PASS.
- Evidence tier and limitation: static/documentation evidence only. Node/npm observations, TypeScript installation, compiler behavior, lockfile resolution, project runtime behavior, cross-host support, and learner understanding remain unproven.

### Decisions and references

- Step 01 remains the proven dependency.
- Local Node `24.21.0` and npm `11.19.0` are recorded as observations only; no repository support version was invented.
- The script cites official Node.js, npm, TypeScript Handbook, and TSConfig guidance.
- Declared commit subject: `WI-0004: 05: doc(typescript): document bootstrap setup`.

### Blockers and risks

- The operator reports a successful user-run of the corrected script; no command
  transcript was supplied, so portability, project correctness beyond that run,
  and learner understanding remain unverified.
- The bootstrap script's npm installation and generated project checks require explicit user-run network/package access outside the iterator.
- The script's example NodeNext configuration is a bootstrap default and must be reviewed against any future project module policy.

### Exact next actions

1. Preserve the user-generated root-level and `typescript/` project metadata as
   uncommitted user state unless the operator separately authorizes a project-
   metadata commit.
2. Retain the user-run success as host evidence only; do not promote it to a
   portability or project-support claim.

## Advisor review

### Review round 1 — corrective static review

- Advisor: API assistant
- Work-item: `WI-0004`
- Step: `05`
- Iterator run: `/tmp/cline-iterate-WI-0004-rebased/run/WI-0004-05-typescript-bootstrap.md/20261014_212733/`
- Reviewed snapshot: topic branch `wi-0004-toolchain-bootstrap`, working tree based on `9d3b4fb`
- Reviewed at (UTC): `2026-10-14T21:42:48Z`

### Findings

- Severity: major
- Finding: The user-run script reached the local TypeScript compiler successfully,
  but `tsconfig.json` matched no source files and TypeScript returned `TS18003`.
  The script also ran its project setup from the repository root, which caused
  the user-run `package.json`, `package-lock.json`, `tsconfig.json`, and
  `node_modules` to be created outside the intended `typescript/` track.
- Evidence: user-provided output; static inspection of the prior script; the
  current repository `.gitignore`; current generated root-level files; isolated
  local-compiler checks distinguishing empty-project `TS18003` from real
  `TS2322` failures; and Bash/static validation after the correction.
- Affected paths: `typescript/bootstrap.sh` and this selected step handoff.

### Correction and disposition

- Decision: accepted for static review; user-run follow-up required
- Correction: the script now changes into its own `typescript/` track directory,
  reports a distinct toolchain-verification success line, uses
  `npx --no-install` for the project-local compiler, treats only an otherwise
  clean `TS18003` no-input diagnostic as an explicit project-check skip, fails
  on all other TypeScript diagnostics, and compiles the temporary runtime check
  with an isolated temporary `tsconfig.json`.
- Validation: `bash -n`, executable-mode checks, required step greps,
  `git diff --check`, empty-project classification, real compiler success/failure
  discrimination, and isolated runtime-check compilation all passed. The
  bootstrap script itself was not rerun; package installation, network access,
  and user-host mutation remain unproven for the corrected version.
- Follow-up: review the previously created root-level Node artifacts manually;
  rerun `./typescript/bootstrap.sh` from the repository root after the selected
  files are committed, then record the resulting host/project evidence separately.

### Review round 2 — user-run outcome reconciliation

- Advisor: API assistant
- Work-item: `WI-0004`
- Step: `05`
- Reviewed at (UTC): `2026-10-14T22:38:05Z`

### Findings and disposition

- Severity: none for the reported outcome; evidence limitation remains.
- Finding: The operator reports that the corrected `typescript/bootstrap.sh`
  completed successfully when run by the user. No command transcript or
  independent host/project artifact was supplied, so this is user-reported host
  evidence only. The earlier root-level generated metadata remains uncommitted;
  the corrected script's intended track-local behavior is recorded but not
  independently replayed by the advisor.
- Decision: proven for static script scope and user-reported execution outcome;
  portability, project correctness, and learner understanding remain unproven.