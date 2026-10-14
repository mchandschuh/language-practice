---
schema: language-practice-step/v1
id: WI-0004
step: "01"
status: proven
skill: language-practice-bootstrap
depends_on: []
target: BOOTSTRAP.md
---

# WI-0004 step 01 — Add the bootstrap contract

## Objective

Define the shared Linux bootstrap contract for the five language tracks before
writing language-specific setup documents. Explain host prerequisites, user-local
versus system installation, version selection, permissions, project-local
environments, verification, cleanup, rollback, reproducibility, and evidence
limits.

## Commit subject

```text
WI-0004: 01: doc(repo): add bootstrap contract
```

## Scope

- Add only `/home/michael/src/github/language-practice/BOOTSTRAP.md` and this
  selected step's worker handoff update.
- Define the common document template and command-safety vocabulary used by
  steps 02–06.
- Define `./{lang}/bootstrap.sh` as the executable target convention, with the
  concrete language paths recorded in the root README.
- Include a language/toolchain matrix without pretending that current local
  versions are repository support decisions.
- State that scripts may contain user-run installation commands, but the
  iterator and advisor must only inspect and statically validate them.

## Required sections

- Purpose, learner profile, and assistance mode;
- host and Linux assumptions;
- state-change and permission classification;
- version-selection and source-authority policy;
- per-language bootstrap matrix;
- verification/evidence tiers;
- cleanup and rollback boundaries;
- troubleshooting/reporting contract; and
- executable `./{lang}/bootstrap.sh` ownership and static-validation rules; and
- explicit non-goals and unproven claims.

## Validation

```bash
test -f BOOTSTRAP.md
grep -q 'permission' BOOTSTRAP.md
grep -q 'rollback' BOOTSTRAP.md
grep -q 'verification' BOOTSTRAP.md
grep -q 'Python' BOOTSTRAP.md
grep -q 'Go' BOOTSTRAP.md
grep -q 'Rust' BOOTSTRAP.md
grep -q 'TypeScript' BOOTSTRAP.md
grep -q 'C++' BOOTSTRAP.md
git diff --check
```

No bootstrap script, installer, package manager, network download, or privileged
command may be run for this contract step.

## Prerequisite

The WI-0004 root documentation commit must be an ancestor of `HEAD`. The current
branch must be `wi-0004-toolchain-bootstrap`, and the worktree must be clean.

## Agent handoff notes

This section is reserved for the iterator worker. Record loaded skills, run
reference, files changed, validation, evidence tier, and limitations without
copying raw transcripts or secrets.

### Corrective implementation handoff

- Loaded skill: `language-practice-bootstrap`.
- Target: `BOOTSTRAP.md`; no language bootstrap script was created or edited.
- Files changed: `BOOTSTRAP.md` and this selected step handoff only.
- Contract coverage: all five `./<lang>/bootstrap.sh` targets, user-run
  installation, visible commands, permissions, idempotence, failure and
  rollback boundaries, verification tiers, troubleshooting, and unproven
  host-installation evidence are explicit.
- Validation: the required contract checks and `git diff --check` were run
  after implementation; no bootstrap script, installer, package manager,
  network download, `sudo`, profile mutation, or destructive cleanup was run.
- Evidence tier: static/documentation evidence only. Host installation,
  toolchain support, project correctness, and learner understanding remain
  unproven.
- Limitations: no repository language versions were invented; the official
  source links are references, not support commitments. Existing unrelated
  work and iterator artifacts were preserved and are not part of this change.

## Advisor review

### Review round 1

- Advisor: API assistant (identity not available)
- Work-item: `WI-0004`
- Step: `01`
- Iterator run: `/home/michael/src/github/language-practice/.cline-iterations/WI-0004-01-bootstrap-contract.md/20261014_045828/`
- Implementation commit: none
- Reviewed snapshot: topic branch `wi-0004-toolchain-bootstrap`, commit `f627883`
- Reviewed at (UTC): `2026-10-14T05:01:00Z`

### Findings

- Severity: blocker
- Finding: The iterator task prompt incorrectly told the worker to write only a
  language script, while the authoritative step target remained `BOOTSTRAP.md`.
  The worker correctly refused to substitute a sibling target and made no
  repository changes. The parent iterator was stopped before any later worker
  could write or commit a target.
- Evidence: completed run manifest/session and event log under
  `/home/michael/src/github/language-practice/.cline-iterations/WI-0004-01-bootstrap-contract.md/20261014_045828/`;
  step-02 inspection-only artifact under
  `/home/michael/src/github/language-practice/.cline-iterations/WI-0004-02-python-bootstrap.md/20261014_045926/`;
  clean Git status; absent `BOOTSTRAP.md` and language bootstrap scripts; no
  new commit after `f627883`.
- Affected paths: none in the repository; preserved iterator artifacts only.

### Disposition

- Decision: corrective step required
- Follow-up: Correct the iterator prompt so step 01 may create only `BOOTSTRAP.md`
  and its handoff, while steps 02–06 own only their declared `./{lang}/bootstrap.sh`
  targets. Rerun step 01 alone, review it, then process dependent script steps.

### Review round 2

- Advisor: API assistant (identity not available)
- Work-item: `WI-0004`
- Step: `01`
- Iterator run: `/home/michael/src/github/language-practice/.cline-iterations/WI-0004-01-bootstrap-contract.md/20261014_050053/`
- Implementation commit: `79f17bf07d0ba73508ce4ef4cd02b607bfd2ac57`
- Reviewed snapshot: topic branch `wi-0004-toolchain-bootstrap`, commit `79f17bf`
- Reviewed at (UTC): `2026-10-14T05:03:00Z`

### Findings

- Severity: none
- Finding: The corrected worker created the shared `BOOTSTRAP.md` contract and
  updated only this selected step handoff. It explicitly assigns the five
  `./{lang}/bootstrap.sh` targets to steps 02–06, requires visible user-run
  commands, and keeps iterator/advisor execution at static validation only.
- Evidence: iterator run `20261014_050053`, commit `79f17bf`, `git show --check`,
  exact step validation output, clean post-commit worktree, and inspection of
  `BOOTSTRAP.md` for permissions, version/source authority, idempotence, failure,
  rollback, verification tiers, and evidence limits. No scripts, installers,
  network access, package manager, `sudo`, profile mutation, or cleanup ran.
- Affected paths: `BOOTSTRAP.md` and this selected step handoff.

### Disposition

- Decision: accepted
- Follow-up: Step 01 is proven for the shared static bootstrap contract. Begin
  steps 02–06 serially; review each executable script statically and do not run
  it. Host installation and project behavior remain unproven pending user-run
  validation.