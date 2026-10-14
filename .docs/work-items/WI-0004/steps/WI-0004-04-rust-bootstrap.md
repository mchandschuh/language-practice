---
schema: language-practice-step/v1
id: WI-0004
step: "04"
status: proven
skill: language-practice-rust
depends_on: ["01"]
target: rust/bootstrap.sh
---

# WI-0004 step 04 — Document Rust bootstrap setup

## Objective

Write the executable `./rust/bootstrap.sh` for Linux Rust installation and
toolchain selection through rustup, Cargo project creation, components, targets,
lockfiles, verification, and cleanup.

## Commit subject

```text
WI-0004: 04: doc(rust): document bootstrap setup
```

## Required boundaries

- Write `./rust/bootstrap.sh` as the sole implementation target for this step;
  the user will inspect and run it explicitly from a terminal.
- Explain stable/channel selection, `rustup show`, toolchain/component/target
  inspection, and the consequences of changing the default toolchain.
- Explain `cargo new`, `Cargo.toml`, `Cargo.lock`, `cargo check`, `cargo test`,
  `cargo fmt`, and `cargo clippy` as separate checks.
- State when a project should pin a toolchain or edition and identify unresolved
  repository decisions explicitly.
- Record that rustc/Cargo are unavailable locally and the Rust Book's 1.97/
  Edition 2024 context is not an execution contract.
- Do not execute the script, run rustup, download components, create Cargo
  projects, or change the user's default toolchain while processing this step.

## Validation

```bash
test -x rust/bootstrap.sh
bash -n rust/bootstrap.sh
grep -q 'rustup show' rust/bootstrap.sh
grep -q 'cargo new' rust/bootstrap.sh
grep -q 'Cargo.toml\|cargo' rust/bootstrap.sh
grep -q 'Cargo.lock\|--locked' rust/bootstrap.sh
grep -q 'rustup\|curl\|apt\|package' rust/bootstrap.sh
git diff --check
```

The script must link to or cite rustup and the official Cargo Book in its
comments or the shared `BOOTSTRAP.md`. Static checks do not execute rustup or
prove toolchain, component, target, or Cargo behavior.

## Prerequisite

Step 01 must be proven and reviewed. No rustup installation, component download,
or default-toolchain mutation is authorized by this step.

## Agent handoff notes

### Agent execution

- Agent identity: API assistant
- Agent/session ID: not available
- Provider/model: not available
- Capability: repository implementation and static validation
- Use case: execute WI-0004 step 04 Act plan
- Repository root: `/home/michael/src/github/language-practice`
- Work-item ID: `WI-0004`
- Step file: `.docs/work-items/WI-0004/steps/WI-0004-04-rust-bootstrap.md`
- Execution mode: Act
- Started at (UTC): 2026-10-14T05:50:00Z
- Finished at (UTC): 2026-10-14T05:50:00Z
- Effective system instructions summary: Implement only the selected step, preserve unrelated work, validate honestly, update this handoff, and commit only explicit in-scope paths.
- Effective task prompt summary: Create the Rust bootstrap script, perform static-only validation without running installers or Rust commands, record evidence limitations, and use the declared commit subject.

### Run reference

- Iterator run directory: `.cline-iterations/WI-0004-04-rust-bootstrap.md/` (current run reference not available)
- Run ID: not available
- Plan session ID: not available
- Act session ID: not available
- Use the run directory's session and event artifacts when available; prior corrective review evidence is under `20261014_054150/`.

### Work performed

- Rechecked that step 01 is proven and that Rust/Cargo are unavailable locally.
- Added only `rust/bootstrap.sh` as the step implementation target.
- The script provides visible prerequisite package installation, rustup installation, toolchain/component/target inspection, explicit default-toolchain mutation, Cargo project creation, lockfile-preserving checks, and bounded cleanup guidance.
- No bootstrap script, rustup installer, package manager, network download, `sudo`, Cargo command, project creation, or default-toolchain mutation was executed.

### Files changed

- `rust/bootstrap.sh`
- `.docs/work-items/WI-0004/steps/WI-0004-04-rust-bootstrap.md`

### Validation and evidence

- Command/result: `test -x`, `bash -n`, all required content greps, and `git diff --check` — PASS.
- Evidence tier and limitation: intended tier-1 static evidence only. Host installation, Rust/Cargo availability, component and target downloads, Cargo project behavior, and learner understanding remain unproven.

### Decisions and references

- Step 01 remains the proven dependency.
- `stable` is the script's explicit default channel, not a repository support decision.
- The script cites the Rust installation page, rustup documentation, and Cargo Book.
- Default toolchain selection is opt-in through `--default-toolchain`; installation and selection are separate operations.
- Declared commit subject: `WI-0004: 04: doc(rust): document bootstrap setup`.

### Blockers and risks

- Runtime execution is intentionally prohibited by the step and remains unverified.
- The repository has not selected a Rust channel or edition; users must make and record that decision before treating the setup as support policy.
- The script records a downloaded installer digest but does not claim independent authenticity verification beyond HTTPS/source visibility.

### Exact next actions

1. Run all required static checks and inspect the complete diff.
2. Stage only `rust/bootstrap.sh` and this selected step file.
3. Run `git diff --cached --check` and inspect staged paths.
4. Commit with the declared subject.
5. Verify the commit summary and clean worktree.

## Advisor review

The advisor appends review rounds here after the worker completes.

### Review round 1 — user-run outcome reconciliation

- Advisor: API assistant
- Work-item: `WI-0004`
- Step: `04`
- Reviewed at (UTC): `2026-10-14T22:38:05Z`

### Findings and disposition

- Severity: none for the reported outcome; evidence limitation remains.
- Finding: The operator reports that `rust/bootstrap.sh` completed successfully
  when run by the user. No command transcript or independent host/project
  artifact was supplied, so this is user-reported host evidence only.
- Decision: proven for static script scope and user-reported execution outcome;
  portability, project correctness, and learner understanding remain unproven.